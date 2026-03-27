package validator

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Result holds the outcome of a validation gate check
type Result struct {
	Success bool
	Output  string // Combined stdout/stderr from the compiler
}

// CheckIntegrity runs the Go toolchain against the project
func CheckIntegrity(projectRoot string) (*Result, error) {
	// 1. Run 'go vet' to check for common semantic errors
	if res := runCommand(projectRoot, "go", "vet", "./..."); !res.Success {
		return res, nil
	}

	// 2. Run 'go build' to ensure the code actually compiles
	// We use -o /dev/null (or equivalent) because we don't need the binary, just the result
	if res := runCommand(projectRoot, "go", "build", "-o", "temp_build_bin", "./..."); !res.Success {
		return res, nil
	}

	return &Result{Success: true, Output: "All checks passed"}, nil
}

// runCommand is a helper to execute shell commands and capture output
func runCommand(dir string, name string, args ...string) *Result {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	
	return &Result{
		Success: err == nil,
		Output:  out.String(),
	}
}

// ParseError cleans up compiler output to give the AI better feedback
func (r *Result) ParseError() string {
	if r.Success {
		return ""
	}
	// Logic to extract specific line/error messages could go here
	return strings.TrimSpace(r.Output)
}
