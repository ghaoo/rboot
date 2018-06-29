package memory

import (
	"fmt"
	"github.com/ghaoo/rboot"
	"sync"
)

type memory struct {
	mu    sync.Mutex
	items map[string][]byte
}

// New constructs memory
func New() rboot.Memorizer {
	return &memory{
		mu:    sync.Mutex{},
		items: make(map[string][]byte),
	}
}

func (m *memory) Open() error {
	return nil
}

func (m *memory) Close() error {
	return nil
}

// save ...
func (m *memory) Save(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.items[key]
	if !ok {
		return fmt.Errorf("key %s already existed, If you want to change its value, please use `Update`.", key)
	}

	m.items[key] = value

	return nil
}

// read ...
func (m *memory) Read(key string) ([]byte, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.items[key]
	if !ok {
		return []byte{}, false
	}
	return v, true
}

// update ...
func (m *memory) Update(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.items[key]
	if !ok {
		m.items[key] = value
	}

	m.items[key] = value
}

// delete ...
func (m *memory) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)

	return nil
}

func init() {
	rboot.RegisterMemorizer(`memory`, New)
}
