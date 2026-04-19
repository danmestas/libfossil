package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

// RepoMergeResolveCmd marks a file conflict as resolved.
type RepoMergeResolveCmd struct {
	File string `arg:"" help:"File to mark as resolved"`
	Dir  string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoMergeResolveCmd) Run(g *Globals) error {
	resolved := false

	// Try standard conflict: reset vfile.chnged from 5 to 1.
	ckout, err := openCheckout(c.Dir)
	if err == nil {
		defer ckout.Close()
		vid, _ := checkoutVid(ckout)
		result, err := ckout.Exec("UPDATE vfile SET chnged=1 WHERE pathname=? AND vid=? AND chnged=5", c.File, vid)
		if err == nil {
			affected, _ := result.RowsAffected()
			if affected > 0 {
				base := filepath.Join(c.Dir, c.File)
				os.Remove(base + ".LOCAL")
				os.Remove(base + ".BASELINE")
				os.Remove(base + ".MERGE")
				fmt.Printf("resolved: %s (conflict markers)\n", c.File)
				resolved = true
			}
		}
	}

	// Try conflict-fork: delete from conflict table.
	r, err := g.OpenRepo()
	if err == nil {
		defer r.Close()
		err := r.ResolveConflictFork(c.File)
		if err == nil {
			fmt.Printf("resolved: %s (conflict-fork)\n", c.File)
			resolved = true
		}
	}

	if !resolved {
		return fmt.Errorf("%s: no conflict found", c.File)
	}
	return nil
}
