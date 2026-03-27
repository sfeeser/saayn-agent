package genome

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/saayn-agent/pkg/model"
)

// RegistryManager handles the lifecycle of the genome.json
type RegistryManager struct {
	Registry *model.Registry
	FilePath string
}

// NewRegistry initializes a new registry from scanned nodes
func NewRegistry(scannedNodes []model.Node, path string) *RegistryManager {
	reg := &model.Registry{
		ProjectName: "CGS-Project",
		Nodes:       make(map[string]model.Node),
	}

	for _, n := range scannedNodes {
		// Generate stable internal identity
		id := uuid.New().String()
		n.UUID = id
		n.Version = 1
		n.LastModified = time.Now()
		
		reg.Nodes[id] = n
	}

	return &RegistryManager{
		Registry: reg,
		FilePath: path,
	}
}

// NormalizeAndHash turns a function body into a stable LogicHash
func (rm *RegistryManager) NormalizeAndHash(fn *ast.FuncDecl) string {
	fset := token.NewFileSet()
	var buf strings.Builder

	// 1. Strip comments and render AST to string
	// We use a config that ignores comments during printing
	conf := &printer.Config{Mode: printer.RawFormat, Tabwidth: 8}
	conf.Fprint(&buf, fset, fn.Body)

	body := buf.String()

	// 2. Normalize Whitespace (Collapse all to single spaces)
	reWhitespace := regexp.MustCompile(`\s+`)
	body = reWhitespace.ReplaceAllString(body, " ")

	// 3. Anonymize Local Identifiers (Phase 1: Simple cleanup)
	// Future versions will implement the v1, v2 mapping here
	body = strings.TrimSpace(body)

	// 4. Generate SHA-256
	hash := sha256.Sum256([]byte(body))
	return fmt.Sprintf("%x", hash)
}

// Save persists the registry to the genome.json file
func (rm *RegistryManager) Save() error {
	data, err := json.MarshalIndent(rm.Registry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(rm.FilePath, data, 0644)
}

// Load reads an existing registry from disk
func Load(path string) (*RegistryManager, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var reg model.Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &RegistryManager{Registry: &reg, FilePath: path}, nil
}
