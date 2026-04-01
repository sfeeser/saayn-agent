package index

import "sort"

// Match represents a ranked search result.
type Match struct {
	UUID     string  `json:"uuid"`
	PublicID string  `json:"public_id"`
	FilePath string  `json:"file_path"`
	NodeType string  `json:"node_type"`
	Score    float32 `json:"score"`
}

// SearchIntent ranks all records in the store against a normalized query vector.
// It expects the query vector to already be normalized to unit length.
// It returns up to k matches, ordered by similarity score.
func (s *IndexStore) SearchIntent(queryVec []float32, k int) []Match {
	// Guard against invalid k or empty query vectors
	if k <= 0 || len(queryVec) == 0 {
		return nil
	}

	// Preallocate to avoid re-allocations during the loop
	matches := make([]Match, 0, len(s.Records))

	for _, rec := range s.Records {
		// Calculate similarity score (Dot product works because vectors are normalized)
		score := Dot(queryVec, rec.Vector)

		matches = append(matches, Match{
			UUID:     rec.UUID,
			PublicID: rec.PublicID,
			FilePath: rec.FilePath,
			NodeType: rec.NodeType,
			Score:    score,
		})
	}

	// Sort results with deterministic tie-breaking
	sort.Slice(matches, func(i, j int) bool {
		// Primary sort: Descending score (highest similarity first)
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		// Secondary sort: Ascending PublicID (alphabetical)
		if matches[i].PublicID != matches[j].PublicID {
			return matches[i].PublicID < matches[j].PublicID
		}
		// Tertiary sort: Ascending UUID (guaranteed unique tie-breaker)
		return matches[i].UUID < matches[j].UUID
	})

	// Truncate to Top-K
	if len(matches) > k {
		return matches[:k]
	}

	return matches
}
