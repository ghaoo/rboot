package Rboot

import "fmt"

type Store interface {
	Name() string
	Open() error
	Close() error
	Get(string) ([]byte, error)
	Set(key string, data []byte) error
	Delete(string) error
}

type BasicStore struct {
	Robot *Robot
}

func (s *BasicStore) SetRobot(r *Robot) {
	s.Robot = r
}

var availableStores = make(map[string]func(*Robot) (Store, error))

func NewStore(robot *Robot) (Store, error) {

	default_store := `memory`

	if _, ok := availableStores[default_store]; !ok {

		return nil, fmt.Errorf("%s is not a registered store", default_store)
	}

	store, err := availableStores[default_store](robot)

	if err != nil {
		return nil, err
	}
	return store, nil
}

func RegisterStore(name string, f func(*Robot) (Store, error)) {
	availableStores[name] = f
}

// 默认实现 memory

type memory struct {
	BasicStore
	data map[string][]byte
}

func NewMemory(robot *Robot) (Store, error) {
	m := &memory{
		data: map[string][]byte{},
	}
	m.SetRobot(robot)
	return m, nil
}

func (m *memory) Name() string {
	return `memory`
}

func (m *memory) Open() error {
	return nil
}

func (m *memory) Close() error {
	return nil
}

func (m *memory) Get(key string) ([]byte, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("key %s was not found", key)
}

func (m *memory) Set(key string, data []byte) error {
	m.data[key] = data
	return nil
}

func (m *memory) Delete(key string) error {
	if _, ok := m.data[key]; !ok {
		return fmt.Errorf("key %s was not found", key)
	}
	delete(m.data, key)
	return nil
}

func init() {
	RegisterStore("memory", NewMemory)
}
