package agent

import "io"

type AgentBet struct {
	token string
	server string
	account string
	name string
	pos int32
	gold uint64
	w io.Writer
}

func (ab *AgentBet) BetRequest()(string, string, int32, uint64){
	return ab.server, ab.name, ab.pos, ab.gold
}

func (ab *AgentBet) BetReply(code int32){

}

func (ab *AgentBet) BetterId() uint64{

	return 0
}

func (ab *AgentBet) Win(gold uint64){

}
