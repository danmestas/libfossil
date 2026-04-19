package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type TestRepo struct {
	Path string
	Dir  string
}

func FossilBinary() string {
	if bin := os.Getenv("FOSSIL_BIN"); bin != "" {
		return bin
	}
	path, err := exec.LookPath("fossil")
	if err != nil {
		return ""
	}
	return path
}

func NewTestRepo(t *testing.T) *TestRepo {
	t.Helper()
	bin := FossilBinary()
	if bin == "" {
		t.Skip("fossil binary not found")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "test.fossil")
	cmd := exec.Command(bin, "new", path)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fossil new failed: %v\n%s", err, out)
	}
	return &TestRepo{Path: path, Dir: dir}
}

func NewTestRepoFromPath(t *testing.T, path string) *TestRepo {
	t.Helper()
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("cannot resolve path %q: %v", path, err)
	}
	return &TestRepo{Path: abs, Dir: filepath.Dir(abs)}
}

func (r *TestRepo) FossilRebuild(t *testing.T) {
	t.Helper()
	cmd := exec.Command(FossilBinary(), "rebuild", r.Path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fossil rebuild failed: %v\n%s", err, out)
	}
}

func (r *TestRepo) FossilArtifact(t *testing.T, uuid string) []byte {
	t.Helper()
	cmd := exec.Command(FossilBinary(), "artifact", uuid, "-R", r.Path)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("fossil artifact %s failed: %v", uuid, err)
	}
	return out
}

func HasFossil() bool {
	return FossilBinary() != ""
}

func FossilRebuild(repoPath string) error {
	cmd := exec.Command(FossilBinary(), "rebuild", repoPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fossil rebuild: %v\n%s", err, out)
	}
	return nil
}

func FossilTimeline(repoPath string) ([]byte, error) {
	cmd := exec.Command(FossilBinary(), "timeline", "-n", "100", "-R", repoPath)
	return cmd.Output()
}

func FossilArtifactByPath(repoPath, uuid string) ([]byte, error) {
	cmd := exec.Command(FossilBinary(), "artifact", uuid, "-R", repoPath)
	return cmd.Output()
}

func (r *TestRepo) FossilSQL(t *testing.T, sql string) string {
	t.Helper()
	cmd := exec.Command(FossilBinary(), "sql", "-R", r.Path, sql)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Fatalf("fossil sql failed: %v\nstderr: %s", err, exitErr.Stderr)
		}
		t.Fatalf("fossil sql failed: %v", err)
	}
	return strings.TrimSpace(string(out))
}
