package agent

import (
	"encoding/binary"
	"bytes"
	"golang.org/x/net/context"
	"io"
	"game-util/publisher"
)

type Reader struct {
	token string
	server string
	name string
	authed bool
	pub *publisher.Publisher
}

func NewReader(pub *publisher.Publisher) *Reader {
	return &Reader{
		authed: false,
		pub: pub,
	}
}

func (h *Reader) Read(ctx context.Context, r io.Reader, w io.Writer) error {
	for{
		select{
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var head PacHead
		if err := binary.Read(r, binary.LittleEndian, &head); err!=nil{
			return err
		}

		buf := new(bytes.Buffer)
		if _, err := io.CopyN(buf, r, int64(head.PayloadLength)); err!=nil{
			return err			
		}

		if !h.authed{
			// TODO: do auth rpc

			// TODO: regist to publisher
			h.pub.Add("md5", w)
			defer h.pub.Remove("md5")
		}

		if err := h.executeBuffer(head.Proto, buf); err!=nil{
			return err
		}
	}	
}

func (h *Reader) executeBuffer(proto uint32, buf *bytes.Buffer) error{
	switch proto{

	}
	return nil
}