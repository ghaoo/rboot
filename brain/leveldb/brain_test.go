package leveldb

import (
	"fmt"
	"testing"
)

func TestGoleveldb_Set(t *testing.T) {
	db := NewLevelDB()
	defer db.Close()

	err := db.Set("test", "test", []byte("GoLevelDB"))
	if err != nil {
		t.Error(err)
	}
}

func TestGoleveldb_Get(t *testing.T) {
	db := NewLevelDB()
	defer db.Close()

	val, err := db.Get("test", "test")

	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(val))
}

func TestGoleveldb_Remove(t *testing.T) {
	db := NewLevelDB()
	defer db.Close()

	err := db.Remove("test", "test")
	if err != nil {
		t.Error(err)
	}
}
