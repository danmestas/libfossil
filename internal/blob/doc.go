// Package blob handles content-addressed blob storage in Fossil
// repository databases.
//
// Fossil's blob format is a 4-byte big-endian uncompressed-size prefix
// followed by zlib-compressed data. [Compress] and [Decompress] handle
// this encoding transparently.
//
// [Store] compresses content, computes its SHA1 hash, and inserts it
// into the blob table. [StoreDelta] does the same for delta-encoded
// blobs and records the delta relationship. [Load] retrieves and
// decompresses a blob by RID.
package blob
