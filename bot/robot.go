package bot

import (
	"fmt"
	"log"
	"strings"
)

const (
	DefaultRobotName = `Rboot`
)

type Robot struct {
	name string
}

func New() *Robot {
	return &Robot{
		name: DefaultRobotName,
	}
}

func (bot *Robot) Name() string {
	return bot.name
}
