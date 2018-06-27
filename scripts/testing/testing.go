package testing

import (
	"rboot"
)

func parse(bot *rboot.Robot) error {
	switch bot.Matcher {
	case `123`:
		bot.Send(`1 or 2 or 3`)

	case `abc`:
		bot.Send(`a or b or c`)

	}

	return nil
}

func hook(bot *rboot.Robot) {
	//
}

func init() {
	rboot.RegisterScript(`testing`, &rboot.Script{
		Action: parse,
		Ruleset: map[string]string{
			`123`: `1|2|3`,
			`abc`: `a|b|c`,
		},
		Hook: hook,
	})
}
