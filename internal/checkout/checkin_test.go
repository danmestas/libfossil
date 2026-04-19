package checkout

import (
	"errors"
	"fmt"
	"testing"

	"github.com/danmestas/libfossil/simio"
)

func TestCommitFullCycle(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer co.Close()

	rid1, uuid1, _ := co.Version()

	// Extract to MemStorage
	mem := simio.NewMemStorage()
	co.env = &simio.Env{Storage: mem, Clock: simio.RealClock{}, Rand: simio.CryptoRand{}}
	co.dir = "/checkout"
	if err := co.Extract(rid1, ExtractOpts{}); err != nil {
		t.Fatal(err)
	}

	// Modify a file
	if err := mem.WriteFile("/checkout/hello.txt", []byte("modified content\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Commit
	rid2, uuid2, err := co.Commit(CommitOpts{
		Message: "edit hello.txt",
		User:    "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	// New version should be different
	if rid2 == rid1 {
		t.Fatal("new RID should differ from old")
	}
	if uuid2 == uuid1 {
		t.Fatal("new UUID should differ from old")
	}
	if rid2 <= 0 {
		t.Fatal("new RID should be positive")
	}

	// Checkout version should be updated
	currentRID, currentUUID, _ := co.Version()
	if currentRID != rid2 || currentUUID != uuid2 {
		t.Fatal("checkout version not updated after commit")
	}

	// Verify the new checkin has correct parent
	var parentRID int64
	err = r.DB().QueryRow("SELECT pid FROM plink WHERE cid=? AND isprim=1", int64(rid2)).Scan(&parentRID)
	if err != nil {
		t.Fatalf("query parent: %v", err)
	}
	if parentRID != int64(rid1) {
		t.Fatalf("parent = %d, want %d", parentRID, rid1)
	}
}

func TestCommitWithEnqueue(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer co.Close()

	rid1, _, _ := co.Version()
	mem := simio.NewMemStorage()
	co.env = &simio.Env{Storage: mem, Clock: simio.RealClock{}, Rand: simio.CryptoRand{}}
	co.dir = "/checkout"
	if err := co.Extract(rid1, ExtractOpts{}); err != nil {
		t.Fatal(err)
	}

	// Modify two files
	if err := mem.WriteFile("/checkout/hello.txt", []byte("changed1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := mem.WriteFile("/checkout/README.md", []byte("changed2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Only enqueue one
	if err := co.Enqueue(EnqueueOpts{Paths: []string{"hello.txt"}}); err != nil {
		t.Fatal(err)
	}

	// Commit should only include hello.txt changes (but all files in manifest)
	rid2, _, err := co.Commit(CommitOpts{Message: "partial", User: "test"})
	if err != nil {
		t.Fatal(err)
	}

	// Verify that README.md still has the old content in the new checkin
	// (because it wasn't enqueued)
	// We'll check by reading the new manifest and verifying README.md content
	// matches the original
	var readmeUUID string
	err = r.DB().QueryRow(`
		SELECT b.uuid FROM mlink m
		JOIN filename f ON f.fnid = m.fnid
		JOIN blob b ON b.rid = m.fid
		WHERE m.mid = ? AND f.name = 'README.md'
	`, rid2).Scan(&readmeUUID)
	if err != nil {
		t.Fatalf("query README.md uuid in new checkin: %v", err)
	}

	// The original README.md hash
	var originalReadmeUUID string
	err = r.DB().QueryRow(`
		SELECT b.uuid FROM mlink m
		JOIN filename f ON f.fnid = m.fnid
		JOIN blob b ON b.rid = m.fid
		WHERE m.mid = ? AND f.name = 'README.md'
	`, rid1).Scan(&originalReadmeUUID)
	if err != nil {
		t.Fatalf("query original README.md uuid: %v", err)
	}

	if readmeUUID != originalReadmeUUID {
		t.Fatalf("README.md should be unchanged in commit (enqueue filtered it out), got %s, want %s", readmeUUID, originalReadmeUUID)
	}
}

func TestCommitNoChanges(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer co.Close()

	rid, _, _ := co.Version()
	mem := simio.NewMemStorage()
	co.env = &simio.Env{
		Storage: mem, Clock: simio.RealClock{}, Rand: simio.CryptoRand{},
	}
	co.dir = "/checkout"
	if err := co.Extract(rid, ExtractOpts{Force: true}); err != nil {
		t.Fatal(err)
	}

	// Commit with no modifications — manifest.Checkin panics with
	// empty Files list, so buildCommitFiles returns the full
	// unchanged set. The commit succeeds but produces an identical
	// manifest. This documents current behavior.
	_, _, err = co.Commit(CommitOpts{
		Message: "empty commit", User: "test",
	})
	// The commit should succeed (all files are still present,
	// just unchanged). If the project later decides to reject
	// no-change commits, this test should be updated.
	if err != nil {
		t.Fatalf("commit with no changes failed: %v", err)
	}
}

func TestCommitWithDelete(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer co.Close()

	rid1, _, _ := co.Version()
	mem := simio.NewMemStorage()
	co.env = &simio.Env{
		Storage: mem, Clock: simio.RealClock{}, Rand: simio.CryptoRand{},
	}
	co.dir = "/checkout"
	if err := co.Extract(rid1, ExtractOpts{Force: true}); err != nil {
		t.Fatal(err)
	}

	// Unmanage hello.txt (marks deleted=1)
	if err := co.Unmanage(UnmanageOpts{
		Paths: []string{"hello.txt"},
	}); err != nil {
		t.Fatal(err)
	}

	// Commit — hello.txt should not appear in the new manifest.
	rid2, _, err := co.Commit(CommitOpts{
		Message: "delete hello.txt", User: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Verify hello.txt is absent from the new checkin's mlink entries.
	var count int
	err = r.DB().QueryRow(`
		SELECT count(*) FROM mlink m
		JOIN filename f ON f.fnid = m.fnid
		WHERE m.mid = ? AND f.name = 'hello.txt'
	`, rid2).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf(
			"hello.txt should be absent from new checkin, got %d mlink rows",
			count,
		)
	}
}

func TestDequeueAll(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer co.Close()

	if err := co.Enqueue(EnqueueOpts{Paths: []string{"a", "b", "c"}}); err != nil {
		t.Fatal(err)
	}

	enqueued, _ := co.IsEnqueued("a")
	if !enqueued {
		t.Fatal("a should be enqueued")
	}

	// Dequeue all (empty paths = set queue to nil)
	if err := co.Dequeue(DequeueOpts{}); err != nil {
		t.Fatal(err)
	}

	// After dequeue all (nil queue), IsEnqueued returns true (implicit all)
	enqueued, _ = co.IsEnqueued("a")
	if !enqueued {
		t.Fatal("nil queue should mean all implicitly enqueued")
	}
}

func TestEnqueueDequeue(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer co.Close()

	// Initially, empty queue means all implicitly enqueued
	enqueued, _ := co.IsEnqueued("hello.txt")
	if !enqueued {
		t.Fatal("empty queue should return true for any file")
	}

	// Enqueue one file
	if err := co.Enqueue(EnqueueOpts{Paths: []string{"hello.txt"}}); err != nil {
		t.Fatal(err)
	}

	// Now only hello.txt is enqueued
	enqueued, _ = co.IsEnqueued("hello.txt")
	if !enqueued {
		t.Fatal("hello.txt should be enqueued")
	}

	enqueued, _ = co.IsEnqueued("README.md")
	if enqueued {
		t.Fatal("README.md should NOT be enqueued")
	}

	// Dequeue hello.txt
	if err := co.Dequeue(DequeueOpts{Paths: []string{"hello.txt"}}); err != nil {
		t.Fatal(err)
	}

	enqueued, _ = co.IsEnqueued("hello.txt")
	if enqueued {
		t.Fatal("hello.txt should NOT be enqueued after dequeue")
	}
}

func TestDiscardQueue(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatal(err)
	}
	defer co.Close()

	// Enqueue some files
	if err := co.Enqueue(EnqueueOpts{Paths: []string{"a", "b"}}); err != nil {
		t.Fatal(err)
	}

	enqueued, _ := co.IsEnqueued("c")
	if enqueued {
		t.Fatal("c should not be enqueued")
	}

	// Discard queue
	if err := co.DiscardQueue(); err != nil {
		t.Fatal(err)
	}

	// Now all implicitly enqueued
	enqueued, _ = co.IsEnqueued("c")
	if !enqueued {
		t.Fatal("after DiscardQueue, all files should be implicitly enqueued")
	}
}

func TestPreCommitCheck_Abort(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	defer co.Close()

	rid1, _, _ := co.Version()
	if err := co.Extract(rid1, ExtractOpts{}); err != nil {
		t.Fatalf("Extract: %v", err)
	}
	co.env.Storage.WriteFile(co.dir+"/abort.txt", []byte("abort"), 0644)

	wantErr := fmt.Errorf("blocked by policy")
	_, _, err = co.Commit(CommitOpts{
		Message:        "should not commit",
		User:           "test",
		PreCommitCheck: func() error { return wantErr },
	})
	if err == nil {
		t.Fatal("Commit succeeded, want error from PreCommitCheck")
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want wrapping %v", err, wantErr)
	}
}

func TestPreCommitCheck_Nil(t *testing.T) {
	r, cleanup := newTestRepoWithCheckin(t)
	defer cleanup()

	dir := t.TempDir()
	co, err := Create(r, dir, CreateOpts{})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	defer co.Close()

	rid1, _, _ := co.Version()
	if err := co.Extract(rid1, ExtractOpts{}); err != nil {
		t.Fatalf("Extract: %v", err)
	}
	co.env.Storage.WriteFile(co.dir+"/ok.txt", []byte("ok"), 0644)

	_, _, err = co.Commit(CommitOpts{
		Message: "normal commit",
		User:    "test",
	})
	if err != nil {
		t.Fatalf("Commit with nil PreCommitCheck failed: %v", err)
	}
}
