package agent

import (
	"game-net/buffer"
	"golang.org/x/net/context"
	"io"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, buf *buffer.Buffer, w io.Writer) error {
	switch buf.Proto() {
	default:

	}

	return nil
}
