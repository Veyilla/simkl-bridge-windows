package simkl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NumericID int64

func (n *NumericID) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*n = 0
		return nil
	}
	s := strings.Trim(string(b), "\"")
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v < 0 {
		return errors.New("Ungültige Medienkennung")
	}
	*n = NumericID(v)
	return nil
}

type IDs struct {
	Simkl NumericID `json:"simkl_id"`
	TVDB  NumericID `json:"tvdb"`
	TMDB  NumericID `json:"tmdb"`
	IMDB  string    `json:"imdb"`
}
type Item struct {
	Title     string `json:"title"`
	Year      int    `json:"year"`
	Type      string `json:"type"`
	AnimeType string `json:"anime_type"`
	IDs       IDs    `json:"ids"`
}
type Counts struct {
	Items int `json:"items"`
}

func (c *Counts) UnmarshalJSON(b []byte) error {
	var v struct {
		Items *int `json:"items"`
	}
	if json.Unmarshal(b, &v) != nil || v.Items == nil || *v.Items < 0 {
		return errors.New("Simkl-Zähler fehlen oder sind ungültig")
	}
	c.Items = *v.Items
	return nil
}

type Owner struct {
	ID int64 `json:"id"`
}
type ListMeta struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	MediaType string  `json:"media_type"`
	Type      string  `json:"type"`
	User      Owner   `json:"user"`
	UpdatedAt string  `json:"updated_at"`
	Counts    *Counts `json:"counts,omitempty"`
}
type Pagination struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func (p *Pagination) UnmarshalJSON(b []byte) error {
	var v struct {
		Limit      *int `json:"limit"`
		Page       *int `json:"page"`
		TotalItems *int `json:"total_items"`
		TotalPages *int `json:"total_pages"`
	}
	if json.Unmarshal(b, &v) != nil || v.Limit == nil || v.Page == nil || v.TotalItems == nil || v.TotalPages == nil || *v.Limit < 1 || *v.Limit > 500 || *v.Page < 1 || *v.TotalItems < 0 || *v.TotalPages < 0 {
		return errors.New("Simkl-Seitendaten fehlen oder sind ungültig")
	}
	p.Limit = *v.Limit
	p.Page = *v.Page
	p.TotalItems = *v.TotalItems
	p.TotalPages = *v.TotalPages
	return nil
}

type List struct {
	ListMeta
	Pagination Pagination `json:"pagination"`
	Items      []Item     `json:"items"`
}
type Activities struct {
	CustomLists struct {
		Lists struct {
			All string `json:"all"`
		} `json:"lists"`
	} `json:"custom_lists"`
}
type TokenSource interface {
	AccessToken(context.Context) (string, error)
	Refresh(context.Context, string) (string, error)
}
type APIError struct {
	Code   string
	Status int
}

func (e *APIError) Error() string {
	switch e.Code {
	case "premium_only":
		return "Simkl PRO oder VIP erforderlich"
	case "oauth2_token_required":
		return "AUTH-V2-Anmeldung erforderlich"
	case "rate_limited":
		return "Simkl-Abfragelimit erreicht; bitte später erneut versuchen"
	case "user_token_required":
		return "Simkl-Anmeldung erforderlich"
	case "private_list":
		return "Kein Zugriff auf diese Simkl-Liste"
	case "not_found":
		return "Simkl-Liste nicht gefunden"
	default:
		return "Simkl-Anfrage fehlgeschlagen"
	}
}

type Client struct {
	clientID, base string
	http           *http.Client
	auth           TokenSource
	mu             sync.Mutex
	next           time.Time
	wait           func(context.Context, time.Duration) error
}

