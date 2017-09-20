package agent

import (
	"bytes"
	"crypto/md5"
	. "game-caidian/internal/gameinfo"
	"game-net/writer"
	"io"
)

type AgentBet struct {
	token   string
	server  string
	account string
	name    string
	pos     int32
	gold    uint64
	w       writer.Writer
	uid     uint64
}

func (ab *AgentBet) BetRequest() (string, string, int32, uint64) {
	return ab.server, ab.name, ab.pos, ab.gold
}

func (ab *AgentBet) BetReply(code int32) {
	var body = AgentOperateReply{
		Operate:  byte(ab.pos),
		Reserved: ab.uid,
		Err:      byte(code),
	}
	copy(body.Account[:], []byte(ab.account))
	copy(body.Name[:], []byte(ab.name))

	ab.w.WriteResponse(SmAgentOPerate, &body)
}

func (ab *AgentBet) BetterId() string {
	m := md5.New()
	io.WriteString(m, ab.token)
	io.WriteString(m, ab.server)
	io.WriteString(m, ab.account)
	io.WriteString(m, ab.name)
	b := bytes.NewBuffer(m.Sum(nil))
	return b.String()
}

func (ab *AgentBet) BetPos() int32 {
	return ab.pos
}

func (ab *AgentBet) UpdateFrom(other BetInfo) {
	_, _, _, gold := other.BetRequest()
	ab.gold += gold
}

func (ab *AgentBet) Win(gold uint64) {

}
