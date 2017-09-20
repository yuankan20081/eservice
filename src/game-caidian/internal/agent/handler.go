package agent

import (
	"bytes"
	"encoding/binary"
	"errors"
	. "game-caidian/internal/gameinfo"
	"game-caidian/internal/logic"
	gamewriter "game-net/writer"
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
	errUnknownProto = errors.New("unknown proto")
)

type Reader struct {
	token  string
	server string
	name   string
	authed bool
	pub    *publisher.Publisher
	ge     *logic.GameEngine
}

func NewReader(pub *publisher.Publisher, ge *logic.GameEngine) *Reader {
	return &Reader{
		authed: false,
		pub:    pub,
		ge:     ge,
	}
}

func (h *Reader) Read(ctx context.Context, r io.Reader, w io.Writer) error {
	var rw = &responseWriter{
		Writer: w,
	}

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
			if err := h.doCenterAuth(ctx, "", string(auth.LicKey[:]), rw); err != nil {
				return err
			}

			h.authed = true

			// regist to publisher
			h.pub.Add(rw)
			defer h.pub.Remove(rw)

			continue
		}

		if err := h.executeBuffer(head.Proto, buf, rw); err != nil {
			return err
		}
	}
}

func (h *Reader) executeBuffer(proto uint16, buf *bytes.Buffer, rw gamewriter.Writer) error {
	switch proto {
	case CmAgentOperate:
		return h.doOperate(buf, rw)
	default:
		return errUnknownProto
	}
}

func (h *Reader) doOperate(buf *bytes.Buffer, rw gamewriter.Writer) error {
	var op AgentOperateReq

	if err := binary.Read(buf, binary.LittleEndian, &op); err != nil {
		return err
	}

	if op.Operate == 0 {
		// banker
		b := &AgentBankering{
			token:   h.token,
			server:  h.server,
			account: string(op.Account[:]),
			name:    string(op.Name[:]),
			gold:    op.OpGold,
			w:       rw,
			uid:     op.Reserved,
		}
		h.ge.AddBankering(b)
	} else {
		// bet
		b := &AgentBet{
			token:   h.token,
			server:  h.server,
			account: string(op.Account[:]),
			name:    string(op.Name[:]),
			gold:    op.OpGold,
			w:       rw,
			uid:     op.Reserved,
			pos:     int32(op.Pos),
		}
		h.ge.AddBet(b)
	}

	return nil
}

func (h *Reader) doCenterAuth(ctx context.Context, ip, ticket string, w gamewriter.Writer) error {
	ctx, _ = context.WithTimeout(ctx, time.Second*2)
	cc, err := grpc.DialContext(ctx, "127.0.0.1:41000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer cc.Close()

	center := centerservice.NewCenterServiceClient(cc)

	var req = centerservice.AgentAuthRequest{
		Ticket: ticket,
		Ip:     ip,
	}

	reply, err := center.AgentAuth(ctx, &req)
	if err != nil {
		return err
	}

	//return reply.Token, reply.Server, reply.Code == centerservice.AgentAuthReply_SUCCESS, nil
	if reply.Code != centerservice.AgentAuthReply_SUCCESS {
		w.WriteResponse(SmAgentAuth, &AgentAuthReply{
			Err:    1,
			EncKey: DefaultEncKey,
		})
		return errAuthFailed
	} else {
		h.token = reply.Token
		h.server = reply.Server
		w.WriteResponse(SmAgentAuth, &AgentAuthReply{
			Err:    0,
			EncKey: DefaultEncKey,
		})
		return nil
	}
}

type agentWriter struct {
	w     io.Writer
	proto uint16
}

func (aw *agentWriter) Write(p []byte) (int, error) {
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

type responseWriter struct {
	io.Writer
}

func (w *responseWriter) WriteResponse(proto uint16, body interface{}) {
	var aw = agentWriter{
		w:     w.Writer,
		proto: proto,
	}

	binary.Write(&aw, binary.LittleEndian, body)
}
