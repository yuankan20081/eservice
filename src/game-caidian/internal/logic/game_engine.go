package logic

import (
	"game-caidian/internal/config/observer"
	. "game-caidian/internal/gameinfo"
	"game-util"
	"game-util/publisher"
	"golang.org/x/net/context"
	"math/rand"
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
	curBankerUid         uint64
	curBankerAccount     string
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

func (ge *GameEngine) Serve(ctx context.Context, co *observer.Observer) error {
	// prepare
	ge.statusChangedChannel <- ge.curStatus

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case status := <-ge.statusChangedChannel:
			ge.curStatus = status
			// broadcast status change
			go ge.broadcastGameStatus(co)

			switch status {
			case IsChoosingBanker:
				ge.beginEventChoosingBanker(ctx, co)
			case BankerChosed:
				ge.beginEventBankerClose(ctx, co)
			case IsBetting:
				ge.beginEventBetting(ctx, co)
			case BettingClosed:
				ge.beginEventBettingClose(ctx, co)
			case Balancing:
				ge.beginEventBalance(ctx, co)
			case Rewarding:
				ge.beginEventReward(ctx, co)
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
					// check limit
					_, _, pos, gold := info.BetRequest()
					if pos == 1 && ge.curBigbetAvail < gold {
						// over limit
						info.BetReply(3)
					} else if pos == 2 && ge.curSmallbetAvail < gold {
						info.BetReply(3)
					} else {
						// update old from info
						old.UpdateFrom(info)
						info.BetReply(0)
						ge.updateBet(info)
					}
				}

				continue
			}

			// if fresh new info
			ge.lstBetting[info.BetterId()] = info
			info.BetReply(0)
			ge.updateBet(info)
		}
	}

	return nil
}

func (ge *GameEngine) beginEventChoosingBanker(ctx context.Context, co *observer.Observer) error {
	game_util.Debug("---当前阶段 %s, %d---", ge.curStatus, co.Config().Volatile.TimeBankering)
	defer time.AfterFunc(time.Second*co.Config().Volatile.TimeBankering, func() {
		ge.statusChangedChannel <- BankerChosed
	})

	// prepare results
	ge.prepareResults()

	ge.curTotalSmallbet = 0
	ge.curTotalBigbet = 0
	ge.curSmallbetAvail = 0
	ge.curBigbetAvail = 0

	return nil
}

func (ge *GameEngine) beginEventBankerClose(ctx context.Context, co *observer.Observer) error {
	game_util.Debug("---当前阶段 %s, %d---", ge.curStatus, co.Config().Volatile.TimeChooseBanker)

	enterNextStep := func() {
		ge.statusChangedChannel <- IsChoosingBanker
	}

	// TODO: choose a banker or go back choosing
	if len(ge.lstBankering) > 0 {
		var find BankeringInfo
		for _, info := range ge.lstBankering {
			if find == nil {
				find = info
			} else if info.BankMoreThan(find) {
				find = info
			}
		}

		if find != nil {
			ge.curBankerServer, ge.curBankerName, ge.curBankerGold = find.BankeringRequest()
			ge.curBigbetAvail = ge.curBankerGold
			ge.curSmallbetAvail = ge.curBankerGold

			// TODO: maybe should broadcast?
			ge.broadcastBanker()

			enterNextStep = func() {
				ge.statusChangedChannel <- IsBetting
			}
		}
	}

	time.AfterFunc(time.Second*co.Config().Volatile.TimeChooseBanker, enterNextStep)

	return nil
}

func (ge *GameEngine) beginEventBetting(ctx context.Context, co *observer.Observer) error {
	game_util.Debug("---当前阶段 %s, %d---", ge.curStatus, co.Config().Volatile.TimeBet)
	defer time.AfterFunc(time.Second*co.Config().Volatile.TimeBet, func() {
		ge.statusChangedChannel <- BettingClosed
	})

	return nil
}

func (ge *GameEngine) beginEventBettingClose(ctx context.Context, co *observer.Observer) error {
	game_util.Debug("---当前阶段 %s, %d---", ge.curStatus, co.Config().Volatile.TimeCloseBet)
	defer time.AfterFunc(time.Second*co.Config().Volatile.TimeCloseBet, func() {
		ge.statusChangedChannel <- Balancing
	})

	return nil
}

