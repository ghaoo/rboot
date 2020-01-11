package rboot

import (
	"regexp"
)

// Rule 消息处理器接口，暂时只支持正则
type Rule interface {
	Match(pattern, msg string) ([]string, bool)
}

// Regex 正则消息处理器
type Regex struct{}

// Match 当匹配失败时返回 false，返回true证明匹配成功并返回消息的匹配文本和子表达式的匹配
func (reg *Regex) Match(pattern, msg string) ([]string, bool) {
	r := regexp.MustCompile(pattern)

	if submatch := r.FindStringSubmatch(msg); submatch != nil {
		return submatch, true
	}

	return nil, false
}
