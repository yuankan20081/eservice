package publisher

import (
	"errors"
	"golang.org/x/net/context"
	"io"
	"log"
	"sync"
)

var (
	ErrDuplicateWriterForToken = errors.New("token already exist")
)

type Publisher struct {
	mu sync.RWMutex
	m  map[string]io.Writer

	errorChannel chan error
}

func New() *Publisher {
	return &Publisher{
		m:            make(map[string]io.Writer),
		errorChannel: make(chan error, 1000),
	}
}

func (pub *Publisher) Add(token string, w io.Writer) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if _, ok := pub.m[token]; ok {
		return ErrDuplicateWriterForToken
	}

	pub.m[token] = w

	return nil
}

func (pub *Publisher) Remove(token string) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()
	delete(pub.m, token)

	return nil
}

func (pub *Publisher) Write(p []byte) (int, error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	for token, w := range pub.m {
		go pub.publish(token, w, p)
	}

	return len(p), nil
}

func (pub *Publisher) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-pub.errorChannel:
			log.Println(err)
		}
	}
}

func (pub *Publisher) publish(token string, w io.Writer, p []byte) {
	if _, err := w.Write(p); err != nil {
		pub.errorChannel <- err
		pub.Remove(token)
	}
}
