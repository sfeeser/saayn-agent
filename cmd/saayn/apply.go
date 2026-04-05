package saayn

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
	"gopkg.in/yaml.v3"

	"github.com/saayn-agent/internal/genome/surgery"
)

var applyInputFile string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a generated V5 surgery patch set to the local filesystem using batch AST splicing",
	RunE:  runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&applyInputFile, "file", "f", "patch.yaml", "Path to the patch set file")
}

func runApply(_ *cobra.Command, _ []string) error {
	fmt.Printf("🩺 Initializing V5 Batch AST Surgeon for %s...\n", applyInputFile)

	// 1. Read and Unmarshal
	yamlBytes, err := os.ReadFile(applyInputFile)
	if err != nil {
		return fmt.Errorf("failed to read patch file: %w", err)
	}

	var patchSet surgery.SurgeryPatchSet
	if err := yaml.Unmarshal(yamlBytes, &patchSet); err != nil {
		return fmt.Errorf("failed to parse patch YAML: %w", err)
	}

	// 2. Safety Gate
	if patchSet.Review == nil || patchSet.Review.Status != surgery.StatusApproved {
		return fmt.Errorf("🛑 safety abort: patch set is not marked as 'approved'")
	}

	// 3. Group records by file to prevent redundant I/O and backup collisions
	fileGroups := make(map[string][]surgery.PatchRecord)
	for _, rec := range patchSet.Records {
		if rec.Action == surgery.ActionNone {
			continue
		}
		fileGroups[rec.TargetNode.FilePath] = append(fileGroups[rec.TargetNode.FilePath], rec)
	}

	if len(fileGroups) == 0 {
		fmt.Println("✅ No changes required. All records were marked as 'none'.")
		return nil
	}

	// 4. Execute Batch Surgery
	for path, records := range fileGroups {
		fmt.Printf("📂 Processing file: %s (%d changes)\n", path, len(records))

		// A. Read once
		currentSrc, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// B. Backup once
		backupPath := path + ".bak"
		if err := os.WriteFile(backupPath, currentSrc, 0644); err != nil {
			fmt.Printf("   ⚠️  Warning: Backup failed for %s\n", path)
		}

		// C. Apply all records for THIS file in memory
		for _, rec := range records {
			// Strict per-record validation at apply-time
			if rec.Action != surgery.ActionReplaceNode {
				return fmt.Errorf("[%s]: unsupported action %q", rec.TargetNode.PublicID, rec.Action)
			}
			if rec.Status != surgery.StatusApproved {
				return fmt.Errorf("[%s]: record is not approved", rec.TargetNode.PublicID)
			}
			if rec.TargetNode.NodeType != "function" {
				return fmt.Errorf("[%s]: V5 apply only supports function targets", rec.TargetNode.PublicID)
			}

			fmt.Printf("   ✂️  Splicing: %s\n", rec.TargetNode.PublicID)

			// We pass currentSrc into the splice and update it with the result
			// This allows multiple edits to accumulate in the buffer
			currentSrc, err = spliceAST(path, rec.TargetNode.PublicID, currentSrc, rec.NewCode)
			if err != nil {
				return fmt.Errorf("splice failed for %s: %w", rec.TargetNode.PublicID, err)
			}
		}

		// D. Write once
		if err := os.WriteFile(path, currentSrc, 0644); err != nil {
			return fmt.Errorf("failed to write modified file %s: %w", path, err)
		}
		fmt.Printf("   ✅ Successfully updated %s\n", filepath.Base(path))
	}

	fmt.Println("\n🎉 Batch Surgery Complete. Files updated on disk.")
	return nil
}

// --- Hardened AST Splicing Logic ---

func spliceAST(filePath string, publicID string, src []byte, newCode string) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("AST parse failed: %w", err)
	}

	var targetNode ast.Node
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if buildFuncPublicID(file, fn, filePath) == publicID {
				targetNode = n
				return false
			}
		}
		return true
	})

	if targetNode == nil {
		return nil, fmt.Errorf("node %s not found in AST", publicID)
	}

	tFile := fset.File(targetNode.Pos())
	startOffset := tFile.Offset(targetNode.Pos())
	endOffset := tFile.Offset(targetNode.End())

	// Perform byte-level surgery
	newSrc := make([]byte, 0, len(src)-int(endOffset-startOffset)+len(newCode))
	newSrc = append(newSrc, src[:startOffset]...)
	newSrc = append(newSrc, []byte(newCode)...)
	newSrc = append(newSrc, src[endOffset:]...)

	// 1. Initial format/syntax check
	formatted, err := format.Source(newSrc)
	if err != nil {
		return nil, fmt.Errorf("formatting failed: %w", err)
	}

	// 2. Resolve imports
	final, err := imports.Process(filePath, formatted, nil)
	if err != nil {
		return nil, fmt.Errorf("import resolution failed: %w", err)
	}

	// 3. FINAL VALIDATION: Parse the result one last time
	verifyFset := token.NewFileSet()
	if _, err := parser.ParseFile(verifyFset, filePath, final, parser.AllErrors); err != nil {
		return nil, fmt.Errorf("post-splice parse validation failed: %w", err)
	}

	return final, nil
}

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
