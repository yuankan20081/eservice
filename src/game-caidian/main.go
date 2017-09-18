package main

import (
	"game-caidian/internal/agent"
	"game-caidian/internal/logic"
	"game-net/buffer/pool"
	"game-net/tcp-server"
	"game-net/tcp-session"
	"golang.org/x/net/context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l, err := net.Listen("tcp", ":12000")
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	errChannel := make(chan error, 1)
	var wg sync.WaitGroup
	ctx, cancelCtx := context.WithCancel(context.Background())

	// start game engine
	wg.Add(1)
	go func() {
		defer wg.Done()

		ge := logic.NewGameEngine()
		if err := ge.Serve(ctx); err != nil {
			errChannel <- err
		}
	}()

	// start server
	wg.Add(1)
	go func() {
		defer wg.Done()

		bufferPool := pool.New()

		s := tcp_server.New(tcp_server.RawConnHandleFunc(func(ctx context.Context, conn net.Conn) error {
			c := tcp_session.New(bufferPool, agent.NewHandler())

			defer c.Cleanup()

			return c.Handle(ctx, conn)
		}))

		if err := s.Serve(ctx, l, 1000); err != nil {
			errChannel <- err
		}
	}()

	// handle signal
	go func() {
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

		<-sigChannel

		cancelCtx()
	}()

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	case err := <-errChannel:
		cancelCtx()
		log.Fatalln(err)
	}

	wg.Wait()
	log.Println("server stopped")

}
