// Package content reconstructs full artifact content from Fossil's
// delta-chain storage.
//
// Fossil stores blobs either as full content or as deltas against a
// source blob. [Expand] walks the delta chain from a given RID back to
// the root, decompresses each link, and applies deltas in sequence to
// produce the original content. Cycle detection prevents infinite loops.
//
// [Verify] expands a blob and compares its hash (SHA1 or SHA3-256)
// against the stored UUID. [IsPhantom] checks whether a blob is a
// placeholder awaiting delivery during sync.
package content
