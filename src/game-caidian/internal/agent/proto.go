package agent

const (
	HeadSize      = 10
	MsgTag        = 0xfa1f2e3d
	DefaultEncKey = 911
)

const (
	CmAgentAuth uint16 = iota
	SmAgentAuth
	CmAgentOperate
	SmAgentOPerate
	SmGameStatus
	SmBroadcastDice
	SmAgentGift
	SmAgentBecomeBanker
	SmBroadcastWager
	SmAgentM2Broadcast
	SmAgentPredictResult
	SmGameConfig
)

type PacHead struct {
	Tag           uint32
	SvcType       uint8
	Proto         uint16
	EncType       uint8
	PayloadLength uint16
}

type AgentAuthReq struct {
	LicKey [38]byte
}

type AgentAuthReply struct {
	Err    byte
	EncKey uint32
}

type AgentOperateReq struct {
	Operate  byte
	Reserved uint64
	OpGold   uint64
	Account  [31]byte
	Name     [31]byte
	Pos      byte
}

type AgentOperateReply struct {
	Operate  byte
	Reserved uint64
	Account  [31]byte
	Name     [31]byte
	Err      byte
}

type GameStatusChanged struct {
	Step byte
	Stay byte
}

type BroadcastDice struct {
	Dice struct {
		DiceVal [3]byte
	}
}

type AgentGift struct {
	Account  [31]byte
	Name     [31]byte
	Reserved uint64
	Gold     uint64
}

type AgentBanker struct {
	Server   [31]byte
	Account  [31]byte
	Name     [31]byte
	Gold     uint64
	Reserved uint64
}

type AgentWagerShow struct {
	BigGold   uint64
	BigLim    uint64
	SmallGold uint64
	SmallLim  uint64
}
