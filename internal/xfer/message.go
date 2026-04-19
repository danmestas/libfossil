package xfer

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

// Message is a sequence of cards forming one xfer request or response.
type Message struct {
	Cards []Card
}

// maxDecompressedBytes caps the decompressed payload size to prevent
// a small compressed payload from expanding to exhaust memory.
const maxDecompressedBytes = 50 * 1024 * 1024 // 50 MB

// Encode serializes all cards and zlib-compresses the result.
// Uses Fossil's compression format: 4-byte big-endian uncompressed size prefix
// followed by standard zlib data.
func (m *Message) Encode() ([]byte, error) {
	if m == nil {
		panic("xfer.Message.Encode: m must not be nil")
	}
	raw, err := m.EncodeUncompressed()
	if err != nil {
		return nil, err
	}
	var zbuf bytes.Buffer
	// 4-byte big-endian uncompressed size prefix (Fossil's blob_compress format).
	var sizePrefix [4]byte
	binary.BigEndian.PutUint32(sizePrefix[:], uint32(len(raw)))
	zbuf.Write(sizePrefix[:])
	zw := zlib.NewWriter(&zbuf)
	if _, err := zw.Write(raw); err != nil {
		return nil, fmt.Errorf("xfer: message zlib write: %w", err)
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("xfer: message zlib close: %w", err)
	}
	return zbuf.Bytes(), nil
}

// EncodeUncompressed serializes all cards without zlib compression.
func (m *Message) EncodeUncompressed() ([]byte, error) {
	var buf bytes.Buffer
	for i, c := range m.Cards {
		if err := EncodeCard(&buf, c); err != nil {
			return nil, fmt.Errorf("xfer: encode card %d (%T): %w", i, c, err)
		}
	}
	return buf.Bytes(), nil
}

// Decode decodes an xfer message. It tries three formats in order:
//  1. Raw zlib (Fossil HTTP sync wire format — no size prefix).
//  2. 4-byte BE size prefix + zlib (Fossil blob compression / our Encode format).
//  3. Uncompressed card data (clone protocol v3, x-fossil-uncompressed).
//
// The first format that successfully decompresses wins.
func Decode(data []byte) (*Message, error) {
	if len(data) == 0 {
		return &Message{}, nil
	}

	// Format 1: raw zlib (Fossil HTTP /xfer wire format).
	if raw, err := decompressBounded(data); err == nil {
		return DecodeUncompressed(raw)
	}

	// Format 2: 4-byte size prefix + zlib (our Encode format).
	if len(data) >= 4 {
		if raw, err := decompressBounded(data[4:]); err == nil {
			return DecodeUncompressed(raw)
		}
	}

	// Format 3: uncompressed card data.
	return DecodeUncompressed(data)
}

// decompressBounded decompresses zlib data with a size cap.
// Returns an error if decompression fails or exceeds maxDecompressedBytes.
func decompressBounded(data []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	lr := io.LimitReader(zr, maxDecompressedBytes+1)
	raw, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if len(raw) > maxDecompressedBytes {
		return nil, fmt.Errorf("xfer: decompressed payload exceeds %d bytes", maxDecompressedBytes)
	}
	return raw, nil
}

// DecodeUncompressed decodes cards from uncompressed data.
func DecodeUncompressed(data []byte) (*Message, error) {
	r := bufio.NewReader(bytes.NewReader(data))
	msg := &Message{}
	for {
		card, err := DecodeCard(r)
		if err == io.EOF {
			return msg, nil
		}
		if err != nil {
			return nil, fmt.Errorf("xfer: decode card %d: %w", len(msg.Cards), err)
		}
		msg.Cards = append(msg.Cards, card)
	}
}
