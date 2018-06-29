package main

import (
	"github.com/ghaoo/rboot"

	_ "github.com/ghaoo/rboot/memorizer"
	_ "github.com/ghaoo/rboot/provider"
	_ "github.com/ghaoo/rboot/scripts"
)

func main() {
	bot := rboot.New()

	bot.Go()
}
