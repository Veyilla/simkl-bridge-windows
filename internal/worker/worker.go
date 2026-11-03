package worker

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"simklbridge/internal/auth"
	"simklbridge/internal/config"
	"simklbridge/internal/diagnostics"
	"simklbridge/internal/feed"
	"simklbridge/internal/simkl"
	"simklbridge/internal/supervisor"
)

type service struct {
	cache  *simkl.Cache
	fatal  chan error
	status *diagnostics.Store
}

func actionRequired(err error) error {
	var failure error
	var api *simkl.APIError
	if errors.Is(err, auth.ErrAuthentication) {
		failure = supervisor.ErrAuthentication
	} else if errors.Is(err, auth.ErrInvalidClient) {
		failure = supervisor.ErrConfiguration
	} else if errors.Is(err, auth.ErrScope) {
		failure = supervisor.ErrConfiguration
	} else if errors.As(err, &api) {
		if api.Status == 401 {
			failure = supervisor.ErrAuthentication
		} else if api.Code == "oauth2_token_required" || api.Code == "premium_only" {
			failure = supervisor.ErrConfiguration
		}
	}
	return failure
}
func (s *service) check(err error) {
	failure := actionRequired(err)
	if failure != nil {
		select {
		case s.fatal <- failure:
		default:
		}
	}
}
func (s *service) Report(id int64, err error) {
	if s.status != nil {
		s.status.Record(id, err)
	}
}
func (s *service) Get(ctx context.Context, id int64, reason simkl.RefreshReason) (simkl.List, error) {
	list, err := s.cache.Get(ctx, id, reason)
	s.check(err)
	return list, err
}
func (s *service) DetailIDs(ctx context.Context, kind string, id int64) (simkl.IDs, error) {
	ids, err := s.cache.DetailIDs(ctx, kind, id)
	s.check(err)
	return ids, err
}
func Run(ctx context.Context, dir string) (result error) {
	logfile, logErr := os.OpenFile(filepath.Join(dir, "bridge.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if logErr != nil {
		return supervisor.ErrConfiguration
	}
	defer logfile.Close()
	status := diagnostics.New(dir, os.Getpid(), logfile)
	status.Record(0, nil)
	defer func() {
		if result != nil {
			status.Record(0, result)
		}
	}()
	c, err := (config.Store{Path: filepath.Join(dir, "settings.json")}).Load()
	if err != nil {
		return supervisor.ErrConfiguration
	}
	if c.ClientID == "" {
		return supervisor.ErrAuthentication
	}
	if len(c.Lists) == 0 {
		return supervisor.ErrConfiguration
	}
	tokens := auth.ProtectedStore{Path: auth.StorePath(dir, c.ClientID), ClientID: c.ClientID}
	manager := auth.New(c.ClientID, "https://api.simkl.com", nil, tokens)
	if _, err = manager.AccessToken(ctx); err != nil {
		if action := actionRequired(err); action != nil {
			return action
		}
		return err
	}
	listener, err := feed.Listen(c.Port)
	if err != nil {
		return supervisor.ErrConfiguration
	}
	defer listener.Close()
	lifetime, cancel := context.WithCancel(ctx)
	defer cancel()
	api := simkl.New(c.ClientID, "", nil, manager)
	account, err := api.AccountID(ctx)
	if err != nil {
		if action := actionRequired(err); action != nil {
			return action
		}
		return err
	}
	cache, err := simkl.NewPersistentCache(api, filepath.Join(dir, "lists-cache.json"), account)
	if err != nil {
		return supervisor.ErrConfiguration
	}
	source := &service{cache: cache, fatal: make(chan error, 1), status: status}
	done := make(chan error, 1)
	go func() { done <- feed.Serve(lifetime, listener, feed.Handler(source, c.Lists)) }()
	select {
	case err = <-source.fatal:
		cancel()
		<-done
		return err
	case err = <-done:
		return err
	case <-ctx.Done():
		cancel()
		<-done
		return nil
	}
}
