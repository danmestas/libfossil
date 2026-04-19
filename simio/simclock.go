package simio

import (
	"sync"
	"time"
)

// SimClock is a virtual clock for deterministic simulation.
// Time only advances when Advance or AdvanceTo is called explicitly.
type SimClock struct {
	mu      sync.Mutex
	now     time.Time
	waiters []waiter
}

type waiter struct {
	deadline time.Time
	ch       chan time.Time
}

// NewSimClock creates a SimClock starting at Unix epoch zero.
func NewSimClock() *SimClock {
	return &SimClock{now: time.Unix(0, 0)}
}

// NewSimClockAt creates a SimClock starting at the given time.
func NewSimClockAt(t time.Time) *SimClock {
	return &SimClock{now: t}
}

func (c *SimClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// After returns a channel that receives the current time once the clock
// has been advanced past the deadline. The channel is never selected
// until Advance/AdvanceTo is called.
func (c *SimClock) After(d time.Duration) <-chan time.Time {
	if d < 0 {
		panic("simclock.After: duration must not be negative")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := make(chan time.Time, 1)
	deadline := c.now.Add(d)
	if !c.now.Before(deadline) {
		ch <- c.now
		return ch
	}
	c.waiters = append(c.waiters, waiter{deadline: deadline, ch: ch})
	return ch
}

// Sleep blocks until the clock has been advanced past now+d.
// In simulation this is typically called from the simulator driving
// the clock forward on another goroutine, or within a single-threaded
// event loop that calls Advance between steps.
func (c *SimClock) Sleep(d time.Duration) {
	<-c.After(d)
}

// Advance moves the clock forward by d and fires any expired waiters.
func (c *SimClock) Advance(d time.Duration) {
	c.AdvanceTo(c.Now().Add(d))
}

// AdvanceTo sets the clock to t (must be >= current time) and fires
// any waiters whose deadline has been reached.
func (c *SimClock) AdvanceTo(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if t.Before(c.now) {
		panic("simclock.AdvanceTo: cannot move time backwards")
	}
	c.now = t
	remaining := c.waiters[:0]
	for _, w := range c.waiters {
		if !c.now.Before(w.deadline) {
			w.ch <- c.now
		} else {
			remaining = append(remaining, w)
		}
	}
	c.waiters = remaining
}

// PendingCount returns the number of waiters that haven't fired yet.
// Useful for the simulator to know if any node is sleeping.
func (c *SimClock) PendingCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.waiters)
}
