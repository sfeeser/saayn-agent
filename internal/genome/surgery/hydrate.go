package surgery

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// HydrateNode reads a file, locates the target node via AST, and extracts its exact raw source code.
func HydrateNode(projectRoot string, anchor TargetAnchor, reason string) (ContextNode, error) {
	fullPath := filepath.Join(projectRoot, anchor.FilePath)

	// 1. Read the raw file bytes
	rawBytes, err := os.ReadFile(fullPath)
	if err != nil {
		return ContextNode{}, fmt.Errorf("failed to read file for hydration %s: %w", fullPath, err)
	}

	// 2. Parse the AST
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fullPath, rawBytes, parser.ParseComments)
	if err != nil {
		return ContextNode{}, fmt.Errorf("failed to parse AST for hydration %s: %w", fullPath, err)
	}

	var foundNode ast.Node

	// 3. Walk the AST to find the exact node
	ast.Inspect(f, func(n ast.Node) bool {
		if foundNode != nil {
			return false // Stop searching if we already found it
		}

		switch node := n.(type) {
		case *ast.FuncDecl:
			identity := extractFuncIdentity(f.Name.Name, filepath.Base(fullPath), node)
			if identity == anchor.PublicID {
				foundNode = node
				return false
			}

		case *ast.GenDecl:
			if genDeclMatches(f.Name.Name, filepath.Base(fullPath), node, anchor.PublicID) {
				foundNode = node
				return false
			}
		}

		return true
	})

	if foundNode == nil {
		return ContextNode{}, fmt.Errorf("node %s not found in AST of %s", anchor.PublicID, anchor.FilePath)
	}

	// 4. The Byte-Slicing Trick with Bounds Safety
	file := fset.File(foundNode.Pos())
	if file == nil {
		return ContextNode{}, fmt.Errorf("failed to resolve token file for %s", anchor.PublicID)
	}

	startByte := file.Offset(foundNode.Pos())
	endByte := file.Offset(foundNode.End())

	if startByte < 0 || endByte < startByte || endByte > len(rawBytes) {
		return ContextNode{}, fmt.Errorf("invalid source offsets for %s (start: %d, end: %d, max: %d)", anchor.PublicID, startByte, endByte, len(rawBytes))
	}

	// Slice the raw file to get the exact original source code
	sourceCode := string(rawBytes[startByte:endByte])

	// 5. Package it into our Schema Contract
	return ContextNode{
		TargetAnchor: anchor,
		Reason:       reason,
		SourceCode:   sourceCode,
	}, nil
}

// --- Extraction Helpers ---

func extractFuncIdentity(pkg string, fileName string, fn *ast.FuncDecl) string {
	receiver := ""
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		switch t := fn.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				// CRITICAL: Do not strip the star, it must match the scanner identity
				receiver = "*" + ident.Name
			}
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

func genDeclMatches(pkg, fileName string, decl *ast.GenDecl, publicID string) bool {
	// Iterate over ALL specs in the declaration block
	for _, spec := range decl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			id := fmt.Sprintf("%s.%s[%s]", pkg, typeSpec.Name.Name, fileName)
			if id == publicID {
				return true
			}
		}
		if valSpec, ok := spec.(*ast.ValueSpec); ok {
			// A single var spec can have multiple names (e.g., var a, b int)
			for _, name := range valSpec.Names {
				id := fmt.Sprintf("%s.%s[%s]", pkg, name.Name, fileName)
				if id == publicID {
					return true
				}
			}
		}
	}
	return false
}
