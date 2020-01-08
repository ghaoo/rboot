package vote

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"strings"
	"time"
)

var timeout = 20 * time.Minute
var active bool // 投票是否进行中

type Vote struct {
	Name      string            // 投票名称
	User      string            // 发起人
	Choices   map[string]int    // 投票选项
	Players   map[string]string // 参与人数
	startTime time.Time         // 开始时间
	ticker    *time.Ticker      // 计时器
	bot       *rboot.Robot
}

// 新建投票
func (v *Vote) New(bot *rboot.Robot, in *rboot.Message, name string, opt string) *rboot.Message {
	// 检查有没有进行中的投票
	if active {
		return rboot.NewMessage("投票进行中，请稍后...")
	}
	opts := strings.Split(opt, " ")

	if len(opts) < 2 {
		return rboot.NewMessage("选项最少两项")
	}

	v.Choices = make(map[string]int, len(opt))
	for _, c := range opts {
		v.Choices[c] = 0
	}

	v.Name = name
	v.User = bot.GetUserName(in.Sender)
	v.Players = make(map[string]string)
	v.startTime = time.Now()
	v.bot = bot
	active = true

	go func() {
		v.ticker = time.NewTicker(timeout)
		select {
		case <-v.ticker.C:
			result := v.Result(in.From)

			bot.Send(result)

			active = false

			v.ticker.Stop()

			bot.SendText(fmt.Sprintf("`%s` 投票结束！", name), in.From)
		}
	}()

	msg := fmt.Sprintf("%s 创建了投票: %s\n> 选项:\n", v.User, name)
	for i, c := range opts {
		msg += fmt.Sprintf("> %d. `%s`\n", i+1, c)
	}

	msg += "\n*投票请直接输入@@选项*"

	return rboot.NewMessage(msg, in.From)
}

func (v *Vote) Voting(user string, opt string) *rboot.Message {
	if !active {
		return rboot.NewMessage("没有正在进行中的投票或投票已经结束！")
	}
	// 检查用户有没有参与
	if iopt, ok := v.Players[user]; ok {
		return rboot.NewMessage(fmt.Sprintf("%s 你已经参与了投票，你选择的是`%s`", user, iopt))
	}

	opt = strings.TrimSpace(opt)

	// 检查选项是否存在
	if _, ok := v.Choices[opt]; !ok {
		return rboot.NewMessage("投票失败！没有这个选项胸弟！")
	}

	v.Players[user] = opt
	v.Choices[opt] += 1

	return rboot.NewMessage("投票成功!")
}

func (v *Vote) Result(to string) *rboot.Message {
	if !active {
		return rboot.NewMessage("没有正在进行中的投票或投票已经结束")
	}

	content := "投票: " + v.Name + "\n      "

	for choice, count := range v.Choices {
		content += fmt.Sprintf(" %d 人选择了 `%s` \n", count, choice)
	}

	content += "\n*发起人*: " + v.User

	fmt.Println(to)

	return rboot.NewMessage(content, to)
}

func (v *Vote) Stop(bot *rboot.Robot, user, to string) *rboot.Message {

	if !active {
		return rboot.NewMessage("没有正在进行中的投票或投票已经结束")
	}

	if user != v.User {
		return rboot.NewMessage("NO!")
	}

	result := v.Result(to)

	v.bot.Send(result)

	active = false

	v.ticker.Stop()

	msg := rboot.NewMessage(fmt.Sprintf("%s 创建的 `%s` 投票结束！", bot.GetUserName(user), v.Name), to)

	return msg
}
