package memory

import (
	"fmt"
	"sync"

	"github.com/ghaoo/rboot"
)

type memory struct {
	mu    sync.Mutex
	items map[string][]byte
	err   error
}

// New constructs memory
func New() rboot.Memorizer {
	return &memory{
		mu:    sync.Mutex{},
		items: make(map[string][]byte),
	}
}

// save ...
func (m *memory) Save(key string, value []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.items[key]
	if !ok {
		m.err = fmt.Errorf("key %s already existed, If you want to change its value, please use `Update`.", key)
		return
	}

	m.items[key] = value
}

// read ...
func (m *memory) Read(key string) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.items[key]
	if !ok {
		return []byte{}
	}
	return v
}

// update ...
func (m *memory) Update(key string, value []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.items[key]
	if !ok {
		m.items[key] = value
	}

	m.items[key] = value
}

// delete ...
func (m *memory) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
}

func (m *memory) Error() error {
	return m.err
}

func init() {
	rboot.RegisterMemorizer(`memory`, New)
}
