package scanner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sfeeser/saayn-agent/pkg/model"
)

// FullScan walks the directory and extracts all functional and structural nodes
func FullScan(root string) ([]*model.Node, error) {
	var nodes []*model.Node
	fset := token.NewFileSet()

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}

		pkgName := f.Name.Name

		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {

			// Target 1: Functions and Methods
			case *ast.FuncDecl:
				nodes = append(nodes, extractFuncMetadata(pkgName, path, node, fset))

			// Target 2: General Declarations (Structs AND Variables)
			case *ast.GenDecl:
				// Handle Structs
				if node.Tok == token.TYPE {
					for _, spec := range node.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
								nodes = append(nodes, extractStructMetadata(pkgName, path, node, typeSpec, fset))
							}
						}
					}
				}
				// --- NEW: Handle Variables (like enrichCmd) ---
				if node.Tok == token.VAR {
					for _, spec := range node.Specs {
						if valSpec, ok := spec.(*ast.ValueSpec); ok {
							nodes = append(nodes, extractVarMetadata(pkgName, path, node, valSpec, fset))
						}
					}
				}
			}
			return true
		})

		return nil
	})

	return nodes, err
}

// extractFuncMetadata turns a raw AST function into a CGS Node
func extractFuncMetadata(pkg string, filePath string, fn *ast.FuncDecl, fset *token.FileSet) *model.Node {
	var receiver string
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		switch t := fn.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			receiver = fmt.Sprintf("*%v", t.X)
		case *ast.Ident:
			receiver = t.Name
		}
	}

	// IDENTITY CONSTRUCTION
	// Goal: pkg.Func OR pkg.Receiver.Func
	identity := pkg
	if receiver != "" {
		identity += "." + receiver
	}
	identity += "." + fn.Name.Name

	fileName := filepath.Base(filePath)
	uniqueID := fmt.Sprintf("%s[%s]", identity, fileName)

	return &model.Node{
		UUID:     uuid.New().String(),
		PublicID: uniqueID,
		NodeType: "function",
		FilePath: filePath,
		AST:      fn,   // Pass the AST though
		Fset:     fset, // Pass the Fset through
	}
}

// extractStructMetadata turns a raw AST struct into a CGS Node
func extractStructMetadata(pkg string, filePath string, decl *ast.GenDecl, typeSpec *ast.TypeSpec, fset *token.FileSet) *model.Node {
	identity := fmt.Sprintf("%s.%s", pkg, typeSpec.Name.Name)
	fileName := filepath.Base(filePath)
	uniqueID := fmt.Sprintf("%s[%s]", identity, fileName)

	return &model.Node{
		UUID:     uuid.New().String(), // Native UUID generation!
		PublicID: uniqueID,
		NodeType: "struct",
		FilePath: filePath,
		// We store the GenDecl so we capture the doc comments above the struct too!
		AST:  decl,
		Fset: fset,
	}
}

// extractVarMetadata turns a global variable (GenDecl + ValueSpec) into a CGS Node
func extractVarMetadata(pkg string, filePath string, decl *ast.GenDecl, vSpec *ast.ValueSpec, fset *token.FileSet) *model.Node {
	// A VarSpec can have multiple names (var a, b int), but usually Cobra cmds have one
	name := "anonymous_var"
	if len(vSpec.Names) > 0 {
		name = vSpec.Names[0].Name
	}

	identity := fmt.Sprintf("%s.%s", pkg, name)
	fileName := filepath.Base(filePath)
	uniqueID := fmt.Sprintf("%s[%s]", identity, fileName)

	return &model.Node{
		UUID:     uuid.New().String(),
		PublicID: uniqueID,
		NodeType: "struct", // Treating as a data structure/package
		FilePath: filePath,
		AST:      decl, // We store the whole GenDecl to capture all values and comments
		Fset:     fset,
	}
}
