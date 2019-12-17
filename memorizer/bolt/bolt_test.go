package bolt

import (
	"testing"
)

func TestBlot(t *testing.T) {
	b := Bolt()

	err := b.Save(`test`, `key`, []byte(`1`))

	if err != nil {
		t.Error(`bolt save items error:`, err)
	}

	item := b.Find(`test`, `key`)

	if item == nil || string(item) != `1` {
		t.Error(`item not found`)
	}

	err = b.Delete(`test`, `key`)

	if err != nil {
		t.Error(`remove item error`, err)
	}

	err = b.Save(`test`, `key`, []byte(`2222`))

	if err != nil {
		t.Error(`bolt save items error:`, err)
	}

	item = b.Find(`test`, `key`)

	if item == nil || string(item) != `2222` {
		t.Error(`item not found`)
	}

}
