package index

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// IndexMetadata tracks embedding compatibility for the index.
type IndexMetadata struct {
	EmbeddingModel  string `json:"embedding_model"`
	VectorDimension int    `json:"vector_dimension"`
}

// EmbeddingRecord represents a single searchable unit in the index.
type EmbeddingRecord struct {
	UUID              string    `json:"uuid"`
	PublicID          string    `json:"public_id"`
	FilePath          string    `json:"file_path"`
	NodeType          string    `json:"node_type"`
	RetrievalDocument string    `json:"retrieval_document"`
	RetrievalDocHash  string    `json:"retrieval_doc_hash"`
	Vector            []float32 `json:"vector"`
}

// IndexStore is the persistent container for all records and their associated metadata.
type IndexStore struct {
	Metadata IndexMetadata              `json:"metadata"`
	Records  map[string]EmbeddingRecord `json:"records"`
}

// RetrievalInput aggregates the fields required to build a canonical document.
type RetrievalInput struct {
	PublicID string
	NodeType string
	FilePath string
	Purpose  string
}

// NewIndexStore initializes a new store with specific model metadata.
func NewIndexStore(modelName string, dimension int) *IndexStore {
	return &IndexStore{
		Metadata: IndexMetadata{
			EmbeddingModel:  modelName,
			VectorDimension: dimension,
		},
		Records: make(map[string]EmbeddingRecord),
	}
}

// Save persists the IndexStore to disk as a JSON file.
func (s *IndexStore) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create index directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal index: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write index file: %w", err)
	}

	return nil
}

// LoadIndex reads an IndexStore from a JSON file.
func LoadIndex(path string) (*IndexStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read index file: %w", err)
	}

	var store IndexStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("unmarshal index: %w", err)
	}

	// Defensive: Ensure the map is initialized even if JSON had it as null/missing
	if store.Records == nil {
		store.Records = make(map[string]EmbeddingRecord)
	}

	return &store, nil
}

// BuildRetrievalDoc builds the canonical retrieval text used for embedding and hashing.
func BuildRetrievalDoc(in RetrievalInput) string {
	return fmt.Sprintf(
		"PublicID: %s\nNodeType: %s\nFilePath: %s\nBusinessPurpose: %s",
		in.PublicID,
		in.NodeType,
		in.FilePath,
		in.Purpose,
	)
}

// HashRetrievalDoc returns the SHA-256 hash of the canonical retrieval document.
func HashRetrievalDoc(in RetrievalInput) string {
	content := BuildRetrievalDoc(in)
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}

// Upsert adds or updates a record in the index.
func (s *IndexStore) Upsert(record EmbeddingRecord) {
	if s.Records == nil {
		s.Records = make(map[string]EmbeddingRecord)
	}
	s.Records[record.UUID] = record
}

// Get retrieves a record by its UUID.
func (s *IndexStore) Get(uuid string) (EmbeddingRecord, bool) {
	if s.Records == nil {
		return EmbeddingRecord{}, false
	}
	rec, ok := s.Records[uuid]
	return rec, ok
}
