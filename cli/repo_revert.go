package cli

import "fmt"

// RepoRevertCmd undoes staging changes in the checkout.
type RepoRevertCmd struct {
	Files []string `arg:"" optional:"" help:"Files to revert (default: all)"`
	Dir   string   `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoRevertCmd) Run(g *Globals) error {
	ckout, err := openCheckout(c.Dir)
	if err != nil {
		return err
	}
	defer ckout.Close()

	vid, err := checkoutVid(ckout)
	if err != nil {
		return err
	}

	if len(c.Files) == 0 {
		result, err := ckout.Exec("UPDATE vfile SET chnged=0, deleted=0, origname=NULL WHERE vid=? AND (chnged=1 OR deleted=1)", vid)
		if err != nil {
			return err
		}
		affected, _ := result.RowsAffected()

		removed, err := ckout.Exec("DELETE FROM vfile WHERE vid=? AND rid=0", vid)
		if err != nil {
			return err
		}
		removedCount, _ := removed.RowsAffected()

		fmt.Printf("reverted %d files, removed %d staged additions\n", affected, removedCount)
	} else {
		for _, name := range c.Files {
			var id, rid int64
			err := ckout.QueryRow("SELECT id, rid FROM vfile WHERE pathname=? AND vid=?", name, vid).Scan(&id, &rid)
			if err != nil {
				return fmt.Errorf("%s: not tracked in checkout", name)
			}
			if rid == 0 {
				ckout.Exec("DELETE FROM vfile WHERE id=?", id)
				fmt.Printf("UNSTAGED %s\n", name)
			} else {
				ckout.Exec("UPDATE vfile SET chnged=0, deleted=0, origname=NULL WHERE id=?", id)
				fmt.Printf("REVERTED %s\n", name)
			}
		}
	}
	return nil
}
