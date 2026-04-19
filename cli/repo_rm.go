package cli

import "fmt"

// RepoRmCmd stages files for removal from the checkout.
type RepoRmCmd struct {
	Files []string `arg:"" required:"" help:"Files to stage for removal"`
	Dir   string   `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoRmCmd) Run(g *Globals) error {
	ckout, err := openCheckout(c.Dir)
	if err != nil {
		return err
	}
	defer ckout.Close()

	vid, err := checkoutVid(ckout)
	if err != nil {
		return err
	}

	for _, name := range c.Files {
		var id int64
		err := ckout.QueryRow("SELECT id FROM vfile WHERE pathname=? AND vid=?", name, vid).Scan(&id)
		if err != nil {
			return fmt.Errorf("%s: not tracked in checkout", name)
		}
		ckout.Exec("UPDATE vfile SET deleted=1 WHERE id=?", id)
		fmt.Printf("REMOVED  %s\n", name)
	}
	return nil
}
