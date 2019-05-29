package rboot

import "testing"

func TestNewQ(t *testing.T) {
	q := NewQ(10)

	if len(q.items) != 1 {
		t.Error("Error in size of queue")
	}

	if cap(q.items) != 1 {
		t.Error("Error with capacity of queue")
	}
}

