package cli

import "fmt"

// RepoConfigCmd groups configuration operations.
type RepoConfigCmd struct {
	Ls  RepoConfigLsCmd  `cmd:"" help:"List all config entries"`
	Get RepoConfigGetCmd `cmd:"" help:"Get a config value"`
	Set RepoConfigSetCmd `cmd:"" help:"Set a config value"`
}

// RepoConfigLsCmd lists all config entries.
type RepoConfigLsCmd struct{}

func (c *RepoConfigLsCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	db := r.Inner().DB()
	rows, err := db.Query("SELECT name, value FROM config ORDER BY name")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, value string
		rows.Scan(&name, &value)
		fmt.Printf("%-20s %s\n", name, value)
	}
	return rows.Err()
}

// RepoConfigGetCmd gets a single config value.
type RepoConfigGetCmd struct {
	Key string `arg:"" help:"Config key to get"`
}

func (c *RepoConfigGetCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	val, err := r.Config(c.Key)
	if err != nil {
		return fmt.Errorf("config key %q not found", c.Key)
	}
	fmt.Println(val)
	return nil
}

// RepoConfigSetCmd sets a config value.
type RepoConfigSetCmd struct {
	Key   string `arg:"" help:"Config key"`
	Value string `arg:"" help:"Config value"`
}

func (c *RepoConfigSetCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	if err := r.SetConfig(c.Key, c.Value); err != nil {
		return err
	}
	fmt.Printf("%s = %s\n", c.Key, c.Value)
	return nil
}
