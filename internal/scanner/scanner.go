package scanner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/saayn-agent/pkg/model"
)

// FullScan walks the directory and extracts all functional nodes
func FullScan(root string) ([]model.Node, error) {
	var nodes []model.Node
	fset := token.NewFileSet()

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		// 1. Parse the Go file into an AST
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files that don't parse (broken code)
		}

		// 2. Extract the Package Name
		pkgName := f.Name.Name

		// 3. Inspect the AST for Declarations
		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true // Keep looking for functions
			}

			// 4. Build the Semantic Identity (The DNA Address)
			node := extractNodeMetadata(pkgName, fn)
			nodes = append(nodes, node)
			
			return true
		})

		return nil
	})

	return nodes, err
}

// extractNodeMetadata turns a raw AST function into a CGS Node
func extractNodeMetadata(pkg string, fn *ast.FuncDecl) model.Node {
	receiver := ""
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		// It's a method. Extract the receiver type (e.g., "*Service")
		switch t := fn.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			receiver = fmt.Sprintf("*%v", t.X)
		case *ast.Ident:
			receiver = t.Name
		}
	}

	// Format: package.Receiver.Name
	identity := pkg
	if receiver != "" {
		identity += "." + receiver
	}
	identity += "." + fn.Name.Name

	return model.Node{
		PublicID: identity,
		// Note: LogicHash and Fingerprint will be generated in the Genome package
	}
}
