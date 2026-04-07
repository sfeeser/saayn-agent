package validator

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

// PhysicsAudit performs a syntactic "Physics Audit" on generated Go code as
// part of the Surgical Inner Loop. It ensures the code is parseable by the
// Go compiler before it is allowed to be written to disk.
func PhysicsAudit(code string) error {
	// Reject ghost generations early
	if len(strings.TrimSpace(code)) == 0 {
		return fmt.Errorf("physics audit failed: generated code is empty or only whitespace")
	}

	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "generated.go", code, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("physics audit failed: invalid Go syntax\n%w", err)
	}

	return nil
}
