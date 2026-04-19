package libfossil

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/danmestas/libfossil/internal/testdriver"
)

func TestCreateAndOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.fossil")

	r, err := Create(path, CreateOpts{User: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Path() != path {
		t.Errorf("Path() = %q, want %q", r.Path(), path)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal("repo file should exist after Create")
	}

	r2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer r2.Close()
	if r2.Path() != path {
		t.Errorf("Path() = %q, want %q", r2.Path(), path)
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.fossil")

	r, err := Create(path, CreateOpts{User: "test"})
	if err != nil {
		t.Fatal(err)
	}
	r.Close()

	_, err = Create(path, CreateOpts{User: "test"})
	if err == nil {
		t.Error("Create on existing repo should fail")
	}
}

func TestOpenNotFound(t *testing.T) {
	_, err := Open("/nonexistent/path.fossil")
	if err == nil {
		t.Error("Open on nonexistent path should fail")
	}
}
