package main

import (
	"rboot"

	_ "rboot/provider"
	_ "rboot/scripts"
)

func main() {
	bot := rboot.New()

	bot.Go()
}
