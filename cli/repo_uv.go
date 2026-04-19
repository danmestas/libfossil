package cli

import (
	"fmt"
	"os"
	"time"
)

// RepoUVCmd groups unversioned file operations.
type RepoUVCmd struct {
	Ls     RepoUVLsCmd     `cmd:"" help:"List unversioned files"`
	Put    RepoUVPutCmd    `cmd:"" help:"Add or update an unversioned file"`
	Get    RepoUVGetCmd    `cmd:"" help:"Retrieve an unversioned file"`
	Delete RepoUVDeleteCmd `cmd:"" help:"Delete an unversioned file (creates tombstone)"`
}

// RepoUVLsCmd lists unversioned files.
type RepoUVLsCmd struct{}

func (c *RepoUVLsCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	entries, err := r.UVList()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("(no unversioned files)")
		return nil
	}

	for _, e := range entries {
		ts := e.Mtime.UTC().Format("2006-01-02 15:04:05")
		if e.Hash == "" {
			fmt.Printf("%-40s %s  [deleted]\n", e.Name, ts)
		} else {
			hash := e.Hash
			if len(hash) > 10 {
				hash = hash[:10]
			}
			fmt.Printf("%-40s %s  %6d  %s\n", e.Name, ts, e.Size, hash)
		}
	}
	return nil
}

// RepoUVPutCmd adds or updates an unversioned file.
type RepoUVPutCmd struct {
	Name string `arg:"" help:"Name of the unversioned file"`
	File string `arg:"" help:"Local file to upload"`
}

func (c *RepoUVPutCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	data, err := os.ReadFile(c.File)
	if err != nil {
		return fmt.Errorf("read %s: %w", c.File, err)
	}

	if err := r.UVWrite(c.Name, data, time.Now()); err != nil {
		return err
	}

	fmt.Printf("wrote %s (%d bytes)\n", c.Name, len(data))
	return nil
}

// RepoUVGetCmd retrieves an unversioned file.
type RepoUVGetCmd struct {
	Name   string `arg:"" help:"Name of the unversioned file"`
	Output string `short:"o" help:"Output file (default: stdout)"`
}

func (c *RepoUVGetCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	data, _, hash, err := r.UVRead(c.Name)
	if err != nil {
		return err
	}
	if data == nil && hash == "" {
		return fmt.Errorf("unversioned file %q not found", c.Name)
	}
	if hash == "" {
		return fmt.Errorf("unversioned file %q has been deleted", c.Name)
	}

	if c.Output != "" {
		return os.WriteFile(c.Output, data, 0o644)
	}
	_, err = os.Stdout.Write(data)
	return err
}

// RepoUVDeleteCmd deletes an unversioned file (creates a tombstone).
type RepoUVDeleteCmd struct {
	Name string `arg:"" help:"Name of the unversioned file to delete"`
}

func (c *RepoUVDeleteCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	// UV delete is done by writing a zero-length entry.
	if err := r.UVWrite(c.Name, nil, time.Now()); err != nil {
		return err
	}

	fmt.Printf("deleted %s\n", c.Name)
	return nil
}
