package libfossil

import (
	"fmt"
	"time"

	"github.com/danmestas/libfossil/internal/annotate"
	"github.com/danmestas/libfossil/internal/fsltype"
)

// LogOpts configures a log/timeline query.
type LogOpts struct {
	Start int64
	Limit int
}

// LogEntry represents a single checkin in the timeline.
type LogEntry struct {
	RID     int64
	UUID    string
	Comment string
	User    string
	Time    time.Time
	Parents []string
}

// DiffEntry describes a unified diff for a single file.
type DiffEntry struct {
	Name    string
	Unified string
}

// AnnotatedLine is a single line of blame/annotate output.
type AnnotatedLine struct {
	Text string
	UUID string
	User string
	Date time.Time
}

// AnnotateOpts configures an annotate operation.
type AnnotateOpts struct {
	FilePath string
	StartRID int64
}

// StatusOpts configures a working-tree status query.
type StatusOpts struct {
	Dir string
}

// BisectSession holds state for a binary-search bisect operation.
type BisectSession struct {
	inner interface{}
}

// Annotate attributes each line of a file to the commit that last changed it.
func (r *Repo) Annotate(opts AnnotateOpts) ([]AnnotatedLine, error) {
	lines, err := annotate.Annotate(r.inner, annotate.Options{
		FilePath: opts.FilePath,
		StartRID: fsltype.FslID(opts.StartRID),
	})
	if err != nil {
		return nil, fmt.Errorf("libfossil: annotate: %w", err)
	}
	result := make([]AnnotatedLine, len(lines))
	for i, l := range lines {
		result[i] = AnnotatedLine{
			Text: l.Text,
			UUID: l.Version.UUID,
			User: l.Version.User,
			Date: l.Version.Date,
		}
	}
	return result, nil
}

// Diff returns a unified diff for a file between two checkins.
// TODO: wire to internal diff package when available.
func (r *Repo) Diff(ridA, ridB int64, filePath string) ([]DiffEntry, error) {
	return nil, fmt.Errorf("libfossil: diff: not yet implemented")
}