func (ge *GameEngine) beginEventBalance(ctx context.Context, co *observer.Observer) error {
	game_util.Debug("---当前阶段 %s, %d---", ge.curStatus, co.Config().Volatile.TimeBalance)
	defer time.AfterFunc(time.Second*co.Config().Volatile.TimeBalance, func() {
		ge.statusChangedChannel <- Rewarding
	})

	// broadcast result
	ge.broadcastDice()

	return nil
}

func (ge *GameEngine) beginEventReward(ctx context.Context, co *observer.Observer) error {
	game_util.Debug("---当前阶段 %s, %d---", ge.curStatus, co.Config().Volatile.TimeReward)
	// TODO: calc reward
	if ge.results.BankerWin() {

	} else if ge.results.BigbetWin() {

	} else {

	}

	defer time.AfterFunc(time.Second*co.Config().Volatile.TimeReward, func() {
		// TODO: do some cleanup
		ge.lstBetting = make(map[string]BetInfo)
		ge.lstBankering = make(map[string]BankeringInfo)

		ge.statusChangedChannel <- IsChoosingBanker
	})

	return nil
}

func (ge *GameEngine) prepareResults() {
	rand.Seed(time.Now().Unix())

	for i, _ := range ge.results {
		ge.results[i] = byte(rand.Intn(6) + 1)
	}

	game_util.Debug("%v", ge.results)
}

func (ge *GameEngine) broadcastGameStatus(co *observer.Observer) {
	var stay time.Duration

	switch ge.curStatus {
	case IsChoosingBanker:
		stay = co.Config().Volatile.TimeBankering
	case BankerChosed:
		stay = co.Config().Volatile.TimeChooseBanker
	case IsBetting:
		stay = co.Config().Volatile.TimeBet
	case BettingClosed:
		stay = co.Config().Volatile.TimeCloseBet
	case Balancing:
		stay = co.Config().Volatile.TimeBalance
	case Rewarding:
		stay = co.Config().Volatile.TimeReward
	}

	var body = GameStatusChanged{
		Step: byte(ge.curStatus) + 1,
		Stay: byte(stay),
	}

	ge.pub.Publish(SmGameStatus, &body)
}

func (ge *GameEngine) broadcastDice() {
	var body BroadcastDice
	copy(body.Dice.DiceVal[:], ge.results[:])
	if ge.results[0] == ge.results[1] && ge.results[1] == ge.results[2] {
		body.Dice.Result = 2
	} else if ge.results[0]+ge.results[1]+ge.results[2] < 11 {
		body.Dice.Result = 0 //xiao
	} else {
		body.Dice.Result = 1 //da
	}

	ge.pub.Publish(SmBroadcastDice, &body)
}

func (ge *GameEngine) broadcastWager() {
	var body = AgentWagerShow{
		BigGold:   ge.curTotalBigbet,
		BigLim:    ge.curBigbetAvail,
		SmallGold: ge.curTotalSmallbet,
		SmallLim:  ge.curSmallbetAvail,
	}

	ge.pub.Publish(SmBroadcastWager, &body)
}

func (ge *GameEngine) broadcastBanker() {
	var body = AgentBanker{
		Gold:     ge.curBankerGold,
		Reserved: ge.curBankerUid,
	}
	copy(body.Server[:], []byte(ge.curBankerServer))
	copy(body.Account[:], []byte(ge.curBankerAccount))
	copy(body.Name[:], []byte(ge.curBankerName))

	ge.pub.Publish(SmAgentBecomeBanker, &body)
}

func (ge *GameEngine) updateBet(info BetInfo) {
	_, _, _, gold := info.BetRequest()
	if info.BetPos() == 1 {
		ge.curTotalBigbet += gold
		ge.curBigbetAvail -= gold
		ge.curSmallbetAvail += gold
	} else {
		ge.curTotalSmallbet += gold
		ge.curSmallbetAvail -= gold
		ge.curBigbetAvail += gold
	}

	ge.broadcastWager()
}
