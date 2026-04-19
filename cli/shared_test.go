package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	libfossil "github.com/danmestas/libfossil"
	"github.com/danmestas/libfossil/cli"

	_ "github.com/danmestas/libfossil/internal/testdriver"
)

func TestGlobalsOpenRepo(t *testing.T) {
	tmp := t.TempDir()
	repoPath := filepath.Join(tmp, "test.fossil")

	r, err := libfossil.Create(repoPath, libfossil.CreateOpts{User: "test"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	r.Close()

	g := &cli.Globals{Repo: repoPath}
	opened, err := g.OpenRepo()
	if err != nil {
		t.Fatalf("OpenRepo: %v", err)
	}
	defer opened.Close()

	if opened.Path() != repoPath {
		t.Errorf("Path() = %q, want %q", opened.Path(), repoPath)
	}
}

func TestGlobalsOpenRepoNotFound(t *testing.T) {
	g := &cli.Globals{Repo: "/nonexistent/repo.fossil"}
	_, err := g.OpenRepo()
	if err == nil {
		t.Fatal("expected error for nonexistent repo")
	}
}

func TestGlobalsOpenRepoAutoFind(t *testing.T) {
	tmp := t.TempDir()
	repoPath := filepath.Join(tmp, "auto.fossil")

	r, err := libfossil.Create(repoPath, libfossil.CreateOpts{User: "test"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	r.Close()

	// Change to the temp dir so findRepo can discover the .fossil file.
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(tmp)

	g := &cli.Globals{}
	opened, err := g.OpenRepo()
	if err != nil {
		t.Fatalf("OpenRepo auto-find: %v", err)
	}
	defer opened.Close()
}
