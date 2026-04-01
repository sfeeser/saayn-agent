package model

import (
	"go/ast"
	"go/token"
	"time"
)

// TestHealth tracks the biological "immunity" of a function
type TestHealth struct {
	HasTest    bool   `json:"has_test"`
	LastPassed bool   `json:"last_passed"`
	TestedAt   string `json:"tested_at"`
}

// Node represents a single "Twig" in the Code Genome System.
// It is the atomic unit of mutation and identity.
type Node struct {
	// Internal Identity (The Soul)
	UUID string `json:"uuid"`

	// Public Identity (The Address)
	PublicID string `json:"public_id"`

	// Structural Fingerprint (The Skeleton)
	Fingerprint string `json:"fingerprint"`

	// Logic Hash (The Muscle)
	LogicHash string `json:"logic_hash"`

	// Add this so the AI knows if it's mutating logic or data
	NodeType string `json:"node_type"` // e.g., "function" or "struct"
	FilePath string `json:"file_path"` //
	// Dependency Model (The Connections)
	Dependencies DependencyMap `json:"dependencies"`

	// Test Tracking (The Immune System)
	TestHealth TestHealth `json:"test_health"`

	// Metadata (The Context)
	BusinessPurpose string    `json:"business_purpose"`
	LastModified    time.Time `json:"last_modified"`
	Version         int       `json:"version"`

	// 🧬 The "Live" AST Twig (Memory only)
	AST  ast.Node       `json:"-"`
	Fset *token.FileSet `json:"-"`
}

// DependencyMap categorizes how a node interacts with the rest of the system.
type DependencyMap struct {
	Calls  []string `json:"calls"`
	Reads  []string `json:"reads"`
	Writes []string `json:"writes"`
}

// Registry represents the full genome.json file structure.
type Registry struct {
	ProjectName     string          `json:"project_name"`
	RootPath        string          `json:"root_path"`
	LastMutatedUUID string          `json:"last_mutated_uuid,omitempty"`
	Nodes           map[string]Node `json:"nodes"`
}

type SurgeryRequest struct {
	TargetUUID  string
	NewLogic    string
	Registry    *Registry
	ProjectRoot string
}
