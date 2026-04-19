package db_test

import (
	"path/filepath"
	"testing"

	"github.com/danmestas/libfossil/db"
	_ "github.com/danmestas/libfossil/internal/testdriver"
)

func TestSeedNobody(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.fossil")
	d, err := db.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer d.Close()

	if err := db.CreateRepoSchema(d); err != nil {
		t.Fatalf("CreateRepoSchema: %v", err)
	}
	if err := db.SeedNobody(d, "oi"); err != nil {
		t.Fatalf("SeedNobody: %v", err)
	}

	var login, cap string
	err = d.QueryRow("SELECT login, cap FROM user WHERE login='nobody'").Scan(&login, &cap)
	if err != nil {
		t.Fatalf("nobody user not found: %v", err)
	}
	if cap != "oi" {
		t.Errorf("cap = %q, want oi", cap)
	}
}
