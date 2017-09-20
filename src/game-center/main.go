package main

import (
	"game-center/internal/rpc"
	"game-share/centerservice"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":41000")
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	s := grpc.NewServer()
	defer s.Stop()

	service := new(rpc.CenterService)

	centerservice.RegisterCenterServiceServer(s, service)

	log.Fatalln(s.Serve(l))
}
