package rboot

import "testing"

var m = newMemory()

func TestMemory_Set(t *testing.T) {
	err := m.Set("test", "key", []byte("value"))

	if err != nil {
		t.Error(err)
	}
}

func TestMemory_Get(t *testing.T) {
	v := m.Get("test", "key")

	if string(v) != "value" {
		t.Error("failed")
	}
}

func TestMemory_Remove(t *testing.T) {
	err := m.Remove("test", "key")

	if err != nil {
		t.Error(err)
	}

	v := m.Get("test", "key")

	if string(v) == "value" {
		t.Error("failed")
	}
}
