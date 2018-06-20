package main

import (
	"rboot"

	_ "rboot/connecter"
	_ "rboot/scripts"
)

func main() {
	bot := rboot.NewRboot()

	bot.Go()
}
