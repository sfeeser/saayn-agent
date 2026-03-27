package genome

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	_ "go/ast"
	"go/printer"
	_ "go/token"
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
func NewRegistry(scannedNodes []*model.Node, path string) *RegistryManager {
	reg := &model.Registry{
		ProjectName: "CGS-Project",
		Nodes:       make(map[string]model.Node),
	}

	rm := &RegistryManager{Registry: reg, FilePath: path}

	// n is already a *model.Node (a pointer), so no need for &scannedNodes[i]

	for _, n := range scannedNodes {
		id := uuid.New().String()
		n.UUID = id

		// Pass the WHOLE node 'n', not just n.AST
		if n.AST != nil {
			n.LogicHash = rm.NormalizeAndHash(n)
		}

		n.Version = 1
		n.LastModified = time.Now()
		reg.Nodes[id] = *n
	}

	return rm
}

// NormalizeAndHash turns a function body into a stable LogicHash
// NormalizeAndHash turns a function body into a stable LogicHash
// NormalizeAndHash turns a function body into a stable LogicHash
func (rm *RegistryManager) NormalizeAndHash(n *model.Node) string {
	if n.AST == nil || n.Fset == nil {
		return ""
	}

	var buf strings.Builder
	conf := &printer.Config{Mode: printer.RawFormat, Tabwidth: 8}
	if err := conf.Fprint(&buf, n.Fset, n.AST.Body); err != nil {
		return ""
	}

	body := buf.String()

	// COLLAPSE WHITESPACE
	re := regexp.MustCompile(`\s+`)
	body = re.ReplaceAllString(body, " ")
	body = strings.TrimSpace(body)

	// CALCULATE HASH
	hash := sha256.Sum256([]byte(body))

	// Create the hash string HERE so the printer and the return can both use it
	hashStr := fmt.Sprintf("%x", hash)

	// 💅 THE PRETTY PRINT
	// If the hash is empty (shouldn't happen here, but just in case), protect the slice
	if len(hashStr) >= 8 {
		fmt.Printf("  🧬 Indexed %-45s [%s]\n", n.PublicID, hashStr[:8])
	}

	return hashStr
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
