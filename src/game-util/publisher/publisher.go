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
	m  map[io.Writer]int

	errorChannel chan error
}

func New() *Publisher {
	return &Publisher{
		m:            make(map[io.Writer]int),
		errorChannel: make(chan error, 1000),
	}
}

func (pub *Publisher) Add(w io.Writer) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if _, ok := pub.m[w]; ok {
		return ErrDuplicateWriterForToken
	}

	pub.m[w] = 0

	return nil
}

func (pub *Publisher) Remove(w io.Writer) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()
	delete(pub.m, w)

	return nil
}

func (pub *Publisher) Write(p []byte) (int, error) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	for w, _ := range pub.m {
		go pub.publish(w, p)
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

func (pub *Publisher) publish(w io.Writer, p []byte) {
	if _, err := w.Write(p); err != nil {
		pub.errorChannel <- err
		pub.Remove(w)
	}
}
