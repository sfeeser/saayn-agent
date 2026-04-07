package saayn

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/sfeeser/saayn-agent/internal/astutil"
	"github.com/sfeeser/saayn-agent/internal/genome/surgery"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var graphFile string
var graphDir string

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Analyze the blast radius of a planned surgery and hydrate context source code",
	RunE:  runGraph,
}

func init() {
	rootCmd.AddCommand(graphCmd)
	graphCmd.Flags().StringVarP(&graphFile, "file", "f", "surgery.yaml", "Path to the surgery plan")
	graphCmd.Flags().StringVarP(&graphDir, "dir", "d", ".", "Project root directory to trace")
}

func runGraph(cmd *cobra.Command, args []string) error {
	fmt.Printf("🕸️ Initializing V4 Impact Engine for %s...\n", graphFile)

	// 1. Read and validate the drafted surgery plan
	yamlBytes, err := os.ReadFile(graphFile)
	if err != nil {
		return fmt.Errorf("failed to read plan: %w\n💡 Hint: Run 'saayn draft' first", err)
	}

	var plan surgery.SurgeryPlan
	if err := yaml.Unmarshal(yamlBytes, &plan); err != nil {
		return fmt.Errorf("failed to parse plan YAML: %w", err)
	}

	targetID := plan.Target.PublicID
	if targetID == "" {
		return fmt.Errorf("invalid surgery plan: missing target PublicID")
	}

	fmt.Printf("🔎 Tracing callers for target: %s\n", targetID)

	// 2. Delegate to our shared AST Tracer
	callers, err := astutil.TraceCallers(graphDir, targetID)
	if err != nil {
		return fmt.Errorf("impact analysis failed: %w", err)
	}

	// 3. Hydrate Source Code and map to Schema
	seenCallers := make(map[string]bool)
	var newContext []surgery.ContextNode

	for _, caller := range callers {
		// Dedupe by specific caller function, not just the whole file
		dedupeKey := caller.FilePath + ":" + caller.CallingFunction
		if seenCallers[dedupeKey] {
			continue
		}
		seenCallers[dedupeKey] = true

		fmt.Printf("   ↳ Hydrating source code for caller: %s() in %s\n", caller.CallingFunction, filepath.Base(caller.FilePath))

		// The Bridge: Extract the actual source code of the calling function
		sourceCode, pkgName, err := extractCallerSource(caller.FilePath, caller.CallingFunction)
		if err != nil {
			fmt.Printf("      ⚠️ Warning: Failed to hydrate source for %s: %v\n", caller.CallingFunction, err)
			continue
		}

		// Reconstruct the canonical SAAYN PublicID (pkg.Receiver.Func[file.go])
		canonicalID := fmt.Sprintf("%s.%s[%s]", pkgName, caller.CallingFunction, filepath.Base(caller.FilePath))

		// Map directly to the strict Schema
		ctxNode := surgery.ContextNode{
			TargetAnchor: surgery.TargetAnchor{
				PublicID: canonicalID,
				FilePath: caller.FilePath,
				NodeType: "function",
			},
			Reason:       surgery.ReasonKnownCaller,
			ReasonDetail: fmt.Sprintf("Calls target at line %d", caller.LineNumber),
			SourceCode:   sourceCode,
		}

		newContext = append(newContext, ctxNode)
	}

	// 4. Save the hydrated plan back to disk
	plan.Context = append(plan.Context, newContext...)

	outBytes, err := yaml.Marshal(&plan)
	if err != nil {
		return fmt.Errorf("failed to marshal updated plan: %w", err)
	}

	if err := os.WriteFile(graphFile, outBytes, 0644); err != nil {
		return fmt.Errorf("failed to save updated plan: %w", err)
	}

	if len(newContext) == 0 {
		fmt.Printf("✅ Impact Analysis Complete. Node is isolated. No extra context needed.\n")
	} else {
		fmt.Printf("✅ Impact Analysis Complete. Fully hydrated %d caller nodes into %s\n", len(newContext), graphFile)
	}

	return nil
}

// --- AST JIT Hydration ---

// extractCallerSource parses a file and extracts the exact string block of a specific function.
// It also returns the package name so we can build a canonical PublicID.
func extractCallerSource(filePath, funcName string) (sourceCode string, pkgName string, err error) {
	srcBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("could not read file: %w", err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, srcBytes, 0)
	if err != nil {
		return "", "", fmt.Errorf("could not parse file: %w", err)
	}

	pkgName = file.Name.Name
	var foundNode ast.Node

	// Walk the AST to find the specific function declaration
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			// In V4 we match against the enclosing name (e.g. "runGraph" or "*User.Save")
			if getFuncDeclName(fn) == funcName {
				foundNode = fn
				return false // stop searching
			}
		}
		return true
	})

	if foundNode == nil {
		return "", pkgName, fmt.Errorf("function '%s' not found in AST", funcName)
	}

	// Extract the exact byte slice using the token FileSet
	start := fset.Position(foundNode.Pos()).Offset
	end := fset.Position(foundNode.End()).Offset

	if start < 0 || end > len(srcBytes) || start > end {
		return "", pkgName, fmt.Errorf("invalid byte offsets computed for AST node")
	}

	return string(srcBytes[start:end]), pkgName, nil
}

// getFuncDeclName mimics astutil.funcDeclName to ensure we match receiver formats correctly
func getFuncDeclName(fn *ast.FuncDecl) string {
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
