package simio

import (
	"testing"
)

func TestRealEnv(t *testing.T) {
	env := RealEnv()
	if env.Clock == nil {
		t.Fatal("Clock should not be nil")
	}
	if env.Rand == nil {
		t.Fatal("Rand should not be nil")
	}

	// RealClock.Now should return a reasonable time
	now := env.Clock.Now()
	if now.Year() < 2024 {
		t.Fatalf("RealClock.Now returned suspicious time: %v", now)
	}

	// CryptoRand should read bytes
	buf := make([]byte, 16)
	n, err := env.Rand.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 16 {
		t.Fatalf("expected 16 bytes, got %d", n)
	}
}

func TestSimEnv(t *testing.T) {
	env := SimEnv(99)
	if env.Clock == nil {
		t.Fatal("Clock should not be nil")
	}
	if env.Rand == nil {
		t.Fatal("Rand should not be nil")
	}

	// SimClock should be a *SimClock
	sc, ok := env.Clock.(*SimClock)
	if !ok {
		t.Fatalf("expected *SimClock, got %T", env.Clock)
	}
	if !sc.Now().IsZero() && sc.Now().Unix() != 0 {
		t.Fatalf("SimClock should start at epoch, got %v", sc.Now())
	}

	// SeededRand should be deterministic
	env2 := SimEnv(99)
	buf1 := make([]byte, 32)
	buf2 := make([]byte, 32)
	env.Rand.Read(buf1)
	env2.Rand.Read(buf2)

	for i := range buf1 {
		if buf1[i] != buf2[i] {
			t.Fatal("SimEnv with same seed should produce identical Rand output")
		}
	}
}
