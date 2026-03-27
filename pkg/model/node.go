package model

import "time"

// Node represents a single "Twig" in the Code Genome System.
// It is the atomic unit of mutation and identity.
type Node struct {
	// Internal Identity (The Soul)
	// Immutable UUID that survives refactors and renames.
	UUID string `json:"uuid"`

	// Public Identity (The Address)
	// Format: package.Receiver.Name
	// Example: billing.Service.CalculateFees
	PublicID string `json:"public_id"`

	// Structural Fingerprint (The Skeleton)
	// Hash of the canonical signature (params and returns).
	// Detects if the API contract changed.
	Fingerprint string `json:"fingerprint"`

	// Logic Hash (The Muscle)
	// Hash of the normalized, anonymized function body.
	// Detects if the internal behavior changed.
	LogicHash string `json:"logic_hash"`

	// Dependency Model (The Connections)
	// Tracks what this node touches and who touches it.
	Dependencies DependencyMap `json:"dependencies"`

	// Metadata (The Context)
	BusinessPurpose string    `json:"business_purpose"`
	LastModified    time.Time `json:"last_modified"`
	Version         int       `json:"version"`
}

// DependencyMap categorizes how a node interacts with the rest of the system.
type DependencyMap struct {
	Calls  []string `json:"calls"`  // Functions this node invokes
	Reads  []string `json:"reads"`  // Global state or config this node accesses
	Writes []string `json:"writes"` // External state this node modifies
}

// Registry represents the full genome.json file structure.
type Registry struct {
	ProjectName string          `json:"project_name"`
	RootPath    string          `json:"root_path"`
	Nodes       map[string]Node `json:"nodes"` // Keyed by UUID
}
