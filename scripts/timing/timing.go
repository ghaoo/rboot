package timing

import (
	"context"
	"fmt"
	"github.com/ghaoo/rboot"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type timer struct {
	t   *time.Timer
	end time.Time
	cmd string
}

// 定时器
var timerN = 0
var timers = make(map[int]*timer)

func setup(ctx context.Context, bot *rboot.Robot) []rboot.Message {
	in := ctx.Value("input").(rboot.Message)

	switch bot.Ruleset {
	case "timer":
		return start_timer(in, bot)
	case "stop_timer":
		return stop_timer(bot)
	case "status_timer":
		return status_timer(bot)
	case "ticker":
		return start_ticker(in, bot)
	case "stop_ticker":
		return stop_ticker(bot)
	case "status_ticker":
		return status_ticker(bot)
	}

	return nil
}

// 定时器开始定时
func start_timer(in rboot.Message, bot *rboot.Robot) []rboot.Message {

	args := bot.Args

	// 时间
	t, err := strconv.Atoi(args[1])
	if err != nil {
		logrus.Error(err)
		return nil
	}

	// 时间刻度
	arc := args[2]

	// 将时间转换为 time.Duration 类型
	tD, err := toDuration(t, arc)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	// 脚本名称
	script := args[3]

	// 检测脚本是否可执行
	_, err = rboot.DirectiveScript(script)
	if err != nil || script == "timing" {
		return []rboot.Message{{Content: "END"}}
	}

	// 脚本内命令名称
	cmd := args[4]

	myTimer := &timer{t: time.NewTimer(tD), end: time.Now().Add(tD), cmd: script + "." + cmd}

	n := timerN
	timers[timerN] = myTimer
	timerN += 1

	bot.SendText("定时器序号 "+strconv.Itoa(n), in.From)

	for {
		select {
		case <-myTimer.t.C:

			//bot.SendText(fmt.Sprintf("定时器 %d 执行脚本命令 %s.%s", n, script, cmd), in.From)

			delete(timers, n)

			sf, err := rboot.DirectiveScript(script)
			if err != nil {
				return []rboot.Message{{Content: fmt.Sprintf("定时器 %d 执行 %s.%s 命令时发生错误: %v", n, script, cmd, err)}}
			}

			bot.Ruleset = cmd

			return sf(nil, bot)
		}
	}
}

// 定时器结束定时
func stop_timer(bot *rboot.Robot) []rboot.Message {
	args := bot.Args

	tNS := args[1]
	tNum, err := strconv.Atoi(tNS)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if mt, ok := timers[tNum]; ok {
		if !mt.t.Stop() {
			<-mt.t.C
		}

		delete(timers, tNum)

		return []rboot.Message{{Content: "定时器 " + tNS + " 关闭"}}
	} else {
		return []rboot.Message{{Content: "未找到序号为 " + tNS + " 的定时器"}}
	}
}

// 定时器状态
func status_timer(bot *rboot.Robot) []rboot.Message {

	content := ""
	if len(timers) <= 0 {
		content = "没有正在进行的定时器"
	} else {
		content = "定时器状态: \n"

		for k, t := range timers {
			left := t.end.Sub(time.Now()).Seconds()
			content += fmt.Sprintf("定时器 %d， 命令 %s 剩余时间 %s \n", k, t.cmd, toCTime(int64(left)))
		}
	}

	return []rboot.Message{{Content: content}}
}

// Ticker
type ticker struct {
	n    int          // 已经执行的次数
	t    *time.Ticker // 续断器实体
	next time.Time    // 下一次执行时间
	cmd  string       // 脚本命令
}

var tickerN = 0
var tickers = make(map[int]*ticker)

// ticker开始计时
func start_ticker(in rboot.Message, bot *rboot.Robot) []rboot.Message {
	args := bot.Args

	// 时间
	t, err := strconv.Atoi(args[1])
	if err != nil {
		logrus.Error(err)
		return nil
	}

	// 时间刻度
	arc := args[2]

	// 将时间转换为 time.Duration 类型
	tD, err := toDuration(t, arc)
	if err != nil {
		logrus.Errorf("解析时间失败: %v", err)
		return nil
	}

	// 脚本名称
	script := args[3]

	// 检测脚本是否可执行
	_, err = rboot.DirectiveScript(script)
	if err != nil || script == "timing" {
		return []rboot.Message{{Content: "END"}}
	}

	// 脚本内命令名称
	cmd := args[4]

	myTicker := &ticker{n: 0, t: time.NewTicker(tD), next: time.Now().Add(tD), cmd: script + "." + cmd}

	n := tickerN
	tickers[tickerN] = myTicker
	tickerN += 1

	bot.SendText("续断器序号 "+strconv.Itoa(n), in.From)

	for {
		select {
		case <-myTicker.t.C:

			//bot.SendText(fmt.Sprintf("序号 %d 续断器循环执行脚本命令 %s.%s", n, script, cmd), in.From)

			myTicker.n += 1

			myTicker.next = time.Now().Add(tD)

			sf, err := rboot.DirectiveScript(script)
			if err != nil {
				myTicker.t.Stop()
				delete(tickers, n)
				return []rboot.Message{{Content: fmt.Sprintf("续断器 %d 执行 %s.%s 命令时发生错误: %v", n, script, cmd, err)}}
			}

			bot.Ruleset = cmd

			for _, msg := range sf(nil, bot) {
				msg.To = in.From
				bot.Send(msg)
			}
		}
	}

}

// 关闭 ticker
func stop_ticker(bot *rboot.Robot) []rboot.Message {
	args := bot.Args

	tNS := args[1]
	tNum, err := strconv.Atoi(tNS)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if mt, ok := tickers[tNum]; ok {
		mt.t.Stop()

		delete(tickers, tNum)

		return []rboot.Message{{Content: "续断器 " + tNS + " 关闭"}}
	} else {
		return []rboot.Message{{Content: "未找到序号为 " + tNS + " 的续断器"}}
	}
}

func status_ticker(bot *rboot.Robot) []rboot.Message {

	content := ""
	if len(tickers) <= 0 {
		content = "没有正在进行的续断器"
	} else {
		content = "续断器状态: \n"

		for k, t := range tickers {
			content += fmt.Sprintf("续断器 %d, 命令 %s, 已经执行了 %d 次，下次执行时间 %s \n", k, t.cmd, t.n, t.next.Format("15:04:05"))
		}
	}

	return []rboot.Message{{Content: content}}
}

func init() {
	// 注册脚本
	rboot.RegisterScripts(`timing`, rboot.Script{
		Action: setup,
		Ruleset: map[string]string{
			`timer`:         `^!(\d+)([小时|H|h|分|分钟|M|m|秒|S|s]{1,2})后执行(.+)\.(.+)`,
			`stop_timer`:    `^!stop timer (\d+)`,
			`status_timer`:  `^!timer status`,
			`ticker`:        `^!每过(\d+)([小时|H|h|分|分钟|M|m|秒|秒钟|S|s]{1,2})执行(.+)\.(.+)`,
			`stop_ticker`:   `^!stop ticker (\d+)`,
			`status_ticker`: `^!ticker status`,
		},
		Usage: "> `!<N小时|分|秒>后执行<脚本名称>.<命令名称>`: 定时器开始定时任务并在倒计时结束时执行相应命令 \n" +
			"> `!stop timer <N>`: 结束对应序号定时器 \n" +
			"> `!timer status`: 定时器状态 \n" +
			"> `!每过<N小时|分|秒>执行<脚本名称>.<命令名称>`: 每过相应时间执行一次对应脚本命令（循环） \n" +
			"> `!stop ticker <N>`: 结束对应序号续断器 \n" +
			"> `!ticker status`: 续断器状态",
		Description: `定时任务脚本。查看帮助信息: !help timing`,
	})
}

// 将小时、分、秒转换为time.Duration
func toDuration(t int, arc string) (time.Duration, error) {

	switch arc {
	case "小时", "H", "h":
		return time.Duration(t) * time.Hour, nil
	case "分", "分钟", "M", "m":
		return time.Duration(t) * time.Minute, nil
	case "秒", "秒钟", "S", "s":
		return time.Duration(t) * time.Second, nil
	}

	return 0, fmt.Errorf("非法字符串")
}

// 将时间戳转换为 汉字 的时分秒
func toCTime(t int64) string {
	left := ""

	if t > 60*60 {
		left += fmt.Sprintf("%d小时", t/(60*60))
	}

	if t%(60*60) > 60 {
		left += fmt.Sprintf("%d分钟", (t%(60*60))/60)
	}

	if (t%(60*60))%60 > 0 {
		left += fmt.Sprintf("%d秒", (t%(60*60))%60)
	}

	return left
}
