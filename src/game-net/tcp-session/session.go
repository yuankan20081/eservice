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

var (
	ErrHeartbeatStop = errors.New("connection stops heartbeat")
	ErrIllegalPacket = errors.New("illegal packet")
)

type TcpSession struct {
	sendChannel      chan *buffer.Buffer
	recvChannel      chan *buffer.Buffer
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
	timeoutPoint     time.Time
}

func New(p *pool.Pool, bufferProcessor buffer.Handler) *TcpSession {
	return &TcpSession{
		sendChannel:      make(chan *buffer.Buffer, 1000),
		recvChannel:      make(chan *buffer.Buffer, 1000),
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
		timeoutPoint:     time.Now().Add(time.Second * 10),
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

	close(ts.recvChannel)
	if len(ts.recvChannel) > 0 {
		for pac := range ts.recvChannel {
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

	go withRecover(func() {
		ts.doReadPacAsync(ctx, conn)
	})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			{
				var counter int = 0

				// process errors first
				select {
				case err := <-ts.errorChannel:
					// if get an error, just break the loop
					// because somewhere is going wrong
					return err
				default:

				}

				// read write
				select {
				case pac := <-ts.sendChannel:
					select {
					case <-ts.readySendChannel:
						// can send now, so send it right now
						ts.listSend.PushBack(pac)
						go ts.doSendAsync(ctx, conn)
					default:
						ts.listPreSend.PushBack(pac)
					}
				default:
					counter += 1
				}

				select {
				case pac := <-ts.recvChannel:
					// reset heartbeat
					ts.timeoutPoint = time.Now().Add(time.Second * 10)

					select {
					case <-ts.readyProcChannel:
						// can process now, so process it right now
						ts.listProc.PushBack(pac)
						go ts.doProcessAsync(ctx)
					default:
						ts.listPreProc.PushBack(pac)
					}
				default:
					counter += 1
				}

				// send and process
				select {
				case <-ts.readySendChannel:
					// check pre send list
					if ts.listPreSend.Len() > 0 {
						ts.listSend.PushBackList(ts.listPreSend)
						ts.listPreSend = list.New()
					}

					if ts.listSend.Len() > 0 {
						go ts.doSendAsync(ctx, conn)
					}
				default:
					counter += 1
				}

				select {
				case <-ts.readyProcChannel:
					// check pre process list
					if ts.listPreProc.Len() > 0 {
						ts.listProc.PushBackList(ts.listPreProc)
						ts.listPreProc = list.New()
					}

					if ts.listProc.Len() > 0 {
						go ts.doProcessAsync(ctx)
					}
				default:
					counter += 1
				}

				// check heartbeat
				if time.Now().After(ts.timeoutPoint) {
					return ErrHeartbeatStop
				}

				// idle
				if counter == 4 {
					time.Sleep(time.Millisecond * 10)
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

func (ts *TcpSession) doReadPacAsync(ctx context.Context, r io.Reader) {
	var head game_net.PacketHead
	if err := binary.Read(r, binary.LittleEndian, &head); err != nil {
		ts.errorChannel <- err
		return
	}

	// validate head
	if head.PayloadLength > game_net.MaxPayloadLength {
		ts.errorChannel <- ErrIllegalPacket
		return
	}

	pac := ts.bufferPool.Get()

	if _, err := io.CopyN(pac, r, int64(head.PayloadLength)); err != nil {
		ts.bufferPool.Put(pac)
		ts.errorChannel <- err
		return
	} else {
		ts.recvChannel <- pac
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
