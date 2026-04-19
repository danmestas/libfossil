package cli

import "fmt"

// RepoRenameCmd renames a tracked file in the checkout.
type RepoRenameCmd struct {
	From string `arg:"" help:"Current file name"`
	To   string `arg:"" help:"New file name"`
	Dir  string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoRenameCmd) Run(g *Globals) error {
	ckout, err := openCheckout(c.Dir)
	if err != nil {
		return err
	}
	defer ckout.Close()

	vid, err := checkoutVid(ckout)
	if err != nil {
		return err
	}

	var id int64
	err = ckout.QueryRow("SELECT id FROM vfile WHERE pathname=? AND vid=?", c.From, vid).Scan(&id)
	if err != nil {
		return fmt.Errorf("%s: not tracked in checkout", c.From)
	}

	var existing int64
	if ckout.QueryRow("SELECT id FROM vfile WHERE pathname=? AND vid=?", c.To, vid).Scan(&existing) == nil {
		return fmt.Errorf("%s: already tracked in checkout", c.To)
	}

	_, err = ckout.Exec("UPDATE vfile SET pathname=?, origname=?, chnged=1 WHERE id=?", c.To, c.From, id)
	if err != nil {
		return err
	}
	fmt.Printf("RENAMED  %s -> %s\n", c.From, c.To)
	return nil
}
