package main

import (
	"game-caidian/internal/agent"
	"game-caidian/internal/config/observer"
	"game-caidian/internal/logic"
	"game-net/tcp-server"
	"game-net/tcp-session"
	"game-util/publisher"
	"golang.org/x/net/context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	errChannel := make(chan error, 10)
	var wg sync.WaitGroup
	ctx, cancelCtx := context.WithCancel(context.Background())

	co := observer.New()
	wg.Add(1)
	go func() {
		err := co.Watch(ctx, "./setup.json", func(ctx context.Context, co *observer.Observer) {
			wg.Done()
		})
		if err != nil {
			log.Fatalln(err)
		}
	}()
	// wait until config initialize
	wg.Wait()

	l, err := net.Listen("tcp", co.Config().LocalAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	// debug limit
	ctx, cancelCtx = context.WithTimeout(ctx, time.Hour*2)
	log.Println("this is a debug server, will stop after 2 hours!!!")

	// start publisher
	pub := publisher.New()
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := pub.Serve(ctx); err != nil {
			errChannel <- err
		}
	}()

	// start game engine
	ge := logic.NewGameEngine(pub)
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := ge.Serve(ctx, co); err != nil {
			errChannel <- err
		}
	}()

	// start server
	wg.Add(1)
	go func() {
		defer wg.Done()

		s := tcp_server.New(tcp_server.RawConnHandleFunc(func(ctx context.Context, conn net.Conn) error {
			c := tcp_session.New(agent.NewReader(pub, ge, co))

			return c.Serve(ctx, conn)
		}))

		if err := s.Serve(ctx, l, co.Config().MaxConnection); err != nil {
			errChannel <- err
		}
	}()

	// handle signal
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println(ctx.Err())
	case err := <-errChannel:
		cancelCtx()
		log.Fatalln(err)
	case <-sigChannel:
		cancelCtx()
	}

	wg.Wait()
	log.Println("server stopped")

}
