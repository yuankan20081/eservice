package agent

import (
	"bytes"
	"encoding/binary"
	"errors"
	"game-share/centerservice"
	"game-util/publisher"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"time"
)

var (
	errIllegalAgent = errors.New("not a valid agent")
	errAuthFailed   = errors.New("agent ticket invalid")
)

type Reader struct {
	token  string
	server string
	name   string
	authed bool
	pub    *publisher.Publisher
}

func NewReader(pub *publisher.Publisher) *Reader {
	return &Reader{
		authed: false,
		pub:    pub,
	}
}

func (h *Reader) Read(ctx context.Context, r io.Reader, w io.Writer) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var head PacHead
		if err := binary.Read(r, binary.LittleEndian, &head); err != nil {
			return err
		}

		p := make([]byte, head.PayloadLength)

		if _, err := io.ReadFull(r, p); err != nil {
			return err
		}

		DecBuff(p, DefaultEncKey)

		buf := bytes.NewBuffer(p)

		if !h.authed {
			if head.Proto != CmAgentAuth {
				return errIllegalAgent
			}
			var auth AgentAuthReq
			binary.Read(buf, binary.LittleEndian, &auth)

			// do auth rpc
			if token, server, success, err := h.doCenterAuth(ctx, "", string(auth.LicKey[:])); err != nil {
				return err
			} else if !success {
				// do reply fail
				h.authReply(w, 1)
				return errAuthFailed
			} else {
				// do replay success
				h.authReply(w, 0)
				h.token = token
				h.server = server
			}

			h.authed = true

			// regist to publisher
			h.pub.Add(w)
			defer h.pub.Remove(w)

			continue
		}

		if err := h.executeBuffer(head.Proto, buf, w); err != nil {
			return err
		}
	}
}

func (h *Reader) executeBuffer(proto uint16, buf *bytes.Buffer, w io.Writer) error {
	switch proto {
	case CmAgentOperate:
	default:

	}
	return nil
}

func (h *Reader) doCenterAuth(ctx context.Context, ip, ticket string) (string, string, bool, error) {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	cc, err := grpc.DialContext(ctx, "127.0.0.1:41000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return "", "", false, err
	}
	defer cc.Close()

	center := centerservice.NewCenterServiceClient(cc)

	var req = centerservice.AgentAuthRequest{
		Ticket: ticket,
		Ip:     ip,
	}

	reply, err := center.AgentAuth(ctx, &req)
	if err != nil {
		return "", "", false, err
	}

	return reply.Token, reply.Server, reply.Code == centerservice.AgentAuthReply_SUCCESS, nil
}

func (h *Reader) authReply(w io.Writer, code byte) {
	var body = AgentAuthReply{
		Err:    code,
		EncKey: DefaultEncKey,
	}

	var aw = AgentWriter{
		w:     w,
		proto: SmAgentAuth,
	}

	binary.Write(&aw, binary.LittleEndian, &body)
}

type AgentWriter struct {
	w     io.Writer
	proto uint16
}

func (aw *AgentWriter) Write(p []byte) (int, error) {
	EncBuff(p, DefaultEncKey)

	var head = PacHead{
		Tag:           MsgTag,
		Proto:         aw.proto,
		PayloadLength: uint16(len(p)) + HeadSize,
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, &head)
	buf.Write(p)

	if _, err := buf.WriteTo(aw.w); err != nil {
		return 0, err
	} else {
		return len(p), nil
	}
}
