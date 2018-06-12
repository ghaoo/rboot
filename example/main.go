package main

import (
	"rboot"

	_ "rboot/adapter/cli"
)

func main() {
	bot := rboot.NewRboot()

	bot.Go()
}
