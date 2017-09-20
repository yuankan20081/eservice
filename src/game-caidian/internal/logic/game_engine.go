package logic

import (
	. "game-caidian/internal/gameinfo"
	"game-util"
	"game-util/publisher"
	"golang.org/x/net/context"
	"time"
)

type GameEngineStatus int32

const (
	IsChoosingBanker GameEngineStatus = iota // 抢庄
	BankerChosed                             // 抽庄
	IsBetting                                // 下注
	BettingClosed                            // 下注结束
	Balancing                                // 结算
	Rewarding                                // 开奖
)

func (gs GameEngineStatus) String() string {
	switch gs {
	case IsChoosingBanker:
		return "正在抢庄"
	case BankerChosed:
		return "正在抽庄"
	case IsBetting:
		return "正在下注"
	case BettingClosed:
		return "下注结束"
	case Balancing:
		return "正在结算"
	case Rewarding:
		return "正在开奖"
	default:
		return "unknown"
	}
}

type GameResult [3]byte

func (gr *GameResult) BankerWin() bool {

	return false
}

func (gr *GameResult) BigbetWin() bool {

	return false
}

type GameEngine struct {
	curStatus            GameEngineStatus
	statusChangedChannel chan GameEngineStatus
	bankerChannel        chan BankeringInfo
	betChannel           chan BetInfo
	lstBankering         map[string]BankeringInfo
	lstBetting           map[string]BetInfo
	results              GameResult
	pub                  *publisher.Publisher
	curBankerServer      string
	curBankerName        string
	curBankerGold        uint64
	curTotalBigbet       uint64
	curTotalSmallbet     uint64
	curBigbetAvail       uint64
	curSmallbetAvail     uint64
}

func NewGameEngine(pub *publisher.Publisher) *GameEngine {
	return &GameEngine{
		curStatus:            IsChoosingBanker,
		statusChangedChannel: make(chan GameEngineStatus, 1),
		bankerChannel:        make(chan BankeringInfo, 1000),
		betChannel:           make(chan BetInfo, 1000),
		lstBankering:         make(map[string]BankeringInfo),
		lstBetting:           make(map[string]BetInfo),
		pub:                  pub,
	}
}

func (ge *GameEngine) AddBankering(info BankeringInfo) {
	ge.bankerChannel <- info
}

func (ge *GameEngine) AddBet(info BetInfo) {
	ge.betChannel <- info
}

func (ge *GameEngine) Serve(ctx context.Context) error {
	// prepare
	ge.statusChangedChannel <- ge.curStatus

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case status := <-ge.statusChangedChannel:
			ge.curStatus = status
			// broadcast status change
			go ge.broadcastGameStatus()

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
				// wrong timing
				info.BankeringReply(1)
				continue
			}

			id := info.BankeringId()

			if _, ok := ge.lstBankering[id]; ok {
				// already in queue
				info.BankeringReply(2)
				continue
			}

			// TODO: check satisfy server config

			// if all ok
			ge.lstBankering[id] = info
			info.BankeringReply(0)
		case info := <-ge.betChannel:
			// wrong timing
			if ge.curStatus != IsBetting {
				info.BetReply(1)
				continue
			}

			id := info.BetterId()

			if old, ok := ge.lstBetting[id]; ok {
				if old.BetPos() != info.BetPos() {
					// ambigous bet
					info.BetReply(2)
				} else {
					// update old from info
					old.UpdateFrom(info)
					info.BetReply(0)
				}

				continue
			}

			// if fresh new info
			ge.lstBetting[info.BetterId()] = info
			info.BetReply(0)
		}
	}

	return nil
}

func (ge *GameEngine) beginEventChoosingBanker(ctx context.Context) error {
	game_util.Debug("---当前阶段 %s---", ge.curStatus)
	defer time.AfterFunc(time.Second*5, func() {
		ge.statusChangedChannel <- BankerChosed
	})

	// prepare results
	ge.prepareResults()
	return nil
}

func (ge *GameEngine) beginEventBankerClose(ctx context.Context) error {
	game_util.Debug("---当前阶段 %s---", ge.curStatus)
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
	game_util.Debug("---当前阶段 %s---", ge.curStatus)
	defer time.AfterFunc(time.Second*10, func() {
		ge.statusChangedChannel <- BettingClosed
	})

	return nil
}

func (ge *GameEngine) beginEventBettingClose(ctx context.Context) error {
	game_util.Debug("---当前阶段 %s---", ge.curStatus)
	defer time.AfterFunc(time.Second*2, func() {
		ge.statusChangedChannel <- Balancing
	})

	return nil
}

func (ge *GameEngine) beginEventBalance(ctx context.Context) error {
	game_util.Debug("---当前阶段 %s---", ge.curStatus)
	defer time.AfterFunc(time.Second*2, func() {
		ge.statusChangedChannel <- Rewarding
	})

	// broadcast result
	ge.broadcastDice()

	return nil
}

func (ge *GameEngine) beginEventReward(ctx context.Context) error {
	game_util.Debug("---当前阶段 %s---", ge.curStatus)
	// TODO: calc reward
	if ge.results.BankerWin() {

	} else if ge.results.BigbetWin() {

	} else {

	}

	defer time.AfterFunc(time.Second*2, func() {
		// TODO: do some cleanup
		ge.lstBetting = make(map[string]BetInfo)
		ge.lstBankering = make(map[string]BankeringInfo)

		ge.statusChangedChannel <- IsChoosingBanker
	})

	return nil
}

func (ge *GameEngine) prepareResults() {
	for i, _ := range ge.results {
		ge.results[i] = byte(i)
	}
}

func (ge *GameEngine) broadcastGameStatus() {
	var body = GameStatusChanged{
		Step: byte(ge.curStatus),
		Stay: 5,
	}

	ge.pub.Publish(SmGameStatus, &body)
}

func (ge *GameEngine) broadcastDice() {
	var body BroadcastDice
	copy(body.Dice.DiceVal[:], ge.results[:])

	ge.pub.Publish(SmBroadcastDice, &body)
}
