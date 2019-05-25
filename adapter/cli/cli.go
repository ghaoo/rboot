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

// New returns an initialized adapter
func New(bot *rboot.Rboot) rboot.Adapter {

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

func (c *cli) Send(msg rboot.Message) error {

	c.writeString(msg.Content)
	return nil
}

func (c *cli) Incoming() chan rboot.Message {
	return c.in
}

// Run executes the adapter run loop
func (c *cli) run() {

	go func() {
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {

			c.in <- rboot.Message{
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

	msg := fmt.Sprintf(" > %s\n", strings.TrimSpace(str))

	if _, err := c.writer.WriteString(msg); err != nil {
		return err
	}

	return c.writer.Flush()
}

func init() {
	rboot.RegisterAdapter(`cli`, New)
}
