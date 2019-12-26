package rboot

import (
	"fmt"
)

type Brain interface {
	Set(key string, value []byte) error
	Get(key string) []byte
	Remove(key string) error
}

var brains = make(map[string]func() Brain)

// 注册存储器
func RegisterBrain(name string, m func() Brain) {

	if name == "" {
		panic("RegisterBrain: brain must have a name")
	}
	if _, ok := brains[name]; ok {
		panic("RegisterBrain: brains named " + name + " already registered. ")
	}
	brains[name] = m
}

func DetectBrain(name string) (func() Brain, error) {
	if brain, ok := brains[name]; ok {
		return brain, nil
	}

	if len(brains) == 0 {
		return nil, fmt.Errorf("no Brain available")
	}

	if name == "" {
		if len(brains) == 1 {
			for _, brain := range brains {
				return brain, nil
			}
		}
		return nil, fmt.Errorf("multiple brains available; must choose one")
	}
	return nil, fmt.Errorf("unknown brain '%s'", name)
}


