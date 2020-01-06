package rboot

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Adapter interface {
	Name() string           // 适配器名称
	Incoming() chan Message // 接收到的消息
	Outgoing() chan Message // 回复的消息
}

type adapterF func(*Robot) Adapter

var adapters = make(map[string]adapterF)

// RegisterAdapter 注册适配器，名称不可重复
// 适配器需实现 Adapter 接口
func RegisterAdapter(name string, adp adapterF) {
	if name == "" {
		panic("RegisterAdapter: adapter must have a name")
	}
	if _, ok := adapters[name]; ok {
		panic("RegisterAdapter: adapter named " + name + " already registered. ")
	}
	adapters[name] = adp
}

// DetectAdapter 根据适配器名称获取适配器实例
func DetectAdapter(name string) (adapterF, error) {
	if adp, ok := adapters[name]; ok {
		return adp, nil
	}

	if len(adapters) == 0 {
		return nil, errors.New("no adapter available")
	}

	if name == "" {
		if len(adapters) == 1 {
			for _, adp := range adapters {
				return adp, nil
			}
		}
		return nil, errors.New("multiple adapters available; must choose one")
	}
	return nil, errors.New("unknown adapter " + name)
}

var (
	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
)

type cli struct {
	in     chan Message
	out    chan Message
	writer *bufio.Writer
}

// New returns an initialized adapter
func newCli(bot *Robot) Adapter {

	c := &cli{
		in:     make(chan Message),
		out:    make(chan Message),
		writer: bufio.NewWriter(stdout),
	}

	go c.run()

	return c
}

func (c *cli) Name() string {
	return `cli`
}

func (c *cli) Incoming() chan Message {
	return c.in
}

func (c *cli) Outgoing() chan Message {
	return c.out
}

// Run executes the adapter run loop
func (c *cli) run() {

	go func() {
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {

			msg := NewMessage(scanner.Text())

			c.in <- msg

		forLoop:
			for {
				select {
				case msg := <-c.out:
					c.writeString(msg.String())
				default:
					break forLoop
				}
			}
		}
	}()

	go func() {
		for msg := range c.out {
			c.writeString(msg.String())
		}
	}()
}

func (c *cli) writeString(str string) error {

	name := os.Getenv(`RBOOT_ALIAS`)
	if name == `` {
		name = os.Getenv(`RBOOT_NAME`)
	}

	msg := fmt.Sprintf(name+"> %s\n", strings.TrimSpace(str))

	if _, err := c.writer.WriteString(msg); err != nil {
		return err
	}

	return c.writer.Flush()
}

func init() {
	RegisterAdapter(`cli`, newCli)
}
