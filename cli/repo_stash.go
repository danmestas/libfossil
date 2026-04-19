package cli

import "fmt"

// RepoStashCmd groups stash operations.
type RepoStashCmd struct {
	Save  RepoStashSaveCmd  `cmd:"" help:"Stash working changes"`
	Pop   RepoStashPopCmd   `cmd:"" help:"Apply top stash and drop it"`
	Apply RepoStashApplyCmd `cmd:"" help:"Apply stash without dropping"`
	Ls    RepoStashLsCmd    `cmd:"" help:"List stash entries"`
	Drop  RepoStashDropCmd  `cmd:"" help:"Remove stash entry"`
	Clear RepoStashClearCmd `cmd:"" help:"Remove all stash entries"`
}

// RepoStashSaveCmd saves working changes to the stash.
type RepoStashSaveCmd struct {
	Message string `short:"m" help:"Stash message" default:""`
	Dir     string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoStashSaveCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()
	return r.StashSave(c.Dir, c.Message)
}

// RepoStashPopCmd pops the top stash entry.
type RepoStashPopCmd struct {
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoStashPopCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()
	return r.StashPop(c.Dir)
}

// RepoStashApplyCmd applies a stash entry without removing it.
type RepoStashApplyCmd struct {
	ID  int64  `arg:"" optional:"" help:"Stash ID to apply (default: latest)"`
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoStashApplyCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()
	return r.StashApply(c.Dir, c.ID)
}

// RepoStashLsCmd lists stash entries.
type RepoStashLsCmd struct{}

func (c *RepoStashLsCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	entries, err := r.StashList()
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("no stash entries")
		return nil
	}
	for _, e := range entries {
		fmt.Printf("%3d: %s  %s\n", e.ID, e.Time, e.Comment)
	}
	return nil
}

// RepoStashDropCmd removes a stash entry by ID.
type RepoStashDropCmd struct {
	ID int64 `arg:"" required:"" help:"Stash ID to drop"`
}

func (c *RepoStashDropCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()
	return r.StashDrop(c.ID)
}

// RepoStashClearCmd removes all stash entries.
type RepoStashClearCmd struct{}

func (c *RepoStashClearCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()
	return r.StashClear()
}