func waitFor(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func New(clientID, base string, client *http.Client, auth TokenSource) *Client {
	if base == "" {
		base = "https://api.simkl.com"
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	safe := *client
	safe.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &Client{clientID: clientID, base: strings.TrimRight(base, "/"), http: &safe, auth: auth, wait: waitFor}
}
func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if query == nil {
		query = url.Values{}
	}
	query.Set("client_id", c.clientID)
	query.Set("app-name", "Simkl-Bridge")
	query.Set("app-version", "0.1.0")
	token, err := c.auth.AccessToken(ctx)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 2; attempt++ {
		if d := time.Until(c.next); d > 0 {
			if err = c.wait(ctx, d); err != nil {
				return err
			}
		}
		req, e := http.NewRequestWithContext(ctx, "GET", c.base+path+"?"+query.Encode(), nil)
		if e != nil {
			return errors.New("Ungültige Simkl-Anfrage")
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("User-Agent", "Simkl-Bridge/0.1.0")
		res, e := c.http.Do(req)
		c.next = time.Now().Add(110 * time.Millisecond)
		if e != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return errors.New("Simkl ist nicht erreichbar")
		}
		b, e := io.ReadAll(io.LimitReader(res.Body, (16<<20)+1))
		res.Body.Close()
		if e != nil || len(b) > 16<<20 {
			return errors.New("Simkl-Antwort unvollständig")
		}
		if res.StatusCode == 401 && attempt == 0 {
			token, err = c.auth.Refresh(ctx, token)
			if err != nil {
				return err
			}
			continue
		}
		if res.StatusCode == 429 {
			delay := time.Minute
			if seconds, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && seconds > 0 {
				delay = time.Duration(seconds) * time.Second
			} else if until, e := http.ParseTime(res.Header.Get("Retry-After")); e == nil && time.Until(until) > 0 {
				delay = time.Until(until)
			}
			c.next = time.Now().Add(delay)
			return &APIError{Code: "rate_limited", Status: 429}
		}
		var status struct {
			Error string `json:"error"`
		}
		json.Unmarshal(b, &status)
		if status.Error != "" || res.StatusCode != 200 {
			code := status.Error
			if res.StatusCode == 404 {
				code = "not_found"
			}
			if res.StatusCode == 401 {
				code = "user_token_required"
			}
			return &APIError{Code: code, Status: res.StatusCode}
		}
		if err = json.Unmarshal(b, out); err != nil {
			return errors.New("Ungültige Simkl-Antwort")
		}
		return nil
	}
	return &APIError{Code: "user_token_required", Status: 401}
}
func (c *Client) AccountID(ctx context.Context) (int64, error) {
	var v struct {
		Account struct {
			ID int64 `json:"id"`
		} `json:"account"`
	}
	if err := c.get(ctx, "/users/settings", nil, &v); err != nil {
		return 0, err
	}
	if v.Account.ID <= 0 {
		return 0, errors.New("Simkl-Konto konnte nicht erkannt werden")
	}
	return v.Account.ID, nil
}
func (c *Client) Activities(ctx context.Context) (Activities, error) {
	var v Activities
	err := c.get(ctx, "/sync/activities", nil, &v)
	return v, err
}
func (c *Client) UserLists(ctx context.Context, userID int64) ([]ListMeta, error) {
	if userID <= 0 {
		return nil, errors.New("Ungültiges Simkl-Konto")
	}
	var lists []ListMeta
	total := -1
	pages := 1
	for page := 1; page <= pages; page++ {
		var v struct {
			Pagination Pagination `json:"pagination"`
			Lists      []ListMeta `json:"lists"`
		}
		q := url.Values{"limit": {"500"}, "page": {strconv.Itoa(page)}, "followed": {"true"}, "collaborants": {"true"}}
		if err := c.get(ctx, fmt.Sprintf("/lists/user/%d", userID), q, &v); err != nil {
			return nil, err
		}
		if v.Lists == nil || v.Pagination.Page != page || v.Pagination.TotalPages < 0 || v.Pagination.TotalPages > 20 || v.Pagination.TotalItems < 0 {
			return nil, errors.New("Listenübersicht unvollständig")
		}
		if page == 1 {
			total = v.Pagination.TotalItems
			pages = max(1, v.Pagination.TotalPages)
		} else if total != v.Pagination.TotalItems || pages != max(1, v.Pagination.TotalPages) {
			return nil, errors.New("Listenübersicht hat sich während des Abrufs geändert")
		}
		lists = append(lists, v.Lists...)
	}
	if len(lists) != total {
		return nil, errors.New("Listenübersicht unvollständig")
	}
	if lists == nil {
		lists = []ListMeta{}
	}
	return lists, nil
}
func (c *Client) ListItems(ctx context.Context, id int64) (List, error) {
	if id <= 0 {
		return List{}, errors.New("Ungültige Listen-ID")
	}
	var result List
	pages := 1
	seen := map[NumericID]bool{}
	for page := 1; page <= pages; page++ {
		var v List
		if err := c.get(ctx, fmt.Sprintf("/lists/%d", id), url.Values{"limit": {"500"}, "page": {strconv.Itoa(page)}}, &v); err != nil {
			return List{}, err
		}
		if v.ID != id || v.Items == nil || v.Counts == nil || v.Counts.Items != v.Pagination.TotalItems || v.Pagination.Page != page || v.Pagination.TotalItems < 0 || v.Pagination.TotalPages < 0 || v.Pagination.TotalPages > 20 || (v.Pagination.TotalItems > 0 && v.Pagination.TotalPages == 0) {
			return List{}, errors.New("Simkl-Liste unvollständig")
		}
		if v.MediaType != "movies" && v.MediaType != "tv" && v.MediaType != "anime" {
			return List{}, errors.New("Medienart dieser Liste wird nicht unterstützt")
		}
		if page == 1 {
			result = v
			result.Items = []Item{}
			pages = max(1, v.Pagination.TotalPages)
		} else if v.UpdatedAt != result.UpdatedAt || v.MediaType != result.MediaType || v.Pagination.TotalItems != result.Pagination.TotalItems || max(1, v.Pagination.TotalPages) != pages {
			return List{}, errors.New("Simkl-Liste hat sich während des Abrufs geändert")
		}
		for _, item := range v.Items {
			if err := ValidateItem(v.MediaType, item); err != nil {
				return List{}, err
			}
			if item.IDs.Simkl <= 0 || seen[item.IDs.Simkl] {
				return List{}, errors.New("Simkl-Liste enthält unvollständige oder doppelte Einträge")
			}
			seen[item.IDs.Simkl] = true
			result.Items = append(result.Items, item)
		}
	}
	if len(result.Items) != result.Pagination.TotalItems {
		return List{}, errors.New("Simkl-Liste unvollständig")
	}
	return result, nil
}
func ValidateItem(media string, item Item) error {
	expected := media
	if media == "movies" {
		expected = "movie"
	}
	if item.Type != expected {
		return errors.New("Medienart des Eintrags widerspricht der Liste")
	}
	if media == "anime" {
		switch item.AnimeType {
		case "tv", "movie", "ova", "ona", "special":
		default:
			return errors.New("Anime-Untertyp wird nicht unterstützt")
		}
	}
	return nil
}
func (c *Client) DetailIDs(ctx context.Context, kind string, id int64) (IDs, error) {
	if (kind != "movies" && kind != "tv" && kind != "anime") || id <= 0 {
		return IDs{}, errors.New("Ungültiger Detailabruf")
	}
	var v struct {
		IDs IDs `json:"ids"`
	}
	err := c.get(ctx, fmt.Sprintf("/%s/%d", kind, id), url.Values{"extended": {"full"}}, &v)
	return v.IDs, err
}
