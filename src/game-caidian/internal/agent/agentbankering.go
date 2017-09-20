package agent

import (
	"bytes"
	"crypto/md5"
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

}

func (ab *AgentBankering) BankerId() string {
	m := md5.New()
	io.WriteString(m, ab.token)
	io.WriteString(m, ab.server)
	io.WriteString(m, ab.account)
	io.WriteString(m, ab.name)
	b := bytes.NewBuffer(m.Sum(nil))
	return b.String()
}

func (ab *AgentBankering) BecomeBanker() {

}

func (ab *AgentBankering) Win(gold uint64) {

}
