package cli

import "fmt"

// RepoResolveCmd resolves a symbolic name to UUID and RID.
type RepoResolveCmd struct {
	Name string `arg:"" help:"Symbolic name, UUID, or prefix to resolve (e.g. trunk, tip, UUID prefix)"`
}

func (c *RepoResolveCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	rid, err := resolveRID(r, c.Name)
	if err != nil {
		return err
	}

	var uuid string
	r.Inner().DB().QueryRow("SELECT uuid FROM blob WHERE rid=?", rid).Scan(&uuid)

	fmt.Printf("rid:  %d\n", rid)
	fmt.Printf("uuid: %s\n", uuid)
	return nil
}
