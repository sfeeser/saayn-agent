package surgeon

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"github.com/saayn-agent/internal/validator"
	"github.com/saayn-agent/pkg/model"
)

// GraftRequest defines the parameters for a mutation
type GraftRequest struct {
	TargetUUID string
	NewBody    string // The raw Go code returned by the AI
	Registry   *model.Registry
}

// PerformSurgery executes the "Locate -> Swap -> Validate -> Commit" loop
func PerformSurgery(projectRoot string, req GraftRequest) error {
	// 1. Locate the Node in the Registry
	node, ok := req.Registry.Nodes[req.TargetUUID]
	if !ok {
		return fmt.Errorf("UUID %s not found in genome", req.TargetUUID)
	}

	// 2. Map PublicID to a Physical File (The Scanner will help here)
	// For MVP, we assume the Scanner has cached the filepath or we re-scan
	filePath := resolvePathFromIdentity(projectRoot, node.PublicID)

	// 3. Parse the existing file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// 4. THE GRAFT: Traverse AST and swap the function body
	mutated := false
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if this is our target based on Semantic Identity
		// (In a hardened version, we'd use the Fingerprint here)
		if getIdentity(f.Name.Name, fn) == node.PublicID {
			// Parse the AI's NewBody into a block statement
			newBlock, err := parser.ParseExpr("{ " + req.NewBody + " }")
			if err == nil {
				fn.Body = newBlock.(*ast.BlockStmt)
				mutated = true
			}
		}
		return true
	})

	if !mutated {
		return fmt.Errorf("failed to locate node %s in file %s", node.PublicID, filePath)
	}

	// 5. THE GATE: Render to a buffer and validate
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return err
	}

	// Temporary write for validation
	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, buf.Bytes(), 0644); err != nil {
		return err
	}
	defer os.Remove(tmpPath) // Always clean up

	res, err := validator.CheckIntegrity(projectRoot)
	if err != nil || !res.Success {
		return fmt.Errorf("validation failed: %s", res.Output)
	}

	// 6. COMMIT: Atomic rename
	return os.Rename(tmpPath, filePath)
}

// Helper to reconstruct identity during inspection
func getIdentity(pkg string, fn *ast.FuncDecl) string {
	// (Reuse logic from scanner.go to ensure perfect matching)
	return pkg + "." + fn.Name.Name 
}

func resolvePathFromIdentity(root, id string) string {
	// Placeholder: In production, the Registry stores the last known path
	// or the Scanner finds it via a quick global walk.
	return "internal/logic/payments.go" 
}
