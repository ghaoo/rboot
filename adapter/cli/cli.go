package cli

import (
	"bufio"
	"Rboot"
	"os"
)

type cli struct {
	quit   chan bool
	writer *bufio.Writer
}

func NewCli(c *rboot.Controller) *cli {
	return &cli{
		quit:   make(chan bool),
		writer: bufio.NewWriter(os.Stdout),
	}
}

//func (c *cli) Send()