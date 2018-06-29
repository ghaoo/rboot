package rboot

import (
	"fmt"
	"log"
	"sync"
)

var (
	availableScripts = make(map[string]*Script)

	availableProviders = make(map[string]ProvFunc)
)

type ProvFunc func() Provider

type Script struct {
	Action      ParseFunc   // 操作函数
	Hook        func(Robot) // 钩子
	Usage       string      // 使用方法
	Description string      // 简介
}

type ParseFunc func(Robot, Message) []Message

// 注册插件
func RegisterScript(name string, script *Script) {

	if name == "" {
		panic("RegisterScript: script must have a name")
	}
	if _, ok := availableScripts[name]; ok {
		panic("RegisterScript: script named " + name + " already registered. ")
	}

	availableScripts[name] = script
}

type Provider interface {
	Incoming() chan Message
	Outgoing() chan Message
	Error() error
}

// register provider
func RegisterProvider(name string, f ProvFunc) {
	availableProviders[name] = f
}

// get provider by name
func Detect(name string) (Provider, error) {
	if f, ok := availableProviders[name]; ok {
		return f(), nil
	}

	if len(availableProviders) == 0 {
		return nil, fmt.Errorf("no provider available")
	}

	if name == "" {
		if len(availableProviders) == 1 {
			for _, f := range availableProviders {
				return f(), nil
			}
		}
		return nil, fmt.Errorf("multiple providers available; must choose one")
	}
	return nil, fmt.Errorf("unknown provider '%s'", name)
}

var processOnce sync.Once

func (bot *Robot) Process() {
	processOnce.Do(func() {

		for _, script := range availableScripts {
			script.Hook(*bot)
		}

		for in := range bot.providerIn {
			go func(bot Robot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()

				for _, script := range availableScripts {
					responses := script.Action(bot, in)

					for _, r := range responses {
						bot.providerOut <- r
					}
				}

			}(*bot, in)
		}
	})
}
