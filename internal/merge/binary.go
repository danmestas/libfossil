package merge

func init() { Register(&Binary{}) }

// Binary always returns a conflict. Binary files cannot be auto-merged.
type Binary struct{}

func (b *Binary) Name() string { return "binary" }

func (b *Binary) Merge(base, local, remote []byte) (*Result, error) {
	return &Result{
		Content: local,
		Clean:   false,
		Conflicts: []Conflict{{
			Local:  local,
			Remote: remote,
			Base:   base,
		}},
	}, nil
}
