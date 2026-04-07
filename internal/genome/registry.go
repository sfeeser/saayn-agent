package genome

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/printer"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sfeeser/saayn-agent/pkg/model"
)

// RegistryManager handles the lifecycle of the genome.json
type RegistryManager struct {
	Registry *model.Registry
	FilePath string
}

// NewRegistry initializes a new registry from scanned nodes
func NewRegistry(scannedNodes []*model.Node, path string) *RegistryManager {
	reg := &model.Registry{
		ProjectName: "SAAYN-Genome",
		Nodes:       make(map[string]*model.Node), // 🧠 Pointer Map
	}

	rm := &RegistryManager{Registry: reg, FilePath: path}

	for _, n := range scannedNodes {
		if n.AST != nil {
			n.LogicHash = rm.NormalizeAndHash(n)
		}

		n.Version = 1
		n.LastModified = time.Now()
		reg.Nodes[n.UUID] = n // 🧠 Store the pointer directly
	}
	return rm
}

// NormalizeAndHash turns a function body into a stable LogicHash
func (rm *RegistryManager) NormalizeAndHash(n *model.Node) string {
	if n.AST == nil || n.Fset == nil {
		return ""
	}

	var buf strings.Builder
	conf := &printer.Config{Mode: printer.RawFormat, Tabwidth: 8}

	switch node := n.AST.(type) {
	case *ast.FuncDecl:
		cp := *node
		cp.Doc = nil
		_ = conf.Fprint(&buf, n.Fset, &cp)
	case *ast.GenDecl:
		cp := *node
		cp.Doc = nil
		_ = conf.Fprint(&buf, n.Fset, &cp)
	default:
		_ = conf.Fprint(&buf, n.Fset, node)
	}

	body := buf.String()

	// COLLAPSE WHITESPACE
	re := regexp.MustCompile(`\s+`)
	body = re.ReplaceAllString(body, " ")
	body = strings.TrimSpace(body)

	// CALCULATE HASH
	hash := sha256.Sum256([]byte(body))
	return fmt.Sprintf("%x", hash)
}

// Save persists the registry to the genome.json file
func (rm *RegistryManager) Save() error {
	return SaveRegistry(rm.Registry, rm.FilePath)
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

// SaveRegistry flushes the current project state to a JSON "genome" file.
func SaveRegistry(reg *model.Registry, path string) error {
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}
