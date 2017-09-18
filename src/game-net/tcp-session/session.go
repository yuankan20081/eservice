package tcp_session

import (
	"log"
	"runtime/debug"
	"sync"
	"context"
	"net"
	"io"
)

type TcpSession struct {
	wg sync.WaitGroup
	errorChannel chan error
}

func New() *TcpSession {
	return &TcpSession{
		errorChannel: make(chan error, 1)
	}
}

func (ts *TcpSession) Serve(ctx context.Context, conn net.Conn) error{

}

func (ts *TcpSession) doRead(ctx context.Context, r io.Reader){
	defer ts.wg.Done()

}



func withRecover(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			debug.PrintStack()
		}
	}()

	fn()
}
