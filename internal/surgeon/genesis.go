package surgeon

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/sfeeser/saayn-agent/internal/genome"
	"github.com/sfeeser/saayn-agent/internal/scanner"
	"github.com/sfeeser/saayn-agent/pkg/model"
)

// ExecuteGenesis implements the Greenfield Protocol: materializes a full Go project
// from a Markdown vision and YAML specbook using the Surgical Inner Loop.
// If a node fails, partial state remains on disk to allow for manual inspection or resumption.
func ExecuteGenesis(visionPath, specPath, targetDir string) error {
	// 1. Directory Setup
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	// 2. Resource Ingestion
	visionBytes, err := os.ReadFile(visionPath)
	if err != nil {
		return fmt.Errorf("failed to read vision file %s: %w", visionPath, err)
	}
	// Senior Tip: Don't just check length; check for actual content
	if len(bytes.TrimSpace(visionBytes)) == 0 {
		return fmt.Errorf("vision file is empty: intent is required for genesis")
	}

	spec, err := loadSpecbook(specPath)
	if err != nil {
		return fmt.Errorf("failed to load specbook: %w", err)
	}

	// 3. Fail-fast Validation
	if spec.Project == "" {
		return fmt.Errorf("specbook validation failed: project name is required")
	}
	if len(spec.Nodes) == 0 {
		return fmt.Errorf("specbook validation failed: no nodes defined")
	}

	// 4. Registry (Birth Certificate) - Using Constructor for future-proofing
	reg := model.NewRegistry(spec.Project, targetDir)

	// 5. Deterministic Build Order (Topological Sort)
	orderedNodes, err := calculateBuildOrder(spec)
	if err != nil {
		return fmt.Errorf("dependency resolution failed: %w", err)
	}
	fmt.Printf("🏗️  Build Order Resolved: %v\n", orderedNodes)

	// 6. Materialization Loop with Surgical Inner Loop
	// Create a local map to track paths by Node Name for instant lookup
	builtPaths := make(map[string]string)

	for _, nodeSpec := range orderedNodes {
		fmt.Printf("\n🌱 MATERIALIZING: %s\n", nodeSpec.Name)

		// --- CONTEXT HYDRATION START ---
		var contextBuilder strings.Builder
		contextBuilder.WriteString(string(visionBytes)) // Always include the Soul

		if len(nodeSpec.DependsOn) > 0 {
			contextBuilder.WriteString("\n\n=== DEPENDENCY CONTEXT ===\n")
			contextBuilder.WriteString("The following dependencies exist. Use their exact signatures:\n")

			for _, depName := range nodeSpec.DependsOn {
				// 1. Look up the path using our local tracker
				depPath, exists := builtPaths[depName]
				if exists {
					// 2. Use the AST Scanner to get just the signature, saving tokens
					signatures, err := scanner.ExtractSignatures(depPath)
					if err == nil {
						contextBuilder.WriteString(fmt.Sprintf("\n// %s\n%s\n", depName, signatures))
					}
				}
			}
		}
		// --- CONTEXT HYDRATION END ---

		// Pass the FULL string (Vision + AST Signatures) into the InnerLoop
		code, err := InnerLoop(nodeSpec, contextBuilder.String())
		if err != nil {
			return fmt.Errorf("failed to materialize node %s: %w", nodeSpec.Name, err)
		}

		fullPath := filepath.Join(targetDir, nodeSpec.Path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
		}

		if err := os.WriteFile(fullPath, []byte(code), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", fullPath, err)
		}

		exec.Command("gofmt", "-w", fullPath).Run()

		// 7. Metadata Generation
		newNode, err := MaterializeNode(nodeSpec, code)
		if err != nil {
			return fmt.Errorf("failed to create metadata for node %s: %w", nodeSpec.Name, err)
		}

		reg.Nodes[newNode.UUID] = newNode

		// Track the path so dependents can find it in the next loop iteration!
		builtPaths[nodeSpec.Name] = fullPath
	}

	// 8. Persist Genome
	genomePath := filepath.Join(targetDir, "genome.json")
	if err := genome.SaveRegistry(reg, genomePath); err != nil {
		return fmt.Errorf("failed to save genome registry: %w", err)
	}

	return nil
}

// calculateBuildOrder performs a topological sort on the spec nodes
// to ensure dependencies are materialized before their dependents.
func calculateBuildOrder(spec *model.Specbook) ([]model.SpecNode, error) {
	// 1. Map for quick lookup
	nodeMap := make(map[string]model.SpecNode)
	for _, n := range spec.Nodes {
		nodeMap[n.Name] = n
	}

	// 2. Adjacency list for dependencies
	// (Assumes SpecNode has a 'DependsOn []string' field)
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	var result []model.SpecNode

	var visit func(string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("cycle detected in dependencies at: %s", name)
		}
		if !visited[name] {
			temp[name] = true
			node, exists := nodeMap[name]
			if !exists {
				return fmt.Errorf("dependency not found in spec: %s", name)
			}

			// Recursively visit dependencies first
			for _, dep := range node.DependsOn {
				if err := visit(dep); err != nil {
					return err
				}
			}

			visited[name] = true
			temp[name] = false
			result = append(result, node)
		}
		return nil
	}

	// 3. Kick off the sort for all nodes
	for _, n := range spec.Nodes {
		if !visited[n.Name] {
			if err := visit(n.Name); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// loadSpecbook reads and parses the YAML skeleton.
func loadSpecbook(path string) (*model.Specbook, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read specbook: %w", err)
	}

	var spec model.Specbook
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("invalid specbook YAML: %w", err)
	}

	return &spec, nil
}
