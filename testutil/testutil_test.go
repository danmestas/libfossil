package testutil

import (
	"os"
	"testing"
)

func TestNewTestRepo(t *testing.T) {
	tr := NewTestRepo(t)
	if _, err := os.Stat(tr.Path); err != nil {
		t.Fatalf("repo file does not exist: %v", err)
	}
}

func TestFossilRebuild(t *testing.T) {
	tr := NewTestRepo(t)
	tr.FossilRebuild(t)
}

func TestFossilSQL(t *testing.T) {
	tr := NewTestRepo(t)
	out := tr.FossilSQL(t, "SELECT count(*) FROM blob;")
	if out != "1" {
		t.Fatalf("FossilSQL count(*) = %q, want %q", out, "1")
	}
}

func TestFossilBinary(t *testing.T) {
	path := FossilBinary()
	if path == "" {
		t.Skip("fossil binary not found in PATH")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("fossil binary not found at %q: %v", path, err)
	}
}
