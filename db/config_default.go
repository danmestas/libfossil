//go:build !wasip1 && !js

package db

func wasmPragmaOverrides() map[string]string {
	return nil
}

const wasmClearPragmas = false

func wasmDSNSuffix() string { return "" }
