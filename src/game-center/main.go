package main

import (
	"game-center/internal/license"
	"game-center/internal/rpc"
	"game-share/centerservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	errorChannel := make(chan error, 10)

	lm := license.NewManager()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lm.Watch(ctx, "./license"); err != nil {
			errorChannel <- err
		}
	}()

	l, err := net.Listen("tcp", ":41000")
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	s := grpc.NewServer()
	defer s.Stop()

	service := rpc.NewCenterService(nil, lm)

	centerservice.RegisterCenterServiceServer(s, service)

	go func() {
		if err := s.Serve(l); err != nil {
			errorChannel <- err
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		wg.Wait()
		log.Println("server stoped")
	}()

	select {
	case <-ctx.Done():
	case err := <-errorChannel:
		cancelCtx()
		log.Fatalln(err)
	case <-sigChannel:
		cancelCtx()
	}
}
