package scanner

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

// ExtractSignatures parses a Go file and returns a clean representation containing
// only declarations (structs, interfaces, functions, methods) with their signatures.
// Function bodies are stripped to save tokens when feeding context to the LLM.
func ExtractSignatures(filePath string) (string, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Remove function bodies while keeping signatures and comments
	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			fn.Body = nil
		}
		return true
	})

	var buf bytes.Buffer
	cfg := printer.Config{
		Mode:     printer.TabIndent | printer.UseSpaces,
		Tabwidth: 4,
	}

	if err := cfg.Fprint(&buf, fset, f); err != nil {
		return "", fmt.Errorf("failed to render signatures from %s: %w", filePath, err)
	}

	return buf.String(), nil
}
