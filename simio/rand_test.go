package simio

import (
	"bytes"
	"testing"
)

func TestCryptoRandReads(t *testing.T) {
	r := CryptoRand{}
	buf := make([]byte, 32)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 32 {
		t.Fatalf("expected 32 bytes, got %d", n)
	}
	if bytes.Equal(buf, make([]byte, 32)) {
		t.Fatal("expected non-zero bytes")
	}
}

func TestSeededRandDeterministic(t *testing.T) {
	r1 := NewSeededRand(42)
	r2 := NewSeededRand(42)

	buf1 := make([]byte, 64)
	buf2 := make([]byte, 64)

	r1.Read(buf1)
	r2.Read(buf2)

	if !bytes.Equal(buf1, buf2) {
		t.Fatal("same seed should produce identical output")
	}
}

func TestSeededRandDifferentSeeds(t *testing.T) {
	r1 := NewSeededRand(1)
	r2 := NewSeededRand(2)

	buf1 := make([]byte, 64)
	buf2 := make([]byte, 64)

	r1.Read(buf1)
	r2.Read(buf2)

	if bytes.Equal(buf1, buf2) {
		t.Fatal("different seeds should produce different output")
	}
}
