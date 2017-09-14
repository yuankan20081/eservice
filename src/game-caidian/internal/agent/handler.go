package agent

import (
	"game-net/buffer"
	"golang.org/x/net/context"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, buf *buffer.Buffer) error {
	switch buf.Proto() {
	default:

	}

	return nil
}
