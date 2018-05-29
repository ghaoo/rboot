package rboot

type Context struct {
	bot *Rboot
}

func (ctx *Context) BotName() string {
	return ctx.bot.name
}

func (ctx *Context) MessageIn() chan Message {
	return ctx.bot.providerIn
}

func (ctx *Context) MessageOut() chan Message {
	return ctx.bot.providerOut
}
