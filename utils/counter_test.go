package utils

import "testing"

func TestCounter_Get(t *testing.T) {
	var count = new(Counter)

	if score := count.Get(); score != 0 {
		t.Errorf("Expected score of 0, but it was %d instead.", score)
	}
}

func TestCounter_Increment(t *testing.T) {
	var count = new(Counter)

	count.Increment()
	if score := count.Get(); score != 1 {
		t.Errorf("Expected score of 1, but it was %d instead.", score)
	}
}

func TestCounter_Reset(t *testing.T) {
	var count = new(Counter)

	count.Increment()
	count.Increment()
	count.Reset()
	if score := count.Get(); score != 0 {
		t.Errorf("Expected score of 0, but it was %d instead.", score)
	}
}
