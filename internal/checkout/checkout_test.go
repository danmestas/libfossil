package checkout

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateAndVersion(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	c, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer c.Close()

	rid, uuid, err := c.Version()
	if err != nil {
		t.Fatalf("Version failed: %v", err)
	}

	if rid <= 0 {
		t.Errorf("expected positive RID, got %d", rid)
	}
	if uuid == "" {
		t.Errorf("expected non-empty UUID, got empty string")
	}

	t.Logf("Checkout version: RID=%d, UUID=%s", rid, uuid)
}

func TestOpenRoundTrip(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()

	// Create
	c1, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	rid1, uuid1, err := c1.Version()
	if err != nil {
		c1.Close()
		t.Fatalf("Version after Create failed: %v", err)
	}

	if err := c1.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Reopen
	c2, err := Open(r, dir, OpenOpts{})
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer c2.Close()

	rid2, uuid2, err := c2.Version()
	if err != nil {
		t.Fatalf("Version after Open failed: %v", err)
	}

	if rid1 != rid2 {
		t.Errorf("RID mismatch: create=%d, open=%d", rid1, rid2)
	}
	if uuid1 != uuid2 {
		t.Errorf("UUID mismatch: create=%s, open=%s", uuid1, uuid2)
	}
}

func TestValidateFingerprint(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	c, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer c.Close()

	if err := c.ValidateFingerprint(); err != nil {
		t.Errorf("ValidateFingerprint failed: %v", err)
	}
}

func TestValidateFingerprintMismatch(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	c, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer c.Close()

	// Tamper with the checkout-hash vvar.
	if err := setVVar(c.db, "checkout-hash", "deadbeef"); err != nil {
		t.Fatal(err)
	}

	err = c.ValidateFingerprint()
	if err == nil {
		t.Fatal("ValidateFingerprint should fail with tampered hash")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf(
			"error should mention mismatch, got: %v", err,
		)
	}
}

func TestFindCheckoutDB(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()

	// Create a checkout
	c, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	c.Close()

	// FindCheckoutDB should find it in the same directory
	env := c.env
	dbPath, err := FindCheckoutDB(env.Storage, dir, false)
	if err != nil {
		t.Fatalf("FindCheckoutDB failed: %v", err)
	}

	expectedPath := filepath.Join(dir, PreferredDBName())
	if dbPath != expectedPath {
		t.Errorf("expected %s, got %s", expectedPath, dbPath)
	}
}

func TestFindCheckoutDBSearchParents(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	rootDir := t.TempDir()
	subDir := filepath.Join(rootDir, "sub", "deep")

	// Create checkout in root
	c, err := Create(r, rootDir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	c.Close()

	// Create subdirectory
	if err := c.env.Storage.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	// Search from subdirectory with searchParents=true
	dbPath, err := FindCheckoutDB(c.env.Storage, subDir, true)
	if err != nil {
		t.Fatalf("FindCheckoutDB with searchParents failed: %v", err)
	}

	expectedPath := filepath.Join(rootDir, PreferredDBName())
	if dbPath != expectedPath {
		t.Errorf("expected %s, got %s", expectedPath, dbPath)
	}

	// Search without searchParents should fail
	_, err = FindCheckoutDB(c.env.Storage, subDir, false)
	if err == nil {
		t.Errorf("FindCheckoutDB without searchParents should have failed")
	}
}

func TestDir(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	c, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer c.Close()

	if got := c.Dir(); got != dir {
		t.Errorf("Dir: expected %s, got %s", dir, got)
	}
}

func TestRepo(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	c, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer c.Close()

	if got := c.Repo(); got != r {
		t.Errorf("Repo: expected same *repo.Repo instance")
	}
}
