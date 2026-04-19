package cli

import "fmt"

// RepoBisectCmd groups binary search operations.
type RepoBisectCmd struct {
	Good   RepoBisectGoodCmd   `cmd:"" help:"Mark version as good"`
	Bad    RepoBisectBadCmd    `cmd:"" help:"Mark version as bad"`
	Next   RepoBisectNextCmd   `cmd:"" help:"Check out midpoint version"`
	Skip   RepoBisectSkipCmd   `cmd:"" help:"Skip current version"`
	Reset  RepoBisectResetCmd  `cmd:"" help:"Clear bisect state"`
	Ls     RepoBisectLsCmd     `cmd:"" help:"Show bisect path"`
	Status RepoBisectStatusCmd `cmd:"" help:"Show bisect state"`
}

// RepoBisectGoodCmd marks a version as good.
type RepoBisectGoodCmd struct {
	Version string `arg:"" optional:"" help:"Version to mark (default: current checkout)"`
	Dir     string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoBisectGoodCmd) Run(g *Globals) error {
	return fmt.Errorf("bisect: not yet implemented (requires checkout DB)")
}

// RepoBisectBadCmd marks a version as bad.
type RepoBisectBadCmd struct {
	Version string `arg:"" optional:"" help:"Version to mark (default: current checkout)"`
	Dir     string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoBisectBadCmd) Run(g *Globals) error {
	return fmt.Errorf("bisect: not yet implemented (requires checkout DB)")
}

// RepoBisectSkipCmd skips the current version.
type RepoBisectSkipCmd struct {
	Version string `arg:"" optional:"" help:"Version to skip (default: current checkout)"`
	Dir     string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoBisectSkipCmd) Run(g *Globals) error {
	return fmt.Errorf("bisect: not yet implemented (requires checkout DB)")
}

// RepoBisectNextCmd checks out the bisect midpoint version.
type RepoBisectNextCmd struct {
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoBisectNextCmd) Run(g *Globals) error {
	return fmt.Errorf("bisect: not yet implemented (requires checkout DB)")
}

// RepoBisectResetCmd clears bisect state.
type RepoBisectResetCmd struct {
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoBisectResetCmd) Run(g *Globals) error {
	return fmt.Errorf("bisect: not yet implemented (requires checkout DB)")
}

// RepoBisectLsCmd shows the bisect path.
type RepoBisectLsCmd struct {
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoBisectLsCmd) Run(g *Globals) error {
	return fmt.Errorf("bisect: not yet implemented (requires checkout DB)")
}

// RepoBisectStatusCmd shows bisect state.
type RepoBisectStatusCmd struct {
	Dir string `short:"d" help:"Checkout directory" default:"."`
}

func (c *RepoBisectStatusCmd) Run(g *Globals) error {
	return fmt.Errorf("bisect: not yet implemented (requires checkout DB)")
}
