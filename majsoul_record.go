package main

import (
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util"
	"strconv"
	"time"
	"sort"
)

type _majsoulRecordAccount struct {
	AccountID int `json:"account_id"`
	// 初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	Seat     int    `json:"seat"` // *重点是拿到自己的座位
	Nickname string `json:"nickname"`
}

// 牌谱基本信息
type majsoulRecordBaseInfo struct {
	UUID      string `json:"uuid"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`

	Accounts []_majsoulRecordAccount `json:"accounts"`
}

func (i *majsoulRecordBaseInfo) sort() {
	sort.Slice(i.Accounts, func(i_, j int) bool {
		return i.Accounts[i_].Seat < i.Accounts[j].Seat
	})
}

var seatNameZH = []string{"东", "南", "西", "北"}

func (i *majsoulRecordBaseInfo) String() string {
	i.sort()

	const timeFormat = "2006-01-02 15:04:05"
	output := fmt.Sprintf("%s\n从 %s\n到 %s\n\n", i.UUID, time.Unix(i.StartTime, 0).Format(timeFormat), time.Unix(i.EndTime, 0).Format(timeFormat))
	maxAccountID := 0
	for _, account := range i.Accounts {
		maxAccountID = util.MaxInt(maxAccountID, account.AccountID)
	}
	accountShownWidth := len(strconv.Itoa(maxAccountID))
	for _, account := range i.Accounts {
		output += fmt.Sprintf("%s %*d %s\n", seatNameZH[account.Seat], accountShownWidth, account.AccountID, account.Nickname)
	}
	return output
}

func (i *majsoulRecordBaseInfo) getSelfSeat(accountID int) (int, error) {
	if len(i.Accounts) == 0 {
		return -1, fmt.Errorf("牌谱基本信息为空")
	}
	if len(i.Accounts) == 3 {
		return -1, fmt.Errorf("暂不支持三人麻将")
	}
	for _, account := range i.Accounts {
		if account.AccountID == accountID {
			return account.Seat, nil
		}
	}
	return -1, fmt.Errorf("找不到用户 %d", accountID)
}

// 获取第一局的庄家：0=自家, 1=下家, 2=对家, 3=上家
func (i *majsoulRecordBaseInfo) getFistRoundDealer(accountID int) (firstRoundDealer int, err error) {
	selfSeat, err := i.getSelfSeat(accountID)
	if err != nil {
		return
	}
	const playerNumber = 4
	return (playerNumber - selfSeat) % playerNumber, nil
}

//

// 牌谱中的单个操作信息
type majsoulRecordAction struct {
	Name   string          `json:"name"`
	Action *majsoulMessage `json:"data"`
}

func parseMajsoulRecordAction(actions []*majsoulRecordAction) (roundActionsList [][]*majsoulRecordAction, err error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("数据异常：拿到的牌谱内容为空")
	}

	var currentRoundActions []*majsoulRecordAction
	for _, action := range actions {
		if action.Name == "RecordNewRound" {
			if len(currentRoundActions) > 0 {
				roundActionsList = append(roundActionsList, currentRoundActions)
			}
			currentRoundActions = []*majsoulRecordAction{action}
		} else {
			if len(currentRoundActions) == 0 {
				return nil, fmt.Errorf("数据异常：未收到 RecordNewRound")
			}
			currentRoundActions = append(currentRoundActions, action)
		}
	}
	if len(currentRoundActions) > 0 {
		roundActionsList = append(roundActionsList, currentRoundActions)
	}
	return
}
