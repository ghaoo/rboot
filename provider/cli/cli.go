package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ghaoo/rboot"
)

var (
	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
)

type cli struct {
	in     chan rboot.Message
	out    chan rboot.Message
	writer *bufio.Writer
}

// 初始化cli连接器
func NewCli() rboot.Provider {
	c := &cli{
		in:     make(chan rboot.Message),
		out:    make(chan rboot.Message),
		writer: bufio.NewWriter(stdout),
	}

	go c.run()
	return c
}

func (c *cli) Incoming() chan rboot.Message {
	return c.in
}

func (c *cli) Outgoing() chan rboot.Message {
	return c.out
}

func (c *cli) Error() error {
	return nil
}

func (c *cli) run() {
	go func() {

		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			prompt()

			header := make(rboot.Header)
			header.Add(`From`, `CLI`)
			header.Add(`To`, `CLI`)

			c.in <- rboot.Message{
				Header: header,
				Content:   scanner.Text(),
			}


		forLoop:
			for {
				select {
				case msg := <-c.out:
					c.writeString(msg.Content)
				default:
					break forLoop
				}
			}
		}
	}()

	go func() {
		for msg := range c.out {
			c.writeString(msg.Content)
			prompt()
		}
	}()

}

func prompt() {
	name := rboot.DefaultRobotName
	if os.Getenv(`ROBOT_NAME`) != `` {
		name = os.Getenv(`ROBOT_NAME`)
	}
	fmt.Print(Color(name, FgRed)+`> `)
}

func (c *cli) writeString(str string) error {
	msg := fmt.Sprintf("%s\n", strings.TrimSpace(str))

	if _, err := c.writer.WriteString(msg); err != nil {
		return err
	}

	return c.writer.Flush()
}

func init() {
	rboot.RegisterProvider(`cli`, NewCli)
}
