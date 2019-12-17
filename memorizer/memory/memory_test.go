package memory

import (
	"testing"
)

func TestMemory(t *testing.T) {
	mom := New()
	mom.Save("t1", "key", []byte("1"))
	mom.Save("t2", "key", []byte("2"))

	item := mom.Find("t1", "key")
	if string(item) != "1" {
		t.Error("read values error.  Got: key")
	}

	mom.Update("t1", "key", []byte("2"))

	item = mom.Find("t1", "key")

	if string(item) != "2" {
		t.Error("update item error.  Got: key")
	}

	mom.Delete("t1", "key")

	if item := mom.Find("t1", "key"); string(item) != `` {
		t.Error("delete item error.")
	}
}
