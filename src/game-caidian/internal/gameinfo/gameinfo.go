package gameinfo

type Winner interface {
	Win(gold uint64)
}

type BankeringInfo interface {
	BankeringRequest() (region, name string, gold uint64)
	BankeringReplay(code int32)
	BankeringId() string
	BecomBanker()
	Winner
}

type BetInfo interface {
	BetRequest() (region, name string, pos int32, gold uint64)
	BetReply(code int32)
	BetterId() string
	BetPos() int32
	UpdateFrom(other BetInfo)
	Winner
}