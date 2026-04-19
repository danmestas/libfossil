package deck

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func VerifyZ(data []byte) error {
	if len(data) < 35 {
		return fmt.Errorf("deck.VerifyZ: manifest too short (%d bytes)", len(data))
	}
	tail := data[len(data)-35:]
	if tail[0] != 'Z' || tail[1] != ' ' || tail[34] != '\n' {
		return fmt.Errorf("deck.VerifyZ: invalid Z-card format")
	}
	stated := string(tail[2:34])
	computed := computeZ(data[:len(data)-35])
	if computed != stated {
		return fmt.Errorf("deck.VerifyZ: checksum mismatch: stated=%s computed=%s", stated, computed)
	}
	return nil
}

func computeZ(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}
