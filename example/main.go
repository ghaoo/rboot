package main

import (
	"rboot"

	_ "rboot/provider"
	_ "rboot/scripts"
)

func main() {
	bot := rboot.NewRboot()

	bot.Go()
}
