// Package fsltype defines shared types used across internal packages.
// These are re-exported by the root libfossil package.
package fsltype

import "time"

// FslID is a row-id in the blob table (content-addressed artifacts).
type FslID int64

// FslSize represents a blob size; negative values indicate phantom blobs.
type FslSize int64

const (
	// PhantomSize is the sentinel size for phantom (not-yet-received) blobs.
	PhantomSize FslSize = -1

	// FossilApplicationID is the SQLite application_id for Fossil repositories.
	FossilApplicationID int32 = 252006673
)

const julianEpoch = 2440587.5

// TimeToJulian converts a time.Time to a Fossil Julian day number.
func TimeToJulian(t time.Time) float64 {
	return julianEpoch + float64(t.UTC().UnixMilli())/(86400.0*1000.0)
}

// JulianToTime converts a Fossil Julian day number to time.Time.
func JulianToTime(j float64) time.Time {
	millis := int64((j - julianEpoch) * 86400.0 * 1000.0)
	return time.UnixMilli(millis).UTC()
}
