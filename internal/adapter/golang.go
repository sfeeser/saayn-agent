package adapter

// SAAYN:CHUNK_START:golang-adapter-imports-v1-h1i2j3k4
// BUSINESS_PURPOSE: Imports the native Go parser and AST tools to perform level 1 syntax validation.
// SPEC_LINK: SpecBook v1.7 Chapter 4 & 6
import (
	"bytes"
	"go/parser"
	"go/token"
	"os/exec"
)
// SAAYN:CHUNK_END:golang-adapter-imports-v1-h1i2j3k4

// SAAYN:CHUNK_START:golang-adapter-struct-v1-l5m6n7o8
// BUSINESS_PURPOSE: Implements the Adapter interface for the Go language.
// SPEC_LINK: SpecBook v1.7 Chapter 4
type GoAdapter struct{}

func (g *GoAdapter) Name() string {
	return "go"
}

func (g *GoAdapter) CommentPrefix() string {
	return "//"
}
// SAAYN:CHUNK_END:golang-adapter-struct-v1-l5m6n7o8

// SAAYN:CHUNK_START:golang-syntax-check-v1-p9q0r1s2
// BUSINESS_PURPOSE: Level 1 SyntaxCheck. Uses go/parser to ensure the chunk is syntactically valid Go code.
// SPEC_LINK: SpecBook v1.7 Chapter 4 (Level 1: Parse-only)
func (g *GoAdapter) SyntaxCheck(code string) error {
	fset := token.NewFileSet()
	// We parse as a 'File' to ensure top-level declarations are valid, 
	// but we use ParseExpr or custom logic if we are chunking inside functions.
	// For universal safety, we check if it parses as a valid block or decl set.
	_, err := parser.ParseDir(fset, "", nil, parser.AllErrors)
	
	// Refined Level 1: Ensure the chunk doesn't break the parser.
	// We wrap in a dummy package to check partial validity if it's a snippet.
	dummy := "package main\n" + code
	_, err = parser.ParseFile(fset, "chunk_test.go", dummy, 0)
	if err != nil {
		return err
	}
	return nil
}
// SAAYN:CHUNK_END:golang-syntax-check-v1-p9q0r1s2

// SAAYN:CHUNK_START:golang-formatter-v1-t3u4v5w6
// BUSINESS_PURPOSE: Executes 'go fmt' logic on the chunk to ensure the codebase remains idiomatic after AI edits.
// SPEC_LINK: SpecBook v1.7 Chapter 4 & 6
func (g *GoAdapter) Format(code string) (string, error) {
	cmd := exec.Command("gofmt")
	cmd.Stdin = bytes.NewBufferString(code)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return code, err // Return unformatted code if gofmt fails
	}
	return out.String(), nil
}

func init() {
	Register(&GoAdapter{})
}
// SAAYN:CHUNK_END:golang-formatter-v1-t3u4v5w6
