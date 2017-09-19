package agent

import "io"

type AgentBankering struct {
	token   string
	server  string
	account string
	name    string
	gold    uint64
	w       io.Writer
}

func (ab *AgentBankering) BankeringRequest() (string, string, uint64) {
	return ab.server, ab.name, ab.gold
}

func (ab *AgentBankering) BankeringReply(code int32) {

}

func (ab *AgentBankering) BankerId() uint64 {

	return 0
}

func (ab *AgentBankering) BecomeBanker() {

}

func (ab *AgentBankering) Win(gold uint64) {

}
