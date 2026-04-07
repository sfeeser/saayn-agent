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
	Status     string `json:"status"` // Required by your MaterializeNode logic
}

// DependencyMap categorizes how a node interacts with the rest of the system.
type DependencyMap struct {
	Calls  []string `json:"calls"`
	Reads  []string `json:"reads"`
	Writes []string `json:"writes"`
}

// Node represents a single unit of materialized logic (The Soul).
type Node struct {
	UUID            string        `json:"uuid"`
	PublicID        string        `json:"public_id"`
	Fingerprint     string        `json:"fingerprint"`
	LogicHash       string        `json:"logic_hash"`
	NodeType        string        `json:"node_type"`
	FilePath        string        `json:"file_path"`
	Dependencies    DependencyMap `json:"dependencies"`
	TestHealth      TestHealth    `json:"test_health"`
	BusinessPurpose string        `json:"business_purpose"`
	LastModified    time.Time     `json:"last_modified"`
	Version         int           `json:"version"`

	// 🧬 The "Live" AST Twig (Memory only)
	AST  ast.Node       `json:"-"`
	Fset *token.FileSet `json:"-"`
}
