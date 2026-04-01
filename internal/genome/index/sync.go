package index

import (
	"context"
	"fmt"

	"github.com/saayn-agent/pkg/model"
)

// SyncStats provides telemetry for the indexing process.
type SyncStats struct {
	Created int
	Updated int
	Deleted int
	Skipped int
}

// SyncIndex synchronizes the IndexStore with the current state of the Genome Registry.
// It handles metadata validation, drift detection, embedding retrieval, and stale record cleanup.
func SyncIndex(
	ctx context.Context,
	store *IndexStore,
	nodes map[string]model.Node,
	apiKey, modelName, baseURL string,
) (SyncStats, error) {
	stats := SyncStats{}

	if store == nil {
		return stats, fmt.Errorf("index store is nil")
	}

	// 1. Metadata Guard: If the model changed, the entire vector space is invalid.
	if store.Metadata.EmbeddingModel != "" && store.Metadata.EmbeddingModel != modelName {
		return stats, fmt.Errorf(
			"index model mismatch: store uses %s, but config specifies %s. Manual rebuild required",
			store.Metadata.EmbeddingModel,
			modelName,
		)
	}

	// Ensure metadata is set for a fresh store.
	if store.Metadata.EmbeddingModel == "" {
		store.Metadata.EmbeddingModel = modelName
	}

	// 2. Process Current Genome Nodes (Create/Update/Delete stale semantic entries)
	processedUUIDs := make(map[string]struct{})

	for uuid, node := range nodes {
		record, exists := store.Get(uuid)

		// If the node is not enriched, it should not remain searchable.
		// This handles the case where a node's logic drifted and its purpose was reset.
		if node.BusinessPurpose == "" {
			if exists {
				delete(store.Records, uuid)
				stats.Deleted++
			} else {
				stats.Skipped++
			}
			processedUUIDs[uuid] = struct{}{}
			continue
		}

		input := RetrievalInput{
			PublicID: node.PublicID,
			NodeType: node.NodeType,
			FilePath: node.FilePath,
			Purpose:  node.BusinessPurpose,
		}

		newHash := HashRetrievalDoc(input)

		// Drift detection: skip if the record exists and nothing has changed.
		if exists && record.RetrievalDocHash == newHash {
			stats.Skipped++
			processedUUIDs[uuid] = struct{}{}
			continue
		}

		// Fetch New Embedding
		doc := BuildRetrievalDoc(input)
		// Fix: Signature now matches our finalized FetchEmbedding (no extra nil client)
		vec, err := FetchEmbedding(ctx, nil, apiKey, modelName, baseURL, doc)
		if err != nil {
			return stats, fmt.Errorf("sync failed at node %s: %w", node.PublicID, err)
		}

		// Dimension guard
		if store.Metadata.VectorDimension == 0 {
			store.Metadata.VectorDimension = len(vec)
		} else if store.Metadata.VectorDimension != len(vec) {
			return stats, fmt.Errorf(
				"vector dimension mismatch: expected %d, got %d from model %s",
				store.Metadata.VectorDimension,
				len(vec),
				modelName,
			)
		}

		// Normalize before storing so search can use dot product as cosine similarity.
		normalizedVec := Normalize(vec)

		store.Upsert(EmbeddingRecord{
			UUID:              uuid,
			PublicID:          node.PublicID,
			FilePath:          node.FilePath,
			NodeType:          node.NodeType,
			RetrievalDocument: doc,
			RetrievalDocHash:  newHash,
			Vector:            normalizedVec,
		})

		if exists {
			stats.Updated++
		} else {
			stats.Created++
		}

		processedUUIDs[uuid] = struct{}{}
	}

	// 3. Stale Record Cleanup (Delete orphaned index entries no longer in genome)
	for uuid := range store.Records {
		if _, ok := processedUUIDs[uuid]; !ok {
			delete(store.Records, uuid)
			stats.Deleted++
		}
	}

	return stats, nil
}
