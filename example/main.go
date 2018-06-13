package main

import (
	"rboot"

	_ "rboot/connecter"
)

func main() {
	bot := rboot.NewRboot()

	bot.Go()
}
