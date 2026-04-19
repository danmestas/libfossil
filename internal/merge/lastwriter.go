package merge

func init() { Register(&LastWriterWins{}) }

// LastWriterWins picks the newer version. The caller passes the newer
// content as remote. Always clean — no conflicts possible.
type LastWriterWins struct{}

func (l *LastWriterWins) Name() string { return "last-writer-wins" }

func (l *LastWriterWins) Merge(base, local, remote []byte) (*Result, error) {
	return &Result{Content: remote, Clean: true}, nil
}
