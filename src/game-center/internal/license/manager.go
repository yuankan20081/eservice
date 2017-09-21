package license

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	licensecipher "game-util/license"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("license not found")
)

type Manager struct {
	mu sync.RWMutex
	m  map[string]*License
}

func NewManager() *Manager {
	return &Manager{
		m: make(map[string]*License),
	}
}

func (m *Manager) Watch(ctx context.Context, dir string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}

	w.Add(dir)

	// init
	go m.loadDirectory(dir)

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-w.Events:
			switch event.Op {
			case fsnotify.Create:
				file := event.Name
				time.AfterFunc(time.Second*3, func() {
					m.addLicense(file)
				})
			case fsnotify.Remove:
				m.removeLicense(event.Name)
			}
		case err := <-w.Errors:
			return err
		}
	}
}

func (m *Manager) Get(token string) (*License, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if l, ok := m.m[token]; ok {
		lc := *l
		return &lc, nil
	} else {
		return nil, ErrNotFound
	}
}

func (m *Manager) loadDirectory(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			m.addLicense(path)
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) addLicense(dir string) {
	bin, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Println(err)
		return
	}

	barr := md5.Sum(bin)
	token := hex.EncodeToString(barr[:])

	bin, err = base64.StdEncoding.DecodeString(string(bin))
	if err != nil {
		log.Println(err)
		return
	}

	bin, err = licensecipher.RsaDecrypt(bin)
	if err != nil {
		log.Println(err)
		return
	}

	l := new(License)
	if err = json.Unmarshal(bin, &l); err != nil {
		log.Println(err)
		return
	}
	l.Dir = dir

	if time.Now().After(l.Expire) {
		log.Println("license expired:", filepath.Base(dir), l.Expire)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.m[token]; ok {
		log.Println("already in memory, should not see this")
		return
	}

	m.m[token] = l
	log.Println("new license add:", l)
	time.AfterFunc(l.Expire.Sub(time.Now()), func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.m, token)
		log.Println("license expired", filepath.Base(l.Dir))
	})
}

func (m *Manager) removeLicense(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for token, l := range m.m {
		if l.Dir == dir {
			delete(m.m, token)
			log.Println("license removed", l)
			return
		}
	}
}
