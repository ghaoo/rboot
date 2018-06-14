package main

import (
	"rboot"

	_ "rboot/connecter"
	_ "rboot/plugins"
)

func main() {
	bot := rboot.NewRboot()

	bot.Go()
}
