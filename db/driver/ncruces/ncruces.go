//go:build !js

package ncruces

import (
	"fmt"
	"strings"

	"github.com/danmestas/libfossil/db"
	_ "github.com/ncruces/go-sqlite3/driver"
)

func init() {
	db.Register(db.DriverConfig{
		Name:     "sqlite3",
		BuildDSN: buildDSN,
	})
}

func buildDSN(path string, pragmas map[string]string) string {
	if path == "" {
		panic("ncruces.buildDSN: path must not be empty")
	}
	if len(pragmas) == 0 {
		return path
	}
	var parts []string
	for k, v := range pragmas {
		parts = append(parts, fmt.Sprintf("_pragma=%s(%s)", k, v))
	}
	return fmt.Sprintf("file:%s?%s", path, strings.Join(parts, "&"))
}
