package simkl

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type RefreshReason int

const (
	Requested RefreshReason = iota
	Explicit
)

type entry struct {
	list                List
	checked, downloaded time.Time
}
type Cache struct {
	mu              sync.Mutex
	client          *Client
	entries         map[int64]entry
	now             func() time.Time
	account         int64
	stamp           string
	index           map[int64]ListMeta
	activityChecked time.Time
	path            string
	details         map[string]detailEntry
}

func NewCache(c *Client) *Cache {
	return &Cache{client: c, entries: map[int64]entry{}, now: time.Now, index: map[int64]ListMeta{}, details: map[string]detailEntry{}}
}
func clone(l List) List {
	l.Items = append([]Item{}, l.Items...)
	if l.Counts != nil {
		v := *l.Counts
		l.Counts = &v
	}
	return l
}
func (c *Cache) Get(ctx context.Context, id int64, reason RefreshReason) (result List, failure error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer func() {
		if failure == nil && c.path != "" {
			if err := c.save(); err != nil {
				result = List{}
				failure = errors.New("Listenzwischenspeicher konnte nicht gespeichert werden")
			}
		}
	}()
	now := c.now()
	cached, exists := c.entries[id]
	if exists && reason != Explicit && now.Sub(cached.checked) < time.Minute {
		return clone(cached.list), nil
	}
	if !exists || reason == Explicit {
		return c.download(ctx, id, now)
	}
	if c.account == 0 {
		account, err := c.client.AccountID(ctx)
		if err != nil {
			return List{}, err
		}
		c.account = account
	}
	own := cached.list.User.ID == c.account
	candidateIndex := c.index
	stamp := c.stamp
	activityChecked := c.activityChecked
	if own {
		if now.Sub(c.activityChecked) >= time.Minute || c.activityChecked.IsZero() {
			activity, err := c.client.Activities(ctx)
			if err != nil {
				return List{}, err
			}
			stamp = activity.CustomLists.Lists.All
			if stamp == "" || stamp != c.stamp {
				lists, err := c.client.UserLists(ctx, c.account)
				if err != nil {
					return List{}, err
				}
				candidateIndex = map[int64]ListMeta{}
				for _, l := range lists {
					candidateIndex[l.ID] = l
				}
			}
			activityChecked = now
		}
	} else {
		// Fremde Listen folgen nicht dem Aktivitätszeitstempel des angemeldeten Kontos.
		lists, err := c.client.UserLists(ctx, cached.list.User.ID)
		if err != nil {
			return List{}, err
		}
		candidateIndex = map[int64]ListMeta{}
		for _, l := range lists {
			candidateIndex[l.ID] = l
		}
	}
	meta, known := candidateIndex[id]
	unchanged := known && meta.UpdatedAt != "" && meta.UpdatedAt == cached.list.UpdatedAt
	dailyDue := cached.list.Type == "auto" && now.Sub(cached.downloaded) >= 24*time.Hour
	if cached.list.Type == "auto" && !dailyDue && known && (meta.UpdatedAt == "" || meta.UpdatedAt == cached.list.UpdatedAt) {
		unchanged = true
	}
	var list List
	var err error
	if unchanged && !dailyDue {
		cached.checked = now
		c.entries[id] = cached
		list = clone(cached.list)
	} else {
		list, err = c.download(ctx, id, now)
		if err != nil {
			return List{}, err
		}
	}
	if own {
		c.index = candidateIndex
		c.stamp = stamp
		c.activityChecked = activityChecked
	}
	return list, nil
}
func (c *Cache) download(ctx context.Context, id int64, now time.Time) (List, error) {
	list, err := c.client.ListItems(ctx, id)
	if err != nil {
		return List{}, err
	}
	c.entries[id] = entry{list: clone(list), checked: now, downloaded: now}
	return clone(list), nil
}
func (c *Cache) DetailIDs(ctx context.Context, kind string, id int64) (IDs, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := fmt.Sprintf("%s/%d", kind, id)
	if entry, ok := c.details[key]; ok && c.now().Sub(entry.Fetched) < 24*time.Hour {
		return entry.IDs, nil
	}
	ids, err := c.client.DetailIDs(ctx, kind, id)
	if err == nil && (ids.TVDB > 0 || ids.TMDB > 0) {
		c.details[key] = detailEntry{IDs: ids, Fetched: c.now()}
	}
	return ids, err
}
