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
}

func (v *Vote) New(bot *rboot.Robot, to rboot.User, name, user string, opt string) []rboot.Message {
	opts := strings.Split(opt, " ")

	fmt.Println(opts)

	if len(opts) < 2 {
		return []rboot.Message{{Content: "选项最少两项"}}
	}

	v.Choices = make(map[string]int, len(opt))
	for _, c := range opts {
		v.Choices[c] = 0
	}

	v.Name = name
	v.User = user
	v.Players = make(map[string]string)
	v.startTime = time.Now()
	active = true

	go func() {
		v.ticker = time.NewTicker(timeout)
		select {
		case <-v.ticker.C:
			result := v.Result()

			result = append(result, rboot.Message{Content: "投票结束！"})

			active = false

			v.ticker.Stop()

			for _, res := range result {
				res.To = to
				bot.Send(res)
			}
		}
	}()

	return []rboot.Message{{Content: "请大家投票: " + name + ", 投票请直接输入@@选项"}}
}

func (v *Vote) Voting(user string, opt string) []rboot.Message {
	if !active {
		return []rboot.Message{{Content: "没有正在进行中的投票"}}
	}
	// 检查用户有没有参与
	if iopt, ok := v.Players[user]; ok {
		return []rboot.Message{{Content: "你已经参与了投票，你选择的是 " + iopt}}
	}

	opt = strings.TrimSpace(opt)

	// 检查选项是否存在
	if _, ok := v.Choices[opt]; !ok {
		return []rboot.Message{{Content: "滚!"}}
	}

	v.Players[user] = opt
	v.Choices[opt] += 1

	return []rboot.Message{{Content: "投票成功!"}}
}

func (v *Vote) Result() []rboot.Message {
	if !active {
		return []rboot.Message{{Content: "没有正在进行中的投票"}}
	}

	content := "投票: " + v.Name + "\n      "

	for choice, count := range v.Choices {
		content += fmt.Sprintf(" %d 人选择了 <%s> , ", count, choice)
	}

	content += "\n发起人: " + v.User

	return []rboot.Message{{Content: content}}
}

func (v *Vote) Stop(user string) []rboot.Message {

	if !active {
		return []rboot.Message{{Content: "没有正在进行中的投票"}}
	}

	if user != v.User {
		return []rboot.Message{{Content: "NO!"}}
	}

	result := v.Result()

	result = append(result, rboot.Message{Content: "投票结束！"})

	active = false

	v.ticker.Stop()

	return result
}
