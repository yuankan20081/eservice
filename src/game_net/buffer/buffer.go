package buffer

import (
	"bytes"
	"io"
)

type Buffer struct {
	proto  uint32
	buffer *bytes.Buffer
}

func New() *Buffer {
	return &Buffer{
		proto:  0,
		buffer: new(bytes.Buffer),
	}
}

func (b *Buffer) Reset() {
	b.buffer.Reset()
}

func (b *Buffer) SetProto(proto uint32) {
	b.proto = proto
}

func (b *Buffer) Write(p []byte) (int, error) {
	return b.buffer.Write(p)
}

func (b *Buffer) Read(p []byte) (int, error) {
	return b.buffer.Read(p)
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	return b.buffer.WriteTo(w)
}

func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	return b.buffer.ReadFrom(r)
}
