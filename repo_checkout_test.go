package libfossil

import (
	"path/filepath"
	"testing"

	_ "github.com/danmestas/libfossil/internal/testdriver"
)

func TestCommitAndTimeline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.fossil")

	r, err := Create(path, CreateOpts{User: "test"})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	rid1, uuid, err := r.Commit(CommitOpts{
		Files: []FileToCommit{
			{Name: "hello.txt", Content: []byte("hello world\n")},
		},
		Comment: "initial commit",
		User:    "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if uuid == "" {
		t.Fatal("Commit returned empty UUID")
	}
	if len(uuid) != 40 && len(uuid) != 64 {
		t.Errorf("UUID length = %d, want 40 or 64", len(uuid))
	}
	if rid1 <= 0 {
		t.Fatalf("Commit returned rid=%d, want > 0", rid1)
	}

	// Second commit to test timeline walking.
	rid2, uuid2, err := r.Commit(CommitOpts{
		ParentID: rid1,
		Files: []FileToCommit{
			{Name: "hello.txt", Content: []byte("hello world v2\n")},
		},
		Comment: "second commit",
		User:    "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	entries, err := r.Timeline(LogOpts{Start: rid2, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Fatalf("Timeline returned %d entries, want >= 2", len(entries))
	}
	if entries[0].UUID != uuid2 {
		t.Errorf("first entry UUID = %q, want %q", entries[0].UUID, uuid2)
	}
	if entries[0].Comment != "second commit" {
		t.Errorf("first entry Comment = %q, want %q", entries[0].Comment, "second commit")
	}
	if entries[1].UUID != uuid {
		t.Errorf("second entry UUID = %q, want %q", entries[1].UUID, uuid)
	}
}

func TestListFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.fossil")

	r, err := Create(path, CreateOpts{User: "test"})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	rid, _, err := r.Commit(CommitOpts{
		Files: []FileToCommit{
			{Name: "a.txt", Content: []byte("aaa\n")},
			{Name: "b.txt", Content: []byte("bbb\n")},
		},
		Comment: "two files",
		User:    "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	files, err := r.ListFiles(rid)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("ListFiles returned %d files, want 2", len(files))
	}
	if files[0].Name != "a.txt" || files[1].Name != "b.txt" {
		t.Errorf("files = %v, want a.txt and b.txt", files)
	}
}
