package tcp_session

import (
	"log"
	"runtime/debug"
	"sync"
	"golang.org/x/net/context"
	"net"
	"io"
)

type TcpSessionReader interface{
	Read(ctx context.Context, r io.Reader, w io.Writer) error
}

type TcpSessionReadFunc func(ctx context.Context, r io.Reader, w io.Writer) error

func (fn TcpSessionReadFunc) Read(ctx context.Context, r io.Reader, w io.Writer)error{
	return fn(ctx, r, w)
}

type TcpSession struct {
	wg sync.WaitGroup
	errorChannel chan error
	sendChannel chan []byte
	CustomReader TcpSessionReader
	closedChannel chan bool
}

func New(cb TcpSessionReader) *TcpSession {
	return &TcpSession{
		errorChannel: make(chan error, 1),
		sendChannel: make(chan []byte, 1000),
		CustomReader: cb,
		closedChannel: make(chan bool),
	}
}

func (ts *TcpSession) Serve(ctx context.Context, conn net.Conn) error{
	ctx, cancelCtx := context.WithCancel(ctx)
	
	ts.wg.Add(2)

	go ts.doRead(ctx, conn)
	go ts.doWrite(ctx, conn)


	defer func(){
		close(ts.closedChannel)
		ts.wg.Wait()		
	}()

	select{
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ts.errorChannel:
		cancelCtx()
		return err
	}
}

func (ts *TcpSession) Write(p []byte) (int, error){
	select{
	case <-ts.closedChannel:
	default:
		ts.sendChannel <- p
	}	
	return len(p), nil
}

func (ts *TcpSession) doRead(ctx context.Context, r io.Reader){
	defer ts.wg.Done()

	if err := ts.CustomReader.Read(ctx, r, ts); err!=nil{
		ts.errorChannel <- err
	}

}

func (ts *TcpSession) doWrite(ctx context.Context, w io.Writer){
	defer ts.wg.Done()
	
	for{
		select{
		case <-ctx.Done():
			return
		case msg := <-ts.sendChannel:
			if _, err := w.Write(msg); err!=nil{
				ts.doError(err)
				return
			}
		}
	}
}

func (ts *TcpSession) doError(err error){
	select{
	case ts.errorChannel<-err:
	default:
	}
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
