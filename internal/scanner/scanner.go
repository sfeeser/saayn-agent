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
// 1. Change the return type to a slice of POINTERS []*model.Node
func FullScan(root string) ([]*model.Node, error) {
	var nodes []*model.Node
	fset := token.NewFileSet() // This is the "Coordinate Map" for the whole scan

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		// Parse the file using the shared fset
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		pkgName := f.Name.Name

		ast.Inspect(f, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Pass the 'path' variable from the WalkDir loop into the helper
			newNode := extractNodeMetadata(pkgName, path, fn, fset)
			nodes = append(nodes, newNode)

			return true
		})

		return nil
	})

	return nodes, err
}

// extractNodeMetadata turns a raw AST function into a CGS Node
func extractNodeMetadata(pkg string, filePath string, fn *ast.FuncDecl, fset *token.FileSet) *model.Node {
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

	// 🚨 THE FIX: Append the filename so multiple init() functions don't collide
	fileName := filepath.Base(filePath)
	uniqueID := fmt.Sprintf("%s[%s]", identity, fileName)

	return &model.Node{
		PublicID: uniqueID,
		AST:      fn,
		Fset:     fset,
	}
}
