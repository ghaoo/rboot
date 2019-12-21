package main

import (
	_ "github.com/ghaoo/rboot/adapter"
	_ "github.com/ghaoo/rboot/memorizer"
	_ "github.com/ghaoo/rboot/scripts"

	"github.com/ghaoo/rboot"
)

func main() {
	bot := rboot.New()

	bot.Go()
}
