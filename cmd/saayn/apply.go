package saayn

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

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/saayn-agent/internal/genome/surgery"
)

var applyInputFile string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a generated surgery patch to the local filesystem using AST splicing",
	RunE:  runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&applyInputFile, "file", "f", "patch.yaml", "Path to the patch file")
}

func runApply(cmd *cobra.Command, args []string) error {
	fmt.Printf("🩺 Initializing AST Surgeon for %s...\n", applyInputFile)

	// 1. Read the Patch File
	yamlBytes, err := os.ReadFile(applyInputFile)
	if err != nil {
		return fmt.Errorf("failed to read patch file: %w", err)
	}

	var patch surgery.SurgeryPatch
	if err := yaml.Unmarshal(yamlBytes, &patch); err != nil {
		return fmt.Errorf("failed to parse patch YAML: %w", err)
	}

	// 2. Validate Patch Content
	if patch.TargetNode.FilePath == "" || patch.TargetNode.PublicID == "" {
		return fmt.Errorf("invalid patch: missing file_path or public_id")
	}
	if patch.Action != surgery.ActionReplaceNode {
		return fmt.Errorf("unsupported patch action: %s", patch.Action)
	}
	if strings.TrimSpace(patch.NewCode) == "" {
		return fmt.Errorf("invalid patch: new_code is empty")
	}
	if patch.TargetNode.NodeType != "function" {
		return fmt.Errorf("V3 apply only supports function targets, got %s", patch.TargetNode.NodeType)
	}

	fmt.Printf("🎯 Target isolated: %s\n", patch.TargetNode.PublicID)
	fmt.Printf("📝 Rationale: %s\n", patch.Rationale)

	// 3. Read Original Source (Before any destructive action)
	originalBytes, err := os.ReadFile(patch.TargetNode.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read original file for backup/surgery: %w", err)
	}

	// 4. Create a Safety Backup immediately
	backupPath := patch.TargetNode.FilePath + ".bak"
	if err := os.WriteFile(backupPath, originalBytes, 0644); err != nil {
		fmt.Printf("⚠️ Warning: Could not create backup at %s\n", backupPath)
		// For V3 we continue, but we warn the user.
	} else {
		fmt.Printf("🛡️  Backup created at: %s\n", backupPath)
	}

	// TODO (V4): Compute logic hash of the located targetNode and compare against patch.TargetNode.LogicHash
	// to ensure the file hasn't drifted since the plan was created.

	// 5. Perform the AST Splice using exact identity matching
	fmt.Println("✂️  Splicing AST...")
	patchedBytes, err := spliceAST(patch.TargetNode.FilePath, patch.TargetNode.PublicID, originalBytes, patch.NewCode)
	if err != nil {
		return fmt.Errorf("AST splicing failed: %w", err)
	}

	// 6. Write the mutated file back to disk
	if err := os.WriteFile(patch.TargetNode.FilePath, patchedBytes, 0644); err != nil {
		return fmt.Errorf("failed to write patched file: %w", err)
	}

	fmt.Println("✅ Surgery successful! Patch applied to disk.")
	return nil
}

// --- AST Splicing Logic ---

// spliceAST reads a Go AST, reconstructs the identities, finds the target matching the full publicID,
// replaces its exact byte range with newCode, and runs go/format on the result.
func spliceAST(filePath string, publicID string, src []byte, newCode string) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse target file AST: %w", err)
	}

	var targetNode ast.Node

	// Walk the AST to find the specific function declaration by full identity
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		if fn, ok := n.(*ast.FuncDecl); ok {
			// Reconstruct full PublicID for the node
			candidateID := buildFuncPublicID(file, fn, filePath)
			if candidateID == publicID {
				targetNode = n
				return false // Stop traversal, we found Patient Zero
			}
		}

		return true
	})

	if targetNode == nil {
		return nil, fmt.Errorf("could not locate target node '%s' in the AST of %s", publicID, filePath)
	}

	// Retrieve exact byte offsets safely
	tFile := fset.File(targetNode.Pos())
	if tFile == nil {
		return nil, fmt.Errorf("failed to resolve token file for target node")
	}

	startOffset := tFile.Offset(targetNode.Pos())
	endOffset := tFile.Offset(targetNode.End())

	// Sanity bounds checks
	if startOffset < 0 || endOffset < startOffset || endOffset > len(src) {
		return nil, fmt.Errorf("invalid AST offsets computed for target node: start=%d, end=%d, fileLen=%d", startOffset, endOffset, len(src))
	}

	// Construct the new file by sandwiching the new code between the old unmutated parts
	var buf bytes.Buffer
	buf.Write(src[:startOffset])
	buf.WriteString(newCode)
	buf.Write(src[endOffset:])

	// Run go/format (equivalent to 'go fmt') to ensure perfect indentation and check for syntax errors
	formattedBytes, err := format.Source(buf.Bytes())
	if err != nil {
		// Fail hard. If the LLM wrote bad syntax, we do NOT want to apply it to disk.
		return nil, fmt.Errorf("code formatting/syntax verification failed: %w\nLLM generated invalid Go code", err)
	}

	return formattedBytes, nil
}

// buildFuncPublicID reconstructs the canonical SAAYN PublicID for a function or method AST node.
// Format: "pkg.[*Receiver.]SymbolName[file.go]"
func buildFuncPublicID(f *ast.File, fn *ast.FuncDecl, filePath string) string {
	pkgName := f.Name.Name
	baseFile := filepath.Base(filePath)

	var receiver string
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		switch exp := fn.Recv.List[0].Type.(type) {
		case *ast.Ident:
			receiver = exp.Name + "."
		case *ast.StarExpr:
			if ident, ok := exp.X.(*ast.Ident); ok {
				receiver = "*" + ident.Name + "."
			}
		}
	}

	return fmt.Sprintf("%s.%s%s[%s]", pkgName, receiver, fn.Name.Name, baseFile)
}
