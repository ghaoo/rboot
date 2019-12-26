package memory

import "testing"

var m = newMemory()

func TestMemory_Set(t *testing.T) {
	err := m.Set("key", []byte("value"))

	if err != nil {
		t.Error(err)
	}
}

func TestMemory_Get(t *testing.T) {
	v := m.Get("key")

	if string(v) != "value" {
		t.Error("failed")
	}
}

func TestMemory_Remove(t *testing.T) {
	err := m.Remove("key")

	if err != nil {
		t.Error(err)
	}

	v := m.Get("key")

	if string(v) == "value" {
		t.Error("failed")
	}
}
