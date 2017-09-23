package observer

import (
	"game-caidian/internal/config"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/net/context"
	"path/filepath"
	"sync/atomic"
)

type Observer struct {
	store *atomic.Value
	w     *fsnotify.Watcher
}

func New() *Observer {
	co := new(Observer)
	co.store = new(atomic.Value)
	co.w, _ = fsnotify.NewWatcher()
	return co
}

func (co *Observer) Watch(ctx context.Context, dir string, onInit func(ctx context.Context, co *Observer)) error {
	if !filepath.IsAbs(dir) {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		dir = absDir
	}

	err := co.w.Add(dir)
	if err != nil {
		return err
	}

	err = co.load(dir)
	if err != nil {
		return err
	}

	onInit(ctx, co)

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-co.w.Events:
			if event.Op == fsnotify.Write {
				co.load(event.Name)
			}
		case err := <-co.w.Errors:
			return err
		}
	}
}

func (co *Observer) Config() *config.Cfg {
	return co.store.Load().(*config.Cfg)
}

func (co *Observer) load(dir string) error {
	cfg, err := config.Load(dir)
	if err != nil {
		return err
	}

	co.store.Store(cfg)
	return nil
}
