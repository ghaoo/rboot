package memory

import (
	"github.com/ghaoo/rboot"
	"sync"
)

// memory brain
type memory struct {
	mu    sync.Mutex
	items map[string][]byte
}

// New constructs memory
func newMemory() rboot.Brain {
	return &memory{
		mu:    sync.Mutex{},
		items: make(map[string][]byte),
	}
}

// save ...
func (m *memory) Set(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = value

	return nil
}

// find ...
func (m *memory) Get(key string) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.items[key]
	if !ok {
		return []byte{}
	}
	return v
}

// delete ...
func (m *memory) Remove(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)

	return nil
}

// register brain ...
func init() {
	rboot.RegisterBrain("memory", newMemory)
}
