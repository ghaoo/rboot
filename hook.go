package rboot

type Hook func(*Robot)

func (bot *Robot) Add(hook Hook) {
	bot.hooks = append(bot.hooks, hook)
}

func (bot *Robot) Hook() {
	if len(bot.hooks) > 0 {
		for _, hook := range bot.hooks {
			hook(bot)
		}
	}
}
