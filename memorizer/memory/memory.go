package memory

import (
	"sync"
	"fmt"
)

type memory struct {
	mu    sync.Mutex
	items map[string][]byte
	err error
}

// New constructs memory
func New() *memory {
	return &memory{
		mu: sync.Mutex{},
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

// error ...
func (m *memory) Error() error {
	return m.err
}

func init() {
}