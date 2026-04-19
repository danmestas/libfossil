package deck

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
)

func (d *Deck) ComputeR(getContent func(uuid string) ([]byte, error)) (string, error) {
	if len(d.F) == 0 {
		return "d41d8cd98f00b204e9800998ecf8427e", nil
	}
	sorted := make([]FileCard, len(d.F))
	copy(sorted, d.F)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })

	h := md5.New()
	for _, f := range sorted {
		if f.UUID == "" {
			continue
		}
		content, err := getContent(f.UUID)
		if err != nil {
			return "", fmt.Errorf("ComputeR: fetching %q: %w", f.UUID, err)
		}
		h.Write([]byte(f.Name))
		h.Write([]byte(fmt.Sprintf(" %d\n", len(content))))
		h.Write(content)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
