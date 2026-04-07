package surgeon

import (
	"fmt"

	"github.com/sfeeser/saayn-agent/internal/validator"
	"github.com/sfeeser/saayn-agent/pkg/model"
)

// InnerLoop executes the "Surgical Inner Loop": Generate -> Physics Audit -> Cognitive Audit -> Remediate.
func InnerLoop(nodeSpec model.SpecNode, vision string) (string, error) {
	var currentCode string

	// Feedback state for the next LLM generation attempt
	var lastPhysicsErr error
	var lastCognitiveFindings string

	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// 1. DRAFTING (The Synthesizer)
		code, err := draftNode(nodeSpec, vision, lastPhysicsErr, lastCognitiveFindings)
		if err != nil {
			return "", fmt.Errorf("generation failed at LLM tier on attempt %d: %w", attempt, err)
		}
		currentCode = code

		// Reset feedback state for this new attempt
		lastPhysicsErr = nil
		lastCognitiveFindings = ""

		// 2. THE PHYSICS AUDIT (The Syntactic Gatekeeper)
		if err := validator.PhysicsAudit(currentCode); err != nil {
			fmt.Printf("    ⚠️  PHYSICS FAILURE (Attempt %d/%d): %v\n", attempt, maxAttempts, err)
			lastPhysicsErr = err
			continue // Remediate
		}
		fmt.Printf("    ✅ PHYSICS PASS\n")

		// 3. THE COGNITIVE AUDIT (The Intent Gatekeeper)
		// This uses a fast LLM to verify the code matches the Spec's "Logic" and the "Vision".
		findings, passed, err := cognitiveAudit(nodeSpec, vision, currentCode)
		if err != nil {
			return "", fmt.Errorf("cognitive audit system failure: %w", err)
		}

		if !passed {
			fmt.Printf("    🧠 INTENT DRIFT DETECTED (Attempt %d/%d):\n       └─ %s\n", attempt, maxAttempts, findings)
			lastCognitiveFindings = findings
			continue // Remediate
		}
		fmt.Printf("    ✅ COGNITIVE PASS: Node intent verified.\n")

		// 4. SUCCESS: Code is both structurally sound and intent-aligned.
		return currentCode, nil
	}

	// 5. Remediation Exhausted
	return "", fmt.Errorf("remediation exhausted for '%s' after %d attempts. Last Physics: %v | Last Cognitive: %s",
		nodeSpec.Name, maxAttempts, lastPhysicsErr, lastCognitiveFindings)
}

// -----------------------------------------------------------------------------
// PLACEHOLDER INTERFACES (To be implemented in materialize.go)
// -----------------------------------------------------------------------------
