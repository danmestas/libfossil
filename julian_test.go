package libfossil

import (
	"math"
	"testing"
	"time"
)

func TestTimeToJulian(t *testing.T) {
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	got := TimeToJulian(epoch)
	if math.Abs(got-2440587.5) > 0.0001 {
		t.Fatalf("TimeToJulian(epoch) = %f, want 2440587.5", got)
	}
}

func TestJulianToTime(t *testing.T) {
	got := JulianToTime(2440587.5)
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	diff := got.Sub(epoch)
	if diff < -time.Second || diff > time.Second {
		t.Fatalf("JulianToTime(2440587.5) = %v, want %v (diff=%v)", got, epoch, diff)
	}
}

func TestJulianRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	julian := TimeToJulian(now)
	back := JulianToTime(julian)
	diff := now.Sub(back)
	if diff < -time.Millisecond || diff > time.Millisecond {
		t.Fatalf("Round-trip failed: %v -> %f -> %v (diff=%v)", now, julian, back, diff)
	}
}

func TestKnownJulianDates(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		julian float64
	}{
		{"J2000", time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC), 2451545.0},
		{"2026-03-14 noon", time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC), 2461114.0},
	}
	for _, tt := range tests {
		got := TimeToJulian(tt.t)
		if math.Abs(got-tt.julian) > 0.001 {
			t.Errorf("%s: TimeToJulian = %f, want %f", tt.name, got, tt.julian)
		}
	}
}
