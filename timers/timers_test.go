package timers

import (
	"slices"
	"testing"
	"time"
)

func TestTimersOrder(t *testing.T) {
	h := NewTimersWorker()
	defer h.Stop()

	t1 := h.SetTimeout(100 * time.Millisecond)
	t2 := h.SetTimeout(200 * time.Millisecond)
	t3 := h.SetTimeout(500 * time.Millisecond)
	t4 := h.SetInterval(300 * time.Millisecond)

	expected := []TimerId{t1, t2, t4, t3, t4}

	got := make([]TimerId, 0, len(expected))
	deadline := time.After(time.Second)

	for len(got) < len(expected) {
		select {
		case id := <-h.C:
			got = append(got, id)
		case <-deadline:
			t.Fatalf("timed out: expected %v, got %v", expected, got)
		}
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}
