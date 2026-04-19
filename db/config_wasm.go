//go:build wasip1 || js

package db

// wasmPragmaOverrides returns browser/WASI-specific pragma overrides.
// OPFS doesn't support shared memory, so WAL is replaced with DELETE journal
// and EXCLUSIVE locking mode (single connection per Worker).
func wasmPragmaOverrides() map[string]string {
	return map[string]string{
		"journal_mode": "DELETE",
		"locking_mode": "EXCLUSIVE",
	}
}

// wasmDSNSuffix appends nolock=1 to disable file locking on WASM.
// WASI runtimes don't support POSIX/OFD/BSD file locks.
func wasmDSNSuffix() string {
	return "nolock=1"
}

// wasmClearPragmas signals that WASM builds should not skip pragmas —
// browser builds need DELETE journal + EXCLUSIVE locking.
const wasmClearPragmas = false
