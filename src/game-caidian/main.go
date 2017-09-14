package main

import (
	"game-caidian/internal/agent"
	"game-net/buffer/pool"
	"game-net/tcp-server"
	"game-net/tcp-session"
	"golang.org/x/net/context"
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":12000")
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	bufferPool := pool.New()

	s := tcp_server.New(tcp_server.RawConnHandleFunc(func(ctx context.Context, conn net.Conn) error {
		c := tcp_session.New(bufferPool, agent.NewHandler())

		defer c.Cleanup()

		return c.Handle(ctx, conn)
	}))

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	if err := s.Serve(ctx, l, 1000); err != nil {
		log.Fatalln(err)
	}
}
