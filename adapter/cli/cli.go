package cli

import (
	"bufio"
	"fmt"
	"github.com/ghaoo/rboot"
	"io"
	"os"
	"strings"
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

// New returns an initialized adapter
func New(bot *rboot.Robot) rboot.Adapter {

	c := &cli{
		in:     make(chan rboot.Message),
		out:    make(chan rboot.Message),
		writer: bufio.NewWriter(stdout),
	}

	go c.run()

	return c
}

func (c *cli) Name() string {
	return `cli`
}

func (c *cli) Incoming() chan rboot.Message {
	return c.in
}

func (c *cli) Outgoing() chan rboot.Message {
	return c.out
}

// Run executes the adapter run loop
func (c *cli) run() {

	go func() {
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {

			c.in <- rboot.Message{
				To:      rboot.User{Name: `cli`},
				From:    rboot.User{Name: `cli`},
				Channel: `cli`,
				Content: scanner.Text(),
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
		}
	}()
}

func (c *cli) writeString(str string) error {

	name := os.Getenv(`RBOOT_ALIAS`)
	if name == `` {
		name = os.Getenv(`RBOOT_NAME`)
	}

	msg := fmt.Sprintf(name+" > %s\n", strings.TrimSpace(str))

	if _, err := c.writer.WriteString(msg); err != nil {
		return err
	}

	return c.writer.Flush()
}

func init() {
	rboot.RegisterAdapter(`cli`, New)
}
