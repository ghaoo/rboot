package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"rboot"
	"strings"
)

type cli struct {
	*rboot.Response
	quit   chan bool
	writer *bufio.Writer
}

func NewCli(res *rboot.Response) rboot.Connecter {
	c := &cli{
		Response: res,
		quit:     make(chan bool),
		writer:   bufio.NewWriter(os.Stdout),
	}
	return c
}

func (c *cli) Name() string {
	return `CLI`
}

func (c *cli) Send(strings ...string) error {
	for _, str := range strings {
		err := c.writeString(str)
		if err != nil {
			log.Printf("send message error: %v", err)
			return err
		}
	}

	return nil
}

func (c *cli) Reply(strings ...string) error {
	for _, str := range strings {
		err := c.writeString(str)
		if err != nil {
			log.Printf("reply message error: %v", err)
			return err
		}
	}

	return nil
}

func (c *cli) Run() error {

	go func() {
		for {

			prompt()
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()

			line := scanner.Bytes()

			header := make(rboot.Header)
			header.Add(`From`, `CLI`)
			header.Add(`To`, `CLI`)

			msg := &rboot.Message{
				Header: header,
				Body:   bytes.NewReader(line),
			}

			c.Receive(msg)

			continue
		}
	}()

	<-c.quit
	return nil
}

func (c *cli) Close() error {
	c.quit <- true
	return nil
}

func prompt() {
	fmt.Print("> ")
}

func (c *cli) writeString(str string) error {
	msg := fmt.Sprintf("%s\n", strings.TrimSpace(str))

	if _, err := c.writer.WriteString(msg); err != nil {
		return err
	}

	if err := c.writer.Flush(); err != nil {
		return err
	}

	return nil
}

func init() {
	rboot.RegisterConnecter(`cli`, NewCli)
}
