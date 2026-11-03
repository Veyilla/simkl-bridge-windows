package simkl

import (
	"encoding/json"
	"errors"
	"os"
	"simklbridge/internal/config"
	"time"
)

type savedEntry struct {
	List       List      `json:"list"`
	Checked    time.Time `json:"checked"`
	Downloaded time.Time `json:"downloaded"`
}
type detailEntry struct {
	IDs     IDs       `json:"ids"`
	Fetched time.Time `json:"fetched"`
}
type snapshot struct {
	Version         int                    `json:"version"`
	ClientID        string                 `json:"client_id"`
	Account         int64                  `json:"account"`
	Entries         map[int64]savedEntry   `json:"entries"`
	Index           map[int64]ListMeta     `json:"index"`
	Stamp           string                 `json:"stamp"`
	ActivityChecked time.Time              `json:"activity_checked"`
	Details         map[string]detailEntry `json:"details"`
}

func NewPersistentCache(client *Client, path string, account int64) (*Cache, error) {
	c := NewCache(client)
	c.path = path
	c.account = account
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return c, nil
	}
	if err != nil {
		return nil, errors.New("Listenzwischenspeicher nicht lesbar")
	}
	var stored snapshot
	if len(raw) > 32<<20 || json.Unmarshal(raw, &stored) != nil || stored.Version != 1 || stored.Account != account || stored.ClientID != client.clientID {
		return c, nil
	}
	now := time.Now()
	for id, e := range stored.Entries {
		if e.List.ID != id || e.List.Counts == nil || e.List.Items == nil || len(e.List.Items) != e.List.Counts.Items || e.Checked.After(now) || e.Downloaded.After(now) {
			continue
		}
		c.entries[id] = entry{list: clone(e.List), checked: e.Checked, downloaded: e.Downloaded}
	}
	if stored.Index != nil {
		c.index = stored.Index
	}
	if stored.Details != nil {
		c.details = stored.Details
	}
	c.stamp = stored.Stamp
	if !stored.ActivityChecked.After(now) {
		c.activityChecked = stored.ActivityChecked
	}
	return c, nil
}
func (c *Cache) save() error {
	data := snapshot{Version: 1, ClientID: c.client.clientID, Account: c.account, Entries: map[int64]savedEntry{}, Index: c.index, Stamp: c.stamp, ActivityChecked: c.activityChecked, Details: c.details}
	for id, e := range c.entries {
		data.Entries[id] = savedEntry{List: e.list, Checked: e.checked, Downloaded: e.downloaded}
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return config.AtomicWrite(c.path, b)
}
