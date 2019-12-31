package rboot

import (
	"regexp"
)

type Rule interface {
	Match(pattern, msg string) ([]string, bool)
}

type Regex struct{}

func (reg *Regex) Match(pattern, msg string) ([]string, bool) {
	r := regexp.MustCompile(pattern)

	if submatch := r.FindStringSubmatch(msg); submatch != nil {
		return submatch, true
	}

	return nil, false
}

