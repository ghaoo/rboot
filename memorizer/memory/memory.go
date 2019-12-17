package memory

import (
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

// save ...
func (m *memory) Save(bucket, key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = value

	return nil
}

// find ...
func (m *memory) Find(bucket, key string) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.items[key]
	if !ok {
		return []byte{}
	}
	return v
}

func (m *memory) Update(bucket, key string, value []byte) error {
	return nil
}

// delete ...
func (m *memory) Delete(bucket, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)

	return nil
}

func (m *memory) FindAll(bucket string) map[string][]byte {
	return nil
}

func init() {
	rboot.RegisterMemorizer(`memory`, New)
}
