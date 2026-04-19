package simio

import (
	"testing"
	"time"
)

func TestSimClockNow(t *testing.T) {
	c := NewSimClock()
	if got := c.Now(); !got.Equal(time.Unix(0, 0)) {
		t.Fatalf("expected epoch, got %v", got)
	}
}

func TestSimClockAt(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	c := NewSimClockAt(start)
	if got := c.Now(); !got.Equal(start) {
		t.Fatalf("expected %v, got %v", start, got)
	}
}

func TestSimClockAdvance(t *testing.T) {
	c := NewSimClock()
	c.Advance(5 * time.Second)
	want := time.Unix(5, 0)
	if got := c.Now(); !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestSimClockAdvanceToBackwards(t *testing.T) {
	c := NewSimClockAt(time.Unix(100, 0))
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for backwards AdvanceTo")
		}
	}()
	c.AdvanceTo(time.Unix(50, 0))
	t.Fatal("should not reach here")
}

func TestSimClockAfterFiresOnAdvance(t *testing.T) {
	c := NewSimClock()
	ch := c.After(10 * time.Second)

	select {
	case <-ch:
		t.Fatal("should not fire before advance")
	default:
	}

	c.Advance(10 * time.Second)

	select {
	case got := <-ch:
		if !got.Equal(time.Unix(10, 0)) {
			t.Fatalf("expected t=10s, got %v", got)
		}
	default:
		t.Fatal("should have fired after advance")
	}
}

func TestSimClockAfterZeroDuration(t *testing.T) {
	c := NewSimClock()
	ch := c.After(0)
	select {
	case <-ch:
		// ok — zero duration fires immediately
	default:
		t.Fatal("After(0) should fire immediately")
	}
}

func TestSimClockMultipleWaiters(t *testing.T) {
	c := NewSimClock()
	ch5 := c.After(5 * time.Second)
	ch10 := c.After(10 * time.Second)
	ch15 := c.After(15 * time.Second)

	c.Advance(10 * time.Second)

	// ch5 and ch10 should fire, ch15 should not
	select {
	case <-ch5:
	default:
		t.Fatal("ch5 should have fired")
	}
	select {
	case <-ch10:
	default:
		t.Fatal("ch10 should have fired")
	}
	select {
	case <-ch15:
		t.Fatal("ch15 should not have fired")
	default:
	}

	if c.PendingCount() != 1 {
		t.Fatalf("expected 1 pending, got %d", c.PendingCount())
	}

	c.Advance(5 * time.Second)
	select {
	case <-ch15:
	default:
		t.Fatal("ch15 should have fired after second advance")
	}

	if c.PendingCount() != 0 {
		t.Fatalf("expected 0 pending, got %d", c.PendingCount())
	}
}

func TestSimClockSleep(t *testing.T) {
	c := NewSimClock()
	done := make(chan struct{})
	go func() {
		c.Sleep(5 * time.Second)
		close(done)
	}()

	// Give goroutine a moment to register the waiter
	time.Sleep(10 * time.Millisecond)
	c.Advance(5 * time.Second)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Sleep did not unblock after Advance")
	}
}
