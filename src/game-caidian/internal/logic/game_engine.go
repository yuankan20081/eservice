package logic

import (
	"game-util"
	"golang.org/x/net/context"
	"time"
)

type BankeringInfo interface {
	BankeringRequest() (region, name string, gold uint64)
	BankeringReplay(code int32)
	BankeringId() uint64
}

type BetInfo interface {
	BetRequest() (region, name string, pos int32, gold uint64)
	BetReply(code int32)
	BetterId() uint64
}

type GameEngineStatus int32

const (
	IsChoosingBanker GameEngineStatus = iota // 抢庄
	BankerChosed                             // 抽庄
	IsBetting                                // 下注
	BettingClosed                            // 下注结束
	Balancing                                // 结算
	Rewarding                                // 开奖
)

func (gs GameEngineStatus) String() string{
	switch gs{
	case IsChoosingBanker:
		return ""
	case BankerChosed:
		return ""
	case IsBetting:
		return ""
	case BettingClosed:
		return ""
	case Balancing:
		return ""
	case Rewarding:
		return ""
	default:
		return "unknown"
	}
}

type GameEngine struct {
	curStatus            GameEngineStatus
	statusChangedChannel chan GameEngineStatus
	bankerChannel        chan BankeringInfo
	betChannel           chan BetInfo
	lstBankering         map[int32]BankeringInfo
	lstBetting           map[int32]BetInfo
	results              [3]int32
}

func NewGameEngine() *GameEngine {
	return &GameEngine{
		curStatus:            IsChoosingBanker,
		statusChangedChannel: make(chan GameEngineStatus, 1),
		bankerChannel:        make(chan BankeringInfo, 1000),
		betChannel:           make(chan BetInfo, 1000),
		lstBankering:         make(map[int32]BankeringInfo),
		lstBetting:           make(map[int32]BetInfo),
	}
}

func (ge *GameEngine) Serve(ctx context.Context) error {
	// prepare
	ge.statusChangedChannel <- ge.curStatus

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case status := <-ge.statusChangedChannel:
			// TODO: broadcast status change

			switch status {
			case IsChoosingBanker:
				ge.beginEventChoosingBanker(ctx)
			case BankerChosed:
				ge.beginEventBankerClose(ctx)
			case IsBetting:
				ge.beginEventBetting(ctx)
			case BettingClosed:
				ge.beginEventBettingClose(ctx)
			case Balancing:
				ge.beginEventBalance(ctx)
			case Rewarding:
				ge.beginEventReward(ctx)
			}

		case info := <-ge.bankerChannel:
			if ge.curStatus != IsChoosingBanker {
				info.BankeringReplay(0)
				continue
			}

			// TODO: check satisfy server config

			// if all ok
			ge.lstBankering[info.BankeringId()] = info
		case info := <-ge.betChannel:
			if ge.curStatus != IsBetting {
				info.BetReply(0)
				continue
			}

			// TODO: check ambigous bet

			// if all ok
			ge.lstBetting[info.BetterId()] = info
		}
	}

	return nil
}

func (ge *GameEngine) beginEventChoosingBanker(ctx context.Context) error {
	game_util.Debug("开始抢庄")
	defer time.AfterFunc(time.Second*5, func() {
		ge.statusChangedChannel <- BankerChosed
	})

	// prepare results
	ge.prepareResults()
	return nil
}

func (ge *GameEngine) beginEventBankerClose(ctx context.Context) error {
	game_util.Debug("开始抽庄")
	// TODO: choose a banker or go back choosing
	if len(ge.lstBankering) == 0 {
		time.AfterFunc(time.Second*2, func() {
			ge.statusChangedChannel <- IsChoosingBanker
		})
	} else {

	}

	return nil
}

func (ge *GameEngine) beginEventBetting(ctx context.Context) error {
	game_util.Debug("开始下注")
	defer time.AfterFunc(time.Second*10, func() {
		ge.statusChangedChannel <- BettingClosed
	})

	return nil
}

func (ge *GameEngine) beginEventBettingClose(ctx context.Context) error {
	game_util.Debug("下注结束")
	defer time.AfterFunc(time.Second*2, func() {
		ge.statusChangedChannel <- Balancing
	})

	return nil
}

func (ge *GameEngine) beginEventBalance(ctx context.Context) error {
	game_util.Debug("开始结算")
	defer time.AfterFunc(time.Second*2, func() {
		ge.statusChangedChannel <- Rewarding
	})

	// TODO: broadcast result

	// TODO: calc reward

	return nil
}

func (ge *GameEngine) beginEventReward(ctx context.Context) error {
	game_util.Debug("开始发奖")
	defer time.AfterFunc(time.Second*2, func() {
		ge.statusChangedChannel <- IsChoosingBanker
	})

	return nil
}

func (ge *GameEngine) prepareResults() {
	for i, _ := range ge.results {
		ge.results[i] = int32(i)
	}
}
