package main

import (
	"strings"
	"fmt"
	"os"
	"time"
	"github.com/fatih/color"
	"math/rand"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// go build -ldflags "-X main.version=$(git describe --abbrev=0 --tags)" -o mahjong-helper
var version = "dev"

var (
	showImproveDetail      bool
	showAgariAboveShanten1 bool
	showScore              bool
	showAllYakuTypes       bool
)

func welcome() int {
	platforms := map[int]string{
		0: "天凤",
		1: "雀魂",
	}

	fmt.Println("使用说明：https://github.com/EndlessCheng/mahjong-helper")
	fmt.Println("问题反馈：https://github.com/EndlessCheng/mahjong-helper/issues")
	fmt.Println("吐槽群：375865038")

	fmt.Println()

	fmt.Println("请输入数字，以选择对应的平台：")
	for k, v := range platforms {
		fmt.Printf("%d - %s\n", k, v)
	}

	choose := 1
	fmt.Scanf("%d", &choose)
	if choose < 0 || choose > 1 {
		choose = 1
	}

	clearConsole()
	platformName := platforms[choose]
	if choose == 1 {
		platformName += "（水晶杠杠版）"
	}
	color.HiGreen("已选择 - %s", platformName)
	if choose == 1 {
		color.HiYellow("提醒：若您已登录游戏，请刷新网页，或者开启一局人机对战\n" +
			"该步骤用于获取您的账号 ID，便于在游戏开始时分析自风，否则程序将无法解析后续数据")
	}

	return choose
}

func main() {
	color.HiGreen("日本麻将助手 %s (by EndlessCheng)", version)
	if version != "dev" {
		go alertNewVersion(version)
	}

	flags, restArgs := parseArgs(os.Args[1:])

	isMajsoul := flags.Bool("majsoul")
	isTenhou := flags.Bool("tenhou")
	isAnalysis := flags.Bool("analysis")
	isInteractive := flags.Bool("i", "interactive")
	showImproveDetail = flags.Bool("detail")
	showAgariAboveShanten1 = flags.Bool("a", "agari")
	showScore = flags.Bool("s", "score")
	showAllYakuTypes = flags.Bool("y", "yaku")

	humanDoraTiles := flags.String("d", "dora")
	humanTiles := strings.Join(restArgs, " ")
	humanTilesInfo := &model.HumanTilesInfo{
		HumanTiles:     humanTiles,
		HumanDoraTiles: humanDoraTiles,
	}

	switch {
	case isMajsoul:
		runServer(true)
	case isTenhou || isAnalysis:
		runServer(false)
	case isInteractive:
		// 交互模式
		if err := interact(humanTilesInfo); err != nil {
			errorExit(err)
		}
	case len(restArgs) > 0:
		// 静态分析
		if _, err := analysisHumanTiles(humanTilesInfo); err != nil {
			errorExit(err)
		}
	default:
		// 服务器模式
		isHTTPS := welcome() == 1
		runServer(isHTTPS)
	}
}
