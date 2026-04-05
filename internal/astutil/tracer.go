package astutil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// CallerInfo represents a single instance where a target node was executed.
type CallerInfo struct {
	FilePath        string
	CallingFunction string
	LineNumber      int
}

// TraceCallers walks a directory tree and returns all locations where the targetNodeID is called.
// It honors package boundaries by parsing the canonical SAAYN PublicID.
func TraceCallers(rootDir string, targetNodeID string) ([]CallerInfo, error) {
	// 1. Validate and Unpack the SAAYN Node Definition
	targetPkg, targetFunc, err := parseNodeID(targetNodeID)
	if err != nil {
		return nil, fmt.Errorf("tracer abort: %w", err)
	}

	var callers []CallerInfo
	fset := token.NewFileSet()

	err = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, walkErr error) error {
		// Stop silently swallowing filesystem errors
		if walkErr != nil {
			return fmt.Errorf("filesystem access error at %s: %w", path, walkErr)
		}

		// Prune directories efficiently using OS-agnostic logic
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" || name == "node_modules" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file (Fail fast: if a file in the project is broken, impact analysis is compromised)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("syntax error parsing %s: %w", path, err)
		}

		currentPkg := file.Name.Name

		// 2. Walk the AST of the current file
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if isCallTo(call, currentPkg, targetPkg, targetFunc) {
					callerFunc := findEnclosingFunction(file, call.Pos())
					lineNum := fset.Position(call.Pos()).Line

					callers = append(callers, CallerInfo{
						FilePath:        path,
						CallingFunction: callerFunc,
						LineNumber:      lineNum,
					})
				}
			}
			return true
		})
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("trace failed: %w", err)
	}

	return callers, nil
}

// parseNodeID breaks down a SAAYN PublicID like "finance.*Calc.Compute[math.go]"
// safely, returning an error if the format is invalid.
func parseNodeID(publicID string) (pkg string, fn string, err error) {
	base := strings.SplitN(publicID, "[", 2)[0]
	parts := strings.Split(base, ".")

	if len(parts) < 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[len(parts)-1]) == "" {
		return "", "", fmt.Errorf("invalid publicID format: %q", publicID)
	}

	pkg = parts[0]
	fn = parts[len(parts)-1]
	return pkg, fn, nil
}

// isCallTo determines if an ast.CallExpr matches our target, honoring package namespaces.
func isCallTo(call *ast.CallExpr, currentPkg, targetPkg, targetFunc string) bool {
	switch fun := call.Fun.(type) {

	case *ast.Ident:
		// A local call (e.g., ComputeTax())
		return fun.Name == targetFunc && currentPkg == targetPkg

	case *ast.SelectorExpr:
		// A selector call (e.g., finance.ComputeTax() or obj.ComputeTax()).
		// Note (V4 Limitation): This uses syntactic matching only. Receiver types and
		// package aliases are not fully disambiguated without a go/types type-checker.
		if pkgIdent, ok := fun.X.(*ast.Ident); ok {
			return pkgIdent.Name == targetPkg && fun.Sel.Name == targetFunc
		}
	}
	return false
}

// findEnclosingFunction walks up the AST to find which function contains a specific byte position.
func findEnclosingFunction(file *ast.File, pos token.Pos) string {
	var enclosing string
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if pos > fn.Pos() && pos < fn.End() {
				enclosing = funcDeclName(fn)
				return false
			}
		}
		return true
	})

	if enclosing == "" {
		return "init/global"
	}
	return enclosing
}

// funcDeclName extracts the canonical SAAYN name for a function, including its receiver if present.
func funcDeclName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return fn.Name.Name
	}

	switch t := fn.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return "*" + ident.Name + "." + fn.Name.Name
		}
	case *ast.Ident:
		return t.Name + "." + fn.Name.Name
	}

	return fn.Name.Name
}
