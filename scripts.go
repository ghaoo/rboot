package rboot

import (
	"context"
	"fmt"
)

var (
	scripts  = make(map[string]Script)
	rulesets = make(map[string]map[string]string)
)

type Script struct {
	Action      SetupFunc         // 执行解析或一些必要加载
	Ruleset     map[string]string // 脚本规则集合
	Usage       string            // 帮助信息
	Description string            // 简介
}

type SetupFunc func(context.Context, *Robot) []Message

func RegisterScripts(name string, script Script) {

	if name == "" {
		panic("RegisterScripts: the script must have a name")
	}
	if _, ok := scripts[name]; ok {
		panic("RegisterScripts: script named " + name + " already registered.")
	}

	scripts[name] = script

	if len(script.Ruleset) > 0 {

		rulesets[name] = script.Ruleset
	}
}

func DirectiveScript(name string) (SetupFunc, error) {

	if script, ok := scripts[name]; ok {
		return script.Action, nil
	}

	return nil, fmt.Errorf("DirectiveScript: no action found in script '%s' (missing a script?)", name)
}

func GetScripts() map[string]Script {
	return scripts
}
