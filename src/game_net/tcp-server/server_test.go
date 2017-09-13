package tcp_server

import (
	"golang.org/x/net/context"
	"net"
	"testing"
	"time"
)

func TestTcpServer_Serve(t *testing.T) {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		t.Error(err)
	}
	defer l.Close()

	s := New(RawConnHandleFunc(func(ctx context.Context, conn net.Conn) error {
		time.Sleep(time.Second * 5)
		conn.Close()
		return nil
	}))

	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	err = s.Serve(ctx, l, 10)
	if err != nil {
		t.Error(err)
	}
}
