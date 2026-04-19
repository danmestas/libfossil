package cli

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// RepoOpenCmd opens a checkout in a directory, creating the .fslckout database.
type RepoOpenCmd struct {
	Dir string `arg:"" optional:"" help:"Checkout directory (default: current dir)" default:"."`
}

func (c *RepoOpenCmd) Run(g *Globals) error {
	if g.Repo == "" {
		return fmt.Errorf("repository required (use -R <path>)")
	}

	absRepo, err := filepath.Abs(g.Repo)
	if err != nil {
		return err
	}

	ckoutPath := filepath.Join(c.Dir, ".fslckout")
	if _, err := os.Stat(ckoutPath); err == nil {
		return fmt.Errorf("checkout already exists: %s", ckoutPath)
	}

	// Create checkout database with vfile, vvar, vmerge tables.
	db, err := sql.Open("sqlite", ckoutPath)
	if err != nil {
		return err
	}
	defer db.Close()

	schema := `
		CREATE TABLE vvar(
			name TEXT PRIMARY KEY NOT NULL,
			value CLOB,
			CHECK(typeof(name)='text' AND length(name)>=1)
		) WITHOUT ROWID;
		CREATE TABLE vfile(
			id INTEGER PRIMARY KEY,
			vid INTEGER,
			chnged INT DEFAULT 0,
			deleted BOOLEAN DEFAULT 0,
			isexe BOOLEAN,
			islink BOOLEAN,
			rid INTEGER,
			mrid INTEGER,
			mtime INTEGER,
			pathname TEXT,
			origname TEXT,
			mhash TEXT,
			UNIQUE(pathname,vid)
		);
		CREATE TABLE vmerge(
			id INTEGER REFERENCES vfile,
			merge INTEGER,
			mhash TEXT
		);
	`
	if _, err := db.Exec(schema); err != nil {
		os.Remove(ckoutPath)
		return fmt.Errorf("creating checkout schema: %w", err)
	}

	// Resolve tip checkin from repo.
	r, err := g.OpenRepo()
	if err != nil {
		os.Remove(ckoutPath)
		return err
	}

	tipRid, err := resolveRID(r, "tip")
	if err != nil {
		tipRid = 0
	}
	var tipHash string
	if tipRid > 0 {
		r.Inner().DB().QueryRow("SELECT uuid FROM blob WHERE rid=?", tipRid).Scan(&tipHash)
	}

	// Populate vvar with checkout metadata.
	vvars := map[string]string{
		"repository":     absRepo,
		"checkout":       fmt.Sprintf("%d", tipRid),
		"checkout-hash":  tipHash,
		"undo_available": "0",
		"undo_checkout":  "0",
	}
	for k, v := range vvars {
		if _, err := db.Exec("INSERT INTO vvar(name,value) VALUES(?,?)", k, v); err != nil {
			return fmt.Errorf("setting vvar %s: %w", k, err)
		}
	}

	// Populate vfile from tip manifest.
	if tipRid > 0 {
		files, err := r.ListFiles(tipRid)
		if err == nil {
			idb := r.Inner().DB()
			for _, f := range files {
				isExe := f.Perm == "x"
				var fileRid int64
				idb.QueryRow("SELECT rid FROM blob WHERE uuid=?", f.UUID).Scan(&fileRid)
				db.Exec(`INSERT INTO vfile(vid, chnged, deleted, isexe, islink, rid, mrid, pathname, mhash)
					VALUES(?, 0, 0, ?, 0, ?, ?, ?, ?)`,
					tipRid, isExe, fileRid, fileRid, f.Name, f.UUID)
			}
		}
	}

	r.Close()

	fmt.Printf("opened checkout in %s (repo: %s)\n", c.Dir, absRepo)
	if tipRid > 0 {
		fmt.Printf("checked out version %s\n", tipHash[:10])
	} else {
		fmt.Println("empty repository -- no checkins yet")
	}
	return nil
}
