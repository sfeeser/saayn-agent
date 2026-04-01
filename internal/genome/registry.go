package genome

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go/ast"
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

	// 🧠 THE FIX: Safely determine what kind of AST node we are hashing
	switch node := n.AST.(type) {
	case *ast.FuncDecl:
		// It's a function! Hash its internal logic (the Body)
		if node.Body != nil {
			_ = conf.Fprint(&buf, n.Fset, node.Body)
		}
	case *ast.GenDecl:
		// It's a struct! Hash the structure itself
		_ = conf.Fprint(&buf, n.Fset, node)
	default:
		// Fallback
		_ = conf.Fprint(&buf, n.Fset, node)
	}

	body := buf.String()

	// COLLAPSE WHITESPACE
	re := regexp.MustCompile(`\s+`)
	body = re.ReplaceAllString(body, " ")
	body = strings.TrimSpace(body)

	// CALCULATE HASH
	hash := sha256.Sum256([]byte(body))
	hashStr := fmt.Sprintf("%x", hash)

	// 💅 THE PRETTY PRINT (Live Pulse Mode)
	if len(hashStr) >= 8 {
		icon := "⚙️ "
		if n.NodeType == "struct" {
			icon = "📦"
		}
		// Use \r to return to the start of the line and %-60s to overwrite old text
		fmt.Printf("\r\033[2K   %s Indexing %-60s", icon, n.PublicID)
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
