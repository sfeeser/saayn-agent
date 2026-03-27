package surgeon

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/saayn-agent/internal/validator"
	"github.com/saayn-agent/pkg/model"
)

// PerformSurgery executes the "Locate -> Swap -> Validate -> Commit/Rollback" loop
func PerformSurgery(req model.SurgeryRequest) error {
	// 1. Locate the Node in the Registry
	node, ok := req.Registry.Nodes[req.TargetUUID]
	if !ok {
		return fmt.Errorf("UUID %s not found in genome", req.TargetUUID)
	}

	// 2. Find the physical file
	filePath, err := findFileForNode(req.ProjectRoot, node.PublicID)
	if err != nil {
		return err
	}

	// 3. Parse the file into an AST
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// 4. THE INCISION: Parse the New Logic safely
	// We wrap the raw code in a dummy function so the Go parser understands it natively
	dummyCode := fmt.Sprintf("package p\nfunc f() {\n%s\n}", req.NewLogic)
	dummyF, err := parser.ParseFile(token.NewFileSet(), "", dummyCode, 0)
	if err != nil {
		return fmt.Errorf("syntax error in new logic block: %v", err)
	}
	newBody := dummyF.Decls[0].(*ast.FuncDecl).Body

	// 5. THE GRAFT: Swap the AST nodes
	mutated := false
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		identity := extractIdentity(f.Name.Name, filepath.Base(filePath), fn)
		if identity == node.PublicID {
			fn.Body = newBody // Inject the new DNA
			mutated = true
		}
		return true
	})

	if !mutated {
		return fmt.Errorf("failed to locate node %s in AST of %s", node.PublicID, filePath)
	}

	// 6. Format the mutated AST back to a source code string
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("formatting failed: %v", err)
	}

	// 7. THE GATE: Atomic Write and Rollback
	originalCode, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Overwrite the live file with the new logic
	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return err
	}

	// Run the compiler/validator on the LIVE project
	res, err := validator.CheckIntegrity(req.ProjectRoot)
	if err != nil || !res.Success {
		// 🚨 ROLLBACK: The code was invalid, revert to the original backup
		os.WriteFile(filePath, originalCode, 0644)
		return fmt.Errorf("surgery rejected by compiler. Rolled back.\nErrors:\n%s", res.Output)
	}

	// If we make it here, the surgery was a success!
	return nil
}

// extractIdentity matches the scanner's exact format: pkg.Receiver.Name[file.go]
func extractIdentity(pkg, fileName string, fn *ast.FuncDecl) string {
	receiver := ""
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		switch t := fn.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			receiver = fmt.Sprintf("*%v", t.X)
		case *ast.Ident:
			receiver = t.Name
		}
	}

	identity := pkg
	if receiver != "" {
		identity += "." + receiver
	}
	identity += "." + fn.Name.Name

	return fmt.Sprintf("%s[%s]", identity, fileName)
}

// findFileForNode walks the directory to locate the exact file containing the target function
func findFileForNode(root, targetID string) (string, error) {
	var foundPath string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		base := filepath.Base(path)
		// Optimization: Only parse files if their filename matches the bracket in the ID
		if strings.Contains(targetID, "["+base+"]") {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return nil
			}

			ast.Inspect(f, func(n ast.Node) bool {
				if fn, ok := n.(*ast.FuncDecl); ok {
					if extractIdentity(f.Name.Name, base, fn) == targetID {
						foundPath = path
					}
				}
				return true
			})
		}
		return nil
	})

	if foundPath != "" {
		return foundPath, nil
	}
	return "", err
}
