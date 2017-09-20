package agent

import (
	"bytes"
	"crypto/md5"
	. "game-caidian/internal/gameinfo"
	"game-net/writer"
	"io"
)

type AgentBankering struct {
	token   string
	server  string
	account string
	name    string
	gold    uint64
	w       writer.Writer
	uid     uint64
}

func (ab *AgentBankering) BankeringRequest() (string, string, uint64) {
	return ab.server, ab.name, ab.gold
}

func (ab *AgentBankering) BankeringReply(code int32) {
	var body = AgentOperateReply{
		Operate:  0,
		Reserved: ab.uid,
		Err:      byte(code),
	}
	copy(body.Account[:], []byte(ab.account))
	copy(body.Name[:], []byte(ab.name))

	ab.w.WriteResponse(SmAgentOPerate, &body)
}

func (ab *AgentBankering) BankeringId() string {
	m := md5.New()
	io.WriteString(m, ab.token)
	io.WriteString(m, ab.server)
	io.WriteString(m, ab.account)
	io.WriteString(m, ab.name)
	b := bytes.NewBuffer(m.Sum(nil))
	return b.String()
}

func (ab *AgentBankering) BecomeBanker() {
	var body = AgentBanker{
		Gold:     ab.gold,
		Reserved: ab.uid,
	}
	copy(body.Server[:], []byte(ab.server))
	copy(body.Account[:], []byte(ab.account))
	copy(body.Name[:], []byte(ab.name))

	ab.w.WriteResponse(SmAgentBecomeBanker, &body)
}

func (ab *AgentBankering) Win(gold uint64) {
	var body = AgentGift{
		Reserved: ab.uid,
		Gold:     gold,
	}
	copy(body.Account[:], []byte(ab.account))
	copy(body.Name[:], []byte(ab.name))

	ab.w.WriteResponse(SmAgentGift, &body)
}
