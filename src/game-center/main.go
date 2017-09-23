package main

import (
	"encoding/json"
	"game-center/internal/license"
	"game-center/internal/rpc"
	"game-share/centerservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Cfg struct {
	LocalAddr  string
	LicenseDir string
}

func main() {
	bin, err := ioutil.ReadFile("./setup.json")
	if err != nil {
		log.Fatalln(err)
	}

	var cfg Cfg
	err = json.Unmarshal(bin, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	errorChannel := make(chan error, 10)

	lm := license.NewManager()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lm.Watch(ctx, cfg.LicenseDir); err != nil {
			errorChannel <- err
		}
	}()

	l, err := net.Listen("tcp", cfg.LocalAddr)
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
