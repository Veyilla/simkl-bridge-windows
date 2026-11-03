package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const testClient = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

type memoryStore struct {
	token Tokens
	err   error
}

func (s *memoryStore) Load() (Tokens, error) { return s.token, s.err }
func (s *memoryStore) Save(v Tokens) error   { s.token = v; return s.err }

type fastClock struct {
	now   time.Time
	waits []time.Duration
}

func (c *fastClock) Now() time.Time { return c.now }
func (c *fastClock) Wait(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.waits = append(c.waits, d)
	c.now = c.now.Add(d)
	return nil
}
func TestDeviceFlowReadOnlyAndPolling(t *testing.T) {
	polls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("client_id") != testClient {
			t.Error("Client-ID fehlt")
		}
		switch r.URL.Path {
		case "/oauth2/device":
			if r.Form.Get("scope") != "media:read" {
				t.Error("Zugriff nicht lesend")
			}
			w.Write([]byte(`{"device_code":"private-device","user_code":"ABCD-1234","verification_uri":"https://simkl.com/pin","expires_in":900,"interval":5}`))
		case "/oauth2/token":
			polls++
			if r.Form.Get("device_code") != "private-device" || r.Form.Get("grant_type") != "urn:ietf:params:oauth:grant-type:device_code" {
				t.Error("Falsches Verfahren")
			}
			if polls < 3 {
				w.WriteHeader(400)
				if polls == 1 {
					w.Write([]byte(`{"error":"authorization_pending"}`))
				} else {
					w.Write([]byte(`{"error":"slow_down"}`))
				}
				return
			}
			w.Write([]byte(`{"access_token":"test-access","refresh_token":"test-refresh","token_type":"Bearer","expires_in":604800,"scope":"media:read"}`))
		default:
			t.Error("Unerwarteter Aufruf", r.URL.Path)
		}
	}))
	defer server.Close()
	store := &memoryStore{}
	clock := &fastClock{now: time.Unix(1000, 0)}
	m := New(testClient, server.URL, server.Client(), store)
	m.clock = clock
	challenge, err := m.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	public, _ := json.Marshal(challenge)
	if strings.Contains(string(public), "private-device") {
		t.Fatal("Geheimer Code in Oberflächendaten")
	}
	if err = m.Complete(context.Background(), challenge); err != nil {
		t.Fatal(err)
	}
	if store.token.AccessToken != "test-access" || polls != 3 {
		t.Fatal("Anmeldung nicht gespeichert")
	}
	if got := clock.waits; len(got) != 3 || got[0] != 5*time.Second || got[1] != 5*time.Second || got[2] != 10*time.Second {
		t.Fatalf("Abfrageintervalle: %v", got)
	}
}
func TestConcurrentRefreshOnlyOnce(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		r.ParseForm()
		if r.Form.Get("refresh_token") != "old-refresh" {
			t.Error("Falscher Refresh-Token")
		}
		w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh","token_type":"Bearer","expires_in":604800,"scope":"media:read"}`))
	}))
	defer server.Close()
	store := &memoryStore{token: Tokens{AccessToken: "expired", RefreshToken: "old-refresh", ExpiresAt: time.Unix(1, 0)}}
	m := New(testClient, server.URL, server.Client(), store)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Go(func() {
			token, err := m.AccessToken(context.Background())
			if err != nil || token != "new-access" {
				t.Errorf("Token: %q %v", token, err)
			}
		})
	}
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatalf("%d Erneuerungen", calls.Load())
	}
}
func TestAuthErrorsDoNotLeakSecrets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"invalid_client","error_description":"secret-from-server"}`))
	}))
	defer server.Close()
	m := New(testClient, server.URL, server.Client(), &memoryStore{})
	_, err := m.Begin(context.Background())
	if !errors.Is(err, ErrInvalidClient) || strings.Contains(err.Error(), "secret-from-server") {
		t.Fatalf("Unsicherer Fehler: %v", err)
	}
}
func TestDeviceCancelAndExpiry(t *testing.T) {
	m := New(testClient, "http://127.0.0.1:1", http.DefaultClient, &memoryStore{})
	m.clock = &fastClock{now: time.Unix(1000, 0)}
	if err := m.Complete(context.Background(), Challenge{expiresAt: time.Unix(999, 0)}); !errors.Is(err, ErrExpired) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := m.Complete(ctx, Challenge{expiresAt: time.Unix(2000, 0), interval: 5 * time.Second}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
