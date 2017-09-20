package publisher

import (
	"errors"
	"game-net/writer"
	"golang.org/x/net/context"
	"log"
	"sync"
)

var (
	ErrDuplicateWriterForToken = errors.New("token already exist")
)

type Publisher struct {
	mu sync.RWMutex
	m  map[writer.Writer]int

	errorChannel chan error
}

func New() *Publisher {
	return &Publisher{
		m:            make(map[writer.Writer]int),
		errorChannel: make(chan error, 1000),
	}
}

func (pub *Publisher) Add(w writer.Writer) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	if _, ok := pub.m[w]; ok {
		return ErrDuplicateWriterForToken
	}

	pub.m[w] = 0

	return nil
}

func (pub *Publisher) Remove(w writer.Writer) error {
	pub.mu.Lock()
	defer pub.mu.Unlock()
	delete(pub.m, w)

	return nil
}

func (pub *Publisher) Publish(proto uint16, body interface{}) {
	pub.mu.Lock()
	defer pub.mu.Unlock()

	for w := range pub.m {
		go w.WriteResponse(proto, body)
	}

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
