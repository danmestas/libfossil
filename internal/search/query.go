package search

import (
	"fmt"
	"strings"
)

const defaultMaxResults = 50

// Search executes a full-text search and returns results ranked by relevance.
// Returns empty results (not error) if term is shorter than 3 characters.
//
// Panics if idx is nil (TigerStyle precondition).
func (idx *Index) Search(q Query) ([]Result, error) {
	if idx == nil {
		panic("search.Search: nil *Index")
	}
	if len(q.Term) < 3 {
		return nil, nil
	}

	maxResults := q.MaxResults
	if maxResults <= 0 {
		maxResults = defaultMaxResults
	}

	escaped := escapeFTS5(q.Term)

	rows, err := idx.repo.DB().Query(
		"SELECT path, content FROM fts_content WHERE fts_content MATCH ? ORDER BY rank LIMIT ?",
		escaped, maxResults,
	)
	if err != nil {
		return nil, fmt.Errorf("search.Search: query: %w", err)
	}
	defer rows.Close()

	var results []Result
	for rows.Next() {
		var path, fileContent string
		if err := rows.Scan(&path, &fileContent); err != nil {
			return nil, fmt.Errorf("search.Search: scan: %w", err)
		}

		hits := findMatches(path, fileContent, q.Term, q.ContextLines)
		results = append(results, hits...)

		if len(results) >= maxResults {
			results = results[:maxResults]
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search.Search: rows: %w", err)
	}

	return results, nil
}

// maxHitsPerFile caps the number of matches returned from a single file.
// Prevents unbounded growth when a term appears on every line.
const maxHitsPerFile = 1000

// findMatches locates all occurrences of term within fileContent and returns Results.
// Case-insensitive because FTS5 trigram matching is case-insensitive — post-filtering
// must match the same semantics to avoid false negatives.
func findMatches(path, fileContent, term string, contextLines int) []Result {
	if path == "" {
		panic("findMatches: empty path")
	}
	if term == "" {
		panic("findMatches: empty term")
	}
	if contextLines < 0 {
		panic("findMatches: negative contextLines")
	}

	lines := strings.Split(fileContent, "\n")
	lowerTerm := strings.ToLower(term)
	var results []Result

	for i, line := range lines {
		lowerLine := strings.ToLower(line)
		col := strings.Index(lowerLine, lowerTerm)
		if col < 0 {
			continue
		}

		r := Result{
			Path:     path,
			Line:     i + 1,
			Column:   col,
			MatchLen: len(term),
			LineText: line,
		}

		if contextLines > 0 {
			start := i - contextLines
			if start < 0 {
				start = 0
			}
			end := i + contextLines + 1
			if end > len(lines) {
				end = len(lines)
			}
			r.Context = strings.Join(lines[start:end], "\n")
		}

		results = append(results, r)
		if len(results) >= maxHitsPerFile {
			break
		}
	}

	return results
}

// escapeFTS5 wraps term in double quotes for literal matching,
// escaping internal double quotes per FTS5 syntax.
func escapeFTS5(term string) string {
	escaped := strings.ReplaceAll(term, `"`, `""`)
	return `"` + escaped + `"`
}
