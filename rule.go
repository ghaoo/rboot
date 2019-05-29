package rboot

import "regexp"

type Rule interface {
	Match(pattern, msg string) bool
}

type Regex struct {
}

func (reg *Regex) Match(pattern, msg string) bool {
	r := regexp.MustCompile(pattern)

	if r.MatchString(msg) {
		return true
	}

	return false
}

func (bot *Rboot) MatchRuleset(pattern string) (plug, match string, matched bool) {

	for plug, rule := range rulesets {
		for m, r := range rule {
			if bot.Rule.Match(r, pattern) {
				return plug, m, true
			}
		}
	}

	return ``, ``, false
}
