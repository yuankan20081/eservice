package tcp_session

import (
	"container/list"
	"encoding/binary"
	"errors"
	"game-net"
	"game-net/buffer"
	"game-net/buffer/pool"
	"golang.org/x/net/context"
	"io"
	"log"
	"net"
	"runtime"
	"time"
)

type TcpSession struct {
	sendChannel      chan *buffer.Buffer
	readyProcChannel chan bool
	readySendChannel chan bool
	listProc         *list.List
	listPreProc      *list.List
	listSend         *list.List
	listPreSend      *list.List
	bufferPool       *pool.Pool
	errorChannel     chan error
	bufferProcessor  buffer.Handler
	cleanupChannel   chan bool
}

func New(p *pool.Pool, bufferProcessor buffer.Handler) *TcpSession {
	return &TcpSession{
		sendChannel:      make(chan *buffer.Buffer, 1000),
		readyProcChannel: make(chan bool, 1),
		readySendChannel: make(chan bool, 1),
		listProc:         list.New(),
		listPreProc:      list.New(),
		listSend:         list.New(),
		listPreSend:      list.New(),
		bufferPool:       p,
		errorChannel:     make(chan error, 1),
		bufferProcessor:  bufferProcessor,
		cleanupChannel:   make(chan bool, 1),
	}
}

func (ts *TcpSession) Cleanup() {
	<-ts.cleanupChannel

	close(ts.sendChannel)
	if len(ts.sendChannel) > 0 {
		for pac := range ts.sendChannel {
			ts.bufferPool.Put(pac)
		}
	}

	for ts.listProc.Len() > 0 {
		pac := ts.listProc.Remove(ts.listProc.Front()).(*buffer.Buffer)
		ts.bufferPool.Put(pac)
	}

	for ts.listPreProc.Len() > 0 {
		pac := ts.listPreProc.Remove(ts.listPreProc.Front()).(*buffer.Buffer)
		ts.bufferPool.Put(pac)
	}

	for ts.listSend.Len() > 0 {
		pac := ts.listSend.Remove(ts.listSend.Front()).(*buffer.Buffer)
		ts.bufferPool.Put(pac)
	}

	for ts.listPreSend.Len() > 0 {
		pac := ts.listPreSend.Remove(ts.listPreSend.Front()).(*buffer.Buffer)
		ts.bufferPool.Put(pac)
	}
}

func (ts *TcpSession) Handle(ctx context.Context, conn net.Conn) error {
	defer func() {
		ts.cleanupChannel <- true
	}()

	// prepare
	ts.readySendChannel <- true
	ts.readyProcChannel <- true

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			{
				select {
				case <-ts.readySendChannel:
					// check pre send list
					if ts.listPreSend.Len() > 0 {
						ts.listSend.PushBackList(ts.listPreSend)
						ts.listPreSend = list.New()
					}

					if ts.listSend.Len() > 0 {
						go ts.doSendAsync(ctx, conn)
						continue
					} else {
						ts.readySendChannel <- true
						// go on do some read
					}
				case <-ts.readyProcChannel:
					// check pre process list
					if ts.listPreProc.Len() > 0 {
						ts.listProc.PushBackList(ts.listPreProc)
						ts.listPreProc = list.New()
					}

					if ts.listProc.Len() > 0 {
						go ts.doProcessAsync(ctx)
						continue
					} else {
						ts.readyProcChannel <- true
						// go on do some read
					}
				case pac := <-ts.sendChannel:
					select {
					case <-ts.readySendChannel:
						// can send now, so send it right now
						ts.listSend.PushBack(pac)
						go ts.doSendAsync(ctx, conn)
					default:
						ts.listProc.PushBack(pac)
					}
				case err := <-ts.errorChannel:
					// if get an error, just break the loop
					// because somewhere is going wrong
					return err
				default:

				}

				// TODO: read and to process list
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				pac, err := ts.readPac(ctx, conn)
				if err != nil {
					return err
				}
				select {
				case <-ts.readyProcChannel:
					// can process now, process right now
					ts.listProc.PushBack(pac)
					go ts.doProcessAsync(ctx)
				default:
					ts.listPreProc.PushBack(pac)
				}
			}
		}
	}
}

func (ts *TcpSession) doProcessAsync(ctx context.Context) {
	defer func() {
		ts.readyProcChannel <- true
	}()

loop:
	for ts.listProc.Len() > 0 {
		pac := ts.listProc.Remove(ts.listProc.Front()).(*buffer.Buffer)
		// process
		ts.bufferProcessor.Handle(ctx, pac)

		ts.bufferPool.Put(pac)

		select {
		case <-ctx.Done():
			break loop
		default:

		}
	}
}

func (ts *TcpSession) doSendAsync(ctx context.Context, w io.Writer) {
	defer func() {
		ts.readySendChannel <- true
	}()

loop:
	for ts.listSend.Len() > 0 {
		pac := ts.listSend.Remove(ts.listSend.Front()).(*buffer.Buffer)

		io.Copy(w, pac)

		ts.bufferPool.Put(pac)

		select {
		case <-ctx.Done():
			break loop
		default:

		}
	}
}

func (ts *TcpSession) readPac(ctx context.Context, r io.Reader) (*buffer.Buffer, error) {
	var head game_net.PacketHead
	if err := binary.Read(r, binary.LittleEndian, &head); err != nil {
		return nil, err
	}

	// validate head
	if head.PayloadLength > game_net.MaxPayloadLength {
		return nil, errors.New("invalid packet head, payload length too big")
	}

	pac := ts.bufferPool.Get()

	if _, err := io.CopyN(pac, r, int64(head.PayloadLength)); err != nil {
		ts.bufferPool.Put(pac)
		return nil, err
	} else {
		return pac, nil
	}
}

func withRecover(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			bufStack := make([]byte, 1024*1000)
			runtime.Stack(bufStack, false)
			log.Println(string(bufStack))
		}
	}()

	fn()
}
