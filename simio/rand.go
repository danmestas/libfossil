package simio

import (
	"crypto/rand"
	"io"
	mrand "math/rand"
)

// Rand abstracts random byte generation for deterministic simulation.
type Rand interface {
	Read(p []byte) (n int, err error)
}

// CryptoRand delegates to crypto/rand (production use).
type CryptoRand struct{}

func (CryptoRand) Read(p []byte) (int, error) { return rand.Read(p) }

// SeededRand uses a seeded math/rand source for deterministic output.
type SeededRand struct {
	src io.Reader
}

func NewSeededRand(seed int64) *SeededRand {
	return &SeededRand{src: mrand.New(mrand.NewSource(seed))}
}

func (r *SeededRand) Read(p []byte) (int, error) {
	return r.src.Read(p)
}
