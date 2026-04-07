package surgeon

import (
	"fmt"

	"github.com/sfeeser/saayn-agent/pkg/model"
)

// CalculateBuildOrder performs a topological sort on the nodes defined in the specbook.
// It returns the nodes in a safe build order: dependencies appear before the nodes
// that depend on them. This ensures foundations are materialized before structures.
func CalculateBuildOrder(spec *model.Specbook) ([]model.SpecNode, error) {
	if spec == nil || len(spec.Nodes) == 0 {
		return nil, fmt.Errorf("specbook is empty: no nodes to order")
	}

	// 1. Quick lookup by name with pre-allocated capacity
	nodeMap := make(map[string]model.SpecNode, len(spec.Nodes))
	for _, n := range spec.Nodes {
		nodeMap[n.Name] = n
	}

	visited := make(map[string]bool) // permanently visited
	temp := make(map[string]bool)    // recursion stack (cycle detection)
	var result []model.SpecNode

	var visit func(string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("dependency cycle detected involving node: %s", name)
		}
		if visited[name] {
			return nil
		}

		node, exists := nodeMap[name]
		if !exists {
			return fmt.Errorf("undefined node referenced: %s", name)
		}

		// Mark as being processed (entering recursion stack)
		temp[name] = true

		// Recursively process all dependencies first
		for _, dep := range node.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}

		// Mark as done and add to result
		temp[name] = false
		visited[name] = true
		result = append(result, node)
		return nil
	}

	// 2. Run DFS from every node to ensure full coverage
	for _, n := range spec.Nodes {
		if !visited[n.Name] {
			if err := visit(n.Name); err != nil {
				return nil, err
			}
		}
	}

	// 3. Reverse the result so dependencies come first in the slice
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}
