package libfossil

import (
	"fmt"

	"github.com/danmestas/libfossil/internal/fsltype"
	"github.com/danmestas/libfossil/internal/merge"
)

// MergeOpts configures a merge operation.
type MergeOpts struct {
	Strategy string
	Dir      string
}

// MergeResult describes the outcome of a merge.
type MergeResult struct {
	Clean     bool
	Conflicts []MergeConflict
}

// MergeConflict describes a conflict region in a file.
type MergeConflict struct {
	Name      string
	StartLine int
	EndLine   int
}

// Fork describes a divergence point between two branches.
type Fork struct {
	Ancestor  int64
	LocalTip  int64
	RemoteTip int64
}

// FindCommonAncestor finds the nearest common ancestor of two checkins.
func (r *Repo) FindCommonAncestor(ridA, ridB int64) (int64, error) {
	ancestor, err := merge.FindCommonAncestor(r.inner, fsltype.FslID(ridA), fsltype.FslID(ridB))
	if err != nil {
		return 0, fmt.Errorf("libfossil: find common ancestor: %w", err)
	}
	return int64(ancestor), nil
}

// DetectForks finds divergent branches in the repository.
func (r *Repo) DetectForks() ([]Fork, error) {
	forks, err := merge.DetectForks(r.inner)
	if err != nil {
		return nil, fmt.Errorf("libfossil: detect forks: %w", err)
	}
	result := make([]Fork, len(forks))
	for i, f := range forks {
		result[i] = Fork{
			Ancestor:  int64(f.Ancestor),
			LocalTip:  int64(f.LocalTip),
			RemoteTip: int64(f.RemoteTip),
		}
	}
	return result, nil
}

// ListConflictForks returns filenames with unresolved conflict-fork entries.
func (r *Repo) ListConflictForks() ([]string, error) {
	names, err := merge.ListConflictForks(r.inner)
	if err != nil {
		return nil, fmt.Errorf("libfossil: list conflict forks: %w", err)
	}
	return names, nil
}

// ResolveConflictFork marks a conflict-fork entry as resolved.
func (r *Repo) ResolveConflictFork(filename string) error {
	return merge.ResolveConflictFork(r.inner, filename)
}
