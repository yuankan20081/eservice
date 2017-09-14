package tcp_server

import (
	"bytes"
	"golang.org/x/net/context"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

type RawConnHandler interface {
	Handle(ctx context.Context, conn net.Conn) error
}
type RawConnHandleFunc func(ctx context.Context, conn net.Conn) error

func (fn RawConnHandleFunc) Handle(ctx context.Context, conn net.Conn) error {
	return fn(ctx, conn)
}

type TcpServer struct {
	ctx     context.Context
	handler RawConnHandler
	wg      sync.WaitGroup
	sem     chan int
}

func New(handler RawConnHandler) *TcpServer {
	return &TcpServer{
		handler: handler,
	}
}

func (ts *TcpServer) Serve(ctx context.Context, listener net.Listener, maxConnection int) (e error) {
	ctx, cancel := context.WithCancel(ctx)
	ts.ctx = ctx

	ts.sem = make(chan int, maxConnection)
	for i := 0; i < maxConnection; i++ {
		ts.sem <- 1
	}

loop:
	for {
		listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Millisecond * 100))
		conn, err := listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
				e = err
				cancel()
				break
			} else {
				select {
				case <-ctx.Done():
					break loop
				default:

				}
			}
		} else {
			select {
			case <-ctx.Done():
				break loop
			default:
				select {
				case <-ts.sem:
					go ts.withRecoverHandleRawConnAsync(conn)
				default:
					conn.Close()
					log.Println("too many connections to handle")
				}
			}
		}
	}

	ts.wg.Wait()
	return
}

func (ts *TcpServer) withRecoverHandleRawConnAsync(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			bufStack := make([]byte, 1024*100)
			runtime.Stack(bufStack, false)
			bf := bytes.NewBuffer(bufStack)
			log.Println(bf.String())
		}
	}()
	defer func() {
		ts.wg.Done()
		ts.sem <- 1
	}()
	defer conn.Close()

	ts.wg.Add(1)

	if err := ts.handler.Handle(ts.ctx, conn); err != nil {
		log.Println(err)
	}
}
