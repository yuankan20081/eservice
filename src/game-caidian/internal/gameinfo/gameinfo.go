package gameinfo

type Winner interface {
	Win(gold uint64)
}

type BankeringInfo interface {
	BankeringRequest() (region, name string, gold uint64)
	BankeringReply(code int32)
	BankeringId() string
	BecomeBanker()
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
