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

// Scheduling new timers must not block while a fired timer is waiting to be read.
func TestSetTimeoutWhileFiredUnread(t *testing.T) {
	h := NewTimersWorker()
	defer h.Stop()

	t1 := h.SetTimeout(0)
	time.Sleep(50 * time.Millisecond)

	ids := make(chan []TimerId)
	go func() {
		t2 := h.SetTimeout(10 * time.Millisecond)
		t3 := h.SetTimeout(20 * time.Millisecond)
		ids <- []TimerId{t2, t3}
	}()

	var expected []TimerId
	select {
	case scheduled := <-ids:
		expected = append([]TimerId{t1}, scheduled...)
	case <-time.After(time.Second):
		t.Fatal("SetTimeout blocked while a fired timer was unread")
	}

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

// A timer cleared after firing, but before being read, must never be delivered.
func TestClearTimeoutAfterFired(t *testing.T) {
	h := NewTimersWorker()
	defer h.Stop()

	t1 := h.SetTimeout(0)
	time.Sleep(50 * time.Millisecond)
	h.ClearTimeout(t1)

	t2 := h.SetTimeout(10 * time.Millisecond)

	select {
	case id := <-h.C:
		if id != t2 {
			t.Fatalf("expected %v, got %v", t2, id)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for timer")
	}
}
