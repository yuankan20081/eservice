package main

import (
	"game-caidian/internal/agent"
	"game-caidian/internal/logic"
	"game-net/tcp-server"
	"game-net/tcp-session"
	"golang.org/x/net/context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"game-util/publisher"
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

	// start publisher
	pub := publisher.New()
	wg.Add(1)
	go func(){
		defer wg.Done()

		if err := pub.Serve(ctx); err!=nil{
			errChannel <- err			
		}
	}()

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

		s := tcp_server.New(tcp_server.RawConnHandleFunc(func(ctx context.Context, conn net.Conn) error {
			c := tcp_session.New(agent.NewReader(pub))

			return c.Serve(ctx, conn)
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
