package cli

import (
	"fmt"
	"strings"
)

// RepoQueryCmd executes raw SQL against the repository database.
type RepoQueryCmd struct {
	SQL string `arg:"" help:"SQL query to execute"`
}

func (c *RepoQueryCmd) Run(g *Globals) error {
	r, err := g.OpenRepo()
	if err != nil {
		return err
	}
	defer r.Close()

	db := r.Inner().DB()
	sql := strings.TrimSpace(c.SQL)

	// SELECT/PRAGMA queries return rows; everything else returns affected count.
	if strings.HasPrefix(strings.ToUpper(sql), "SELECT") || strings.HasPrefix(strings.ToUpper(sql), "PRAGMA") {
		rows, err := db.Query(sql)
		if err != nil {
			return err
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			return err
		}

		// Print header.
		fmt.Println(strings.Join(cols, "|"))

		// Print rows.
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		for rows.Next() {
			rows.Scan(ptrs...)
			parts := make([]string, len(cols))
			for i, v := range vals {
				parts[i] = fmt.Sprintf("%v", v)
			}
			fmt.Println(strings.Join(parts, "|"))
		}
		return rows.Err()
	}

	result, err := db.Exec(sql)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	fmt.Printf("%d rows affected\n", affected)
	return nil
}
