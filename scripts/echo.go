package scripts

import (
	"context"
	"github.com/ghaoo/rboot"
)

func init() {
	rboot.RegisterScripts(`echo`, rboot.Script{
		Action: func(ctx context.Context, bot *rboot.Robot) []rboot.Message {
			bot.SendText("Hello World!")
			return nil
		},
		Ruleset:     map[string]string{"hello": "hello"},
		Description: `返回 Hello World! `,
	})
}
