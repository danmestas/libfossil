package db

// OpenConfig allows callers to customize driver selection and pragmas.
type OpenConfig struct {
	Driver  string            // override driver name (empty = use registered driver)
	Pragmas map[string]string // additional/override pragmas (merged with defaults)
}

// DefaultPragmas returns the default pragma settings.
func DefaultPragmas() map[string]string {
	m := map[string]string{
		"journal_mode": "WAL",
		"busy_timeout": "5000",
		"foreign_keys": "OFF", // ncruces enables FK by default; normalize to OFF for schema compat
	}
	for k, v := range wasmPragmaOverrides() {
		m[k] = v
	}
	return m
}
