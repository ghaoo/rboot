package rboot

import (
	"fmt"
	"sync"
)

// Brain 是Rboot缓存器实现的接口
type Brain interface {
	Set(bucket, key string, value []byte) error
	Get(bucket, key string) []byte
	Remove(bucket, key string) error
}

var brains = make(map[string]func() Brain)

// RegisterBrain 注册存储器，名称须唯一
// 需实现Brain接口
func RegisterBrain(name string, m func() Brain) {

	if name == "" {
		panic("RegisterBrain: brain must have a name")
	}
	if _, ok := brains[name]; ok {
		panic("RegisterBrain: brains named " + name + " already registered. ")
	}
	brains[name] = m
}

// DetectBrain 获取名称为 name 的缓存器
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

// memory the default brain
type memory struct {
	mu    sync.Mutex
	items map[string][]byte
}

// New constructs memory
func newMemory() Brain {
	return &memory{
		mu:    sync.Mutex{},
		items: make(map[string][]byte),
	}
}

// save ...
func (m *memory) Set(bucket, key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[bucket+key] = value

	return nil
}

// find ...
func (m *memory) Get(bucket, key string) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.items[bucket+key]
	if !ok {
		return []byte{}
	}
	return v
}

// delete ...
func (m *memory) Remove(bucket, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, bucket+key)

	return nil
}

// register brain ...
func init() {
	RegisterBrain("memory", newMemory)
}
