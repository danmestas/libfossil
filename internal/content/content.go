package content

import (
	"fmt"

	libfossil "github.com/danmestas/libfossil/internal/fsltype"
	"github.com/danmestas/libfossil/internal/blob"
	"github.com/danmestas/libfossil/db"
	"github.com/danmestas/libfossil/internal/delta"
	"github.com/danmestas/libfossil/internal/hash"
	"github.com/danmestas/libfossil/simio"
)

func Expand(q db.Querier, rid libfossil.FslID) (result []byte, err error) {
	if q == nil {
		panic("content.Expand: q must not be nil")
	}
	if rid <= 0 {
		panic("content.Expand: rid must be > 0")
	}
	defer func() {
		if err == nil && result == nil {
			panic("content.Expand: postcondition violated: result is nil with no error")
		}
	}()

	chain, err := walkDeltaChain(q, rid)
	if err != nil {
		return nil, fmt.Errorf("content.Expand: %w", err)
	}

	content, err := blob.Load(q, chain[0])
	if err != nil {
		return nil, fmt.Errorf("content.Expand load root rid=%d: %w", chain[0], err)
	}

	for i := 1; i < len(chain); i++ {
		deltaBytes, err := blob.Load(q, chain[i])
		if err != nil {
			return nil, fmt.Errorf("content.Expand load delta rid=%d: %w", chain[i], err)
		}
		content, err = delta.Apply(content, deltaBytes)
		if err != nil {
			return nil, fmt.Errorf("content.Expand apply delta rid=%d: %w", chain[i], err)
		}
	}

	// BUGGIFY: flip a byte in expanded content to exercise UUID-mismatch detection.
	if simio.Buggify(0.01) && len(content) > 0 {
		corrupted := make([]byte, len(content))
		copy(corrupted, content)
		corrupted[0] ^= 0xFF
		return corrupted, nil
	}

	return content, nil
}

func walkDeltaChain(q db.Querier, rid libfossil.FslID) (chain []libfossil.FslID, err error) {
	if q == nil {
		panic("content.walkDeltaChain: q must not be nil")
	}
	if rid <= 0 {
		panic("content.walkDeltaChain: rid must be > 0")
	}
	defer func() {
		if err == nil && len(chain) == 0 {
			panic("content.walkDeltaChain: postcondition violated: chain is empty with no error")
		}
	}()

	current := rid
	seen := make(map[libfossil.FslID]bool)

	for {
		if seen[current] {
			return nil, fmt.Errorf("delta chain cycle detected at rid=%d", current)
		}
		seen[current] = true
		chain = append(chain, current)

		var sourceID int64
		err := q.QueryRow("SELECT srcid FROM delta WHERE rid=?", current).Scan(&sourceID)
		if err != nil {
			break
		}
		current = libfossil.FslID(sourceID)
	}

	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain, nil
}

func Verify(q db.Querier, rid libfossil.FslID) error {
	if q == nil {
		panic("content.Verify: q must not be nil")
	}
	if rid <= 0 {
		panic("content.Verify: rid must be > 0")
	}

	var uuid string
	err := q.QueryRow("SELECT uuid FROM blob WHERE rid=?", rid).Scan(&uuid)
	if err != nil {
		return fmt.Errorf("content.Verify query uuid: %w", err)
	}

	content, err := Expand(q, rid)
	if err != nil {
		return fmt.Errorf("content.Verify expand: %w", err)
	}

	var computed string
	if len(uuid) == 64 {
		computed = hash.SHA3(content)
	} else {
		computed = hash.SHA1(content)
	}

	if computed != uuid {
		return fmt.Errorf("content.Verify: hash mismatch for rid=%d: stored=%s computed=%s", rid, uuid, computed)
	}
	return nil
}

func IsPhantom(q db.Querier, rid libfossil.FslID) (bool, error) {
	if q == nil {
		panic("content.IsPhantom: q must not be nil")
	}
	if rid <= 0 {
		panic("content.IsPhantom: rid must be > 0")
	}

	var count int
	err := q.QueryRow("SELECT count(*) FROM phantom WHERE rid=?", rid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
