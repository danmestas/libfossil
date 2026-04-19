package uv

// Status compares a local unversioned file against a remote one and returns
// an action code. Exact port of Fossil's unversioned_status().
//
// localMtime==0 means no local row (returns 0).
// localHash="" with localMtime>0 means tombstone (participates in mtime comparison).
// "-" means deletion marker in either position.
//
// Return codes:
//
//	0 = not present locally (pull)
//	1 = different hash, remote newer or tiebreaker (pull)
//	2 = same hash, remote mtime older (pull mtime only)
//	3 = identical (no action)
//	4 = same hash, remote mtime newer (push mtime only)
//	5 = different hash, local newer or tiebreaker (push)
func Status(localMtime int64, localHash string, remoteMtime int64, remoteHash string) int {
	if localMtime == 0 {
		return 0
	}

	var mtimeCmp int
	switch {
	case localMtime < remoteMtime:
		mtimeCmp = -1
	case localMtime > remoteMtime:
		mtimeCmp = 1
	}

	hashCmp := cmpStr(localHash, remoteHash)

	if hashCmp == 0 {
		return 3 + mtimeCmp
	}

	// Tiebreaker when mtimes are equal and hashes differ:
	// Use hash comparison, but deletion marker "-" always wins
	if mtimeCmp == 0 {
		if localHash == "-" {
			return 5 // Push local deletion
		}
		if remoteHash == "-" {
			return 1 // Pull remote deletion
		}
	}

	if mtimeCmp < 0 || (mtimeCmp == 0 && hashCmp < 0) {
		return 1
	}
	return 5
}

func cmpStr(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
