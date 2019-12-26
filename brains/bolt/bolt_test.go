package bolt

import (
	"testing"
)

func TestBlot(t *testing.T) {
	b := Bolt()

	err := b.Set(`key`, []byte(`1`))

	if err != nil {
		t.Error(`bolt save items error:`, err)
	}

	item := b.Get(`key`)

	if item == nil || string(item) != `1` {
		t.Error(`item not found`)
	}

	err = b.Remove(`key`)

	if err != nil {
		t.Error(`remove item error`, err)
	}

}
