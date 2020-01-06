package richman

import (
	"context"
	"github.com/ghaoo/rboot"
)

func setup(ctx context.Context, bot *rboot.Robot) (msg []rboot.Message) {
	in := ctx.Value("input").(rboot.Message)

	if in.Mate["SendByMySelf"] != nil && in.Mate["SendByMySelf"].(bool) {
		return nil
	}

	switch bot.Ruleset {
	case `start`:
		return StartGame(in, bot)
	case `shake`, `roll`:
		return Dice(in, bot)
	case `create1`, `create2`:
		return CreateGameRoom(in, bot)
	case `join`:
		return JoinRoom(in, bot)
	case `leave`:
		return AFK(in, bot)
	case `quit`:
		return QuitGame(in, bot)
	case `look`:
		return Look(in, bot)
	case `status`:
		return RoomInfo(in, bot)
	case `stop`:
		return StopGame(in)
	}

	return
}

func init() {
	// 创建游戏
	Go()

	if GM.Map != nil {
		rboot.RegisterScripts(`richman`, rboot.Script{
			Action: setup,
			Ruleset: map[string]string{
				`start`:   `^(?:开始游戏|游戏开始)$`,
				`shake`:   `.*[&lt;]{1}gameext type="2" content="(\d+)" [&gt;]{1}`,
				`roll`:    `^#roll#|掷骰子`,
				`create1`: `^创建游戏`,
				`create2`: `^创建(\d+)人游戏`,
				`join`:    `^加入游戏`,
				`leave`:   `^(?:AFK|托管|暂离|暂时离开)`,
				`quit`:    `^退出游戏$`,
				`status`:  `^游戏状态$`,
				`look`:    `^查看$`,
				`stop`:    `^结束游戏`,
			},
			Usage: "> `创建<N人>游戏`: 创建游戏，支持人数限制 \n" +
				"> `加入游戏`: 加入未开始的游戏 \n" +
				"> `开始游戏`: 开启游戏命令 \n" +
				"> `#roll#`: 掷骰子，微信网页版可直接在微信app中使用骰子表情" +
				"> `托管`|`暂离`: 由系统托管游戏 \n" +
				"> `退出游戏`: 退出游戏 \n" +
				"> `游戏状态`: 查看全部玩家资产状态 \n" +
				"> `查看`: 查看自身资产状态 \n" +
				"> `结束游戏`: 结束游戏",
			Description: "模拟大富翁游戏，暂时只支持 rboot 微信",
		})
	}
}
