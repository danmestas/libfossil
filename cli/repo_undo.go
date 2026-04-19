package cli

// RepoUndoCmd undoes the last checkout operation.
type RepoUndoCmd struct {
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoUndoCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()
	return r.Undo(c.Dir)
}

// RepoRedoCmd re-applies the last undone operation.
type RepoRedoCmd struct {
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoRedoCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()
	return r.Redo(c.Dir)
}
