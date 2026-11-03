package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"simklbridge/internal/config"
)

var ErrInvalidClient = errors.New("AUTH-V2-Client-ID prüfen")
var ErrAuthentication = errors.New("Simkl-Anmeldung erforderlich")
var ErrExpired = errors.New("Anmeldecode abgelaufen; bitte erneut verbinden")
var ErrScope = errors.New("Lesende Simkl-Berechtigung erforderlich")

type Tokens struct {
	ClientID     string    `json:"client_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (Tokens) String() string { return "[geschützte Anmeldung]" }

type TokenStore interface {
	Load() (Tokens, error)
	Save(Tokens) error
}
type Clock interface {
	Now() time.Time
	Wait(context.Context, time.Duration) error
}
type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now() }
func (wallClock) Wait(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

type Challenge struct {
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	deviceCode      string
	expiresAt       time.Time
	interval        time.Duration
}

func (Challenge) String() string { return "[Anmeldevorgang]" }

type Manager struct {
	mu             sync.Mutex
	clientID, base string
	http           *http.Client
	store          TokenStore
	clock          Clock
}

func New(clientID, base string, client *http.Client, store TokenStore) *Manager {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	safe := *client
	safe.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &Manager{clientID: clientID, base: strings.TrimRight(base, "/"), http: &safe, store: store, clock: wallClock{}}
}

type response struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	TokenType       string `json:"token_type"`
	Scope           string `json:"scope"`
	Error           string `json:"error"`
}

func (m *Manager) request(ctx context.Context, path string, form url.Values) (response, error) {
	if !config.ValidClientID(m.clientID) {
		return response{}, ErrInvalidClient
	}
	form.Set("client_id", m.clientID)
	req, err := http.NewRequestWithContext(ctx, "POST", m.base+path, strings.NewReader(form.Encode()))
	if err != nil {
		return response{}, errors.New("Simkl-Anfrage konnte nicht erstellt werden")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Simkl-Bridge/0.1.0")
	res, err := m.http.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return response{}, ctx.Err()
		}
		return response{}, errors.New("Simkl ist nicht erreichbar")
	}
	defer res.Body.Close()
	var body response
	if err = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&body); err != nil {
		return response{}, errors.New("Ungültige Simkl-Antwort")
	}
	if body.Error == "invalid_client" || body.Error == "unauthorized_client" {
		return body, ErrInvalidClient
	}
	if body.Error != "" {
		return body, nil
	}
	if res.StatusCode != 200 {
		return response{}, errors.New("Simkl-Anfrage fehlgeschlagen")
	}
	return body, nil
}
func (m *Manager) Begin(ctx context.Context) (Challenge, error) {
	body, err := m.request(ctx, "/oauth2/device", url.Values{"scope": {"media:read"}})
	if err != nil {
		return Challenge{}, err
	}
	uri, parseErr := url.Parse(body.VerificationURI)
	if body.Error != "" || body.DeviceCode == "" || body.UserCode == "" || body.ExpiresIn <= 0 || body.ExpiresIn > 900 || body.Interval < 1 || parseErr != nil || uri.Scheme != "https" || uri.Host != "simkl.com" || uri.Path != "/pin" || uri.User != nil || uri.RawQuery != "" || uri.Fragment != "" {
		return Challenge{}, errors.New("Ungültige Simkl-Anmeldeantwort")
	}
	interval := time.Duration(body.Interval) * time.Second
	if interval < 5*time.Second {
		interval = 5 * time.Second
	}
	return Challenge{UserCode: body.UserCode, VerificationURI: body.VerificationURI, deviceCode: body.DeviceCode, expiresAt: m.clock.Now().Add(time.Duration(body.ExpiresIn) * time.Second), interval: interval}, nil
}
func (m *Manager) Complete(ctx context.Context, ch Challenge) error {
	interval := ch.interval
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if !m.clock.Now().Before(ch.expiresAt) {
			return ErrExpired
		}
		if interval <= 0 {
			return ErrExpired
		}
		if err := m.clock.Wait(ctx, interval); err != nil {
			return err
		}
		if !m.clock.Now().Before(ch.expiresAt) {
			return ErrExpired
		}
		body, err := m.request(ctx, "/oauth2/token", url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:device_code"}, "device_code": {ch.deviceCode}})
		if err != nil {
			return err
		}
		switch body.Error {
		case "authorization_pending":
			continue
		case "slow_down":
			interval += 5 * time.Second
			continue
		case "expired_token":
			return ErrExpired
		case "":
			m.mu.Lock()
			err = m.persist(body)
			m.mu.Unlock()
			return err
		default:
			return ErrAuthentication
		}
	}
}
func (m *Manager) persist(body response) error {
	if body.AccessToken == "" || body.RefreshToken == "" || body.ExpiresIn <= 0 || !strings.EqualFold(body.TokenType, "Bearer") {
		return ErrAuthentication
	}
	if strings.TrimSpace(body.Scope) != "media:read" {
		return ErrScope
	}
	token := Tokens{ClientID: m.clientID, AccessToken: body.AccessToken, RefreshToken: body.RefreshToken, ExpiresAt: m.clock.Now().Add(time.Duration(body.ExpiresIn) * time.Second)}
	if err := m.store.Save(token); err != nil {
		return errors.New("Simkl-Anmeldung konnte nicht sicher gespeichert werden")
	}
	return nil
}
func (m *Manager) AccessToken(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.access(ctx, "")
}

// Refresh erneuert nach einem 401 höchstens einmal pro bisherigem Token.
func (m *Manager) Refresh(ctx context.Context, previous string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.access(ctx, previous)
}
func (m *Manager) access(ctx context.Context, previous string) (string, error) {
	token, err := m.store.Load()
	if err != nil || (token.ClientID != "" && token.ClientID != m.clientID) {
		return "", ErrAuthentication
	}
	if token.AccessToken != "" && token.ExpiresAt.After(m.clock.Now().Add(time.Minute)) && (previous == "" || previous != token.AccessToken) {
		return token.AccessToken, nil
	}
	if token.RefreshToken == "" {
		return "", ErrAuthentication
	}
	body, err := m.request(ctx, "/oauth2/token", url.Values{"grant_type": {"refresh_token"}, "refresh_token": {token.RefreshToken}})
	if err != nil {
		return "", err
	}
	if body.Error != "" {
		return "", ErrAuthentication
	}
	if err = m.persist(body); err != nil {
		return "", err
	}
	return body.AccessToken, nil
}
