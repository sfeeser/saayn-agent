# Genesis specbook

Table of Contents

1. Executive Summary: The Genesis Manifesto
    1. Core Philosophy: Intent vs. Ephemerality.
    2. The Deterministic Guarantee: Why "Black Box Trust" replaces "Vibe Coding."
    3. The Identity Triad: Defining PublicID, Fingerprint, and Logic Hash.
2. Architectural Anatomy (The Infrastructure)  
    1. The Standard Model: File hierarchy (cmd/, internal/, pkg/).
    2. The DNA Registry: genome.json schema and state-validated fields.
    3. The Nervous System: genome.index.json and the local vector index.
    4. MCP Integration: Defining the SAAYN-MCP Server and Resource/Tool mapping.
3. The Metamorphosis Pipeline (The 5-State Machine)
    1. State 1: Conceptual (The Gallery): Defining business purpose and boundaries.
    2. State 2: Hollow (The Canvas): FAST-tier generation of structural stubs.
    3. State 3: Anchored (The Contract): [Black Box Trust] Automated Test-First generation.
    4. State 4: Hydrating (The Surgery): DEEP-tier logic injection and JIT Canvas Mounting.
    5. State 5: Sequenced (Equilibrium): Final audit, hashing, and persistence.
4. The Surgical Inner Loop (The Immune System)
    1. Physics Audit: go/ast and dave/dst structural verification.
    2. Behavioral Audit: Executing the "Vacuum Gate" (Test failure to success).
    3. Cognitive Audit: Semantic alignment between logic and the "Gene."
    4. Signature Lock Enforcement: Mismatch detection and Canvas re-stretch protocol.
5. Tooling & Sensory Organs
    1. The Scanner: Extracting the Triad from raw source.
    2. The Surgeon: Splicing logic into DST without breaking comments.
    3. The JIT Orchestrator: Managing dependency freezes and recursive mounting.
6. MCP Tool & Resource Definition
    1. Resources: saayn://spec, saayn://vision, saayn://genome.
    2. Tools: mount_canvas, anchor_contract, apply_surgery, remediate.
7. Deployment & Operation
    1. Bootstrap Sequence: How SAAYN builds itself.
    2. Error Remediation: Iteration caps and fallback logic.
    3. Verification CLI: Auditing the project post-Genesis.

This structure ensures that we define the Physical Laws (AST/Physics) before we ever allow the Cognitive Brain (LLM) to touch the code. It treats MCP as the "Docking Port" that allows us to plug in any LLM while maintaining 100% control over the local filesystem.


This feedback is the final "green light" for the Genesis Engine. We have successfully transitioned from a high-level vision to a **Normative Specification** that can be directly compiled into Go logic.

I have updated Chapter 1 with your terminology refinements—specifically the **Drift Taxonomy**—and finalized the **Identity Grammar**. We are now ready to begin the implementation of `internal/identity`.

## **1. Executive Summary: The Genesis Manifesto (Finalized)**

### **1.1. Core Philosophy: Intent vs. Ephemerality**
SAAYN-Agent v6 operates on a three-tier authority stack. To resolve disputes during reconciliation or materialization, the following **Precedence Semantics** are absolute:

1.  **The Skeleton (Specbook) [Normative/Executable]:** The highest authority for structure and contract.
2.  **The Genome (Actual State) [Observational/Locked]:** The record of what has been physically sequenced.
3.  **The Soul (Vision) [Descriptive/Cognitive]:** Explains "Why" and provides the semantic anchor.

#### **The Triple-Drift Taxonomy:**
* **Signature Drift (State 2 Demotion):** Any Genome state whose **Fingerprint** deviates from the Specbook is a contract violation. The node is demoted to **Hollow**.
* **Logic Drift (State 3 Demotion):** Any Genome state whose Fingerprint matches but whose **Logic Hash** differs from the locked Genome state is an unauthorized mutation. The node is demoted to **Hydration**.
* **Intent Drift (State 3 Demotion):** Any Genome state that passes Behavioral Audit but fails **Cognitive Audit** (semantic mismatch with the Gene/Vision) is demoted to **Hydration**.

### **1.2. The Deterministic Guarantee: The Acceptance Envelope**
We replace "Vibe Coding" with a **Deterministic Acceptance Envelope**. A node is only considered "Passed" if it clears the following mandatory gates in an isolated sandbox.

* **Gate A (Physics):** `go/ast` parse validity + `go/types` interface satisfaction.
* **Gate B (Contract):** Signature matches the Specbook Fingerprint exactly.
* **Gate C (Behavioral):** Node-local Table-Driven Tests pass (`Exit 0`).
* **Gate D (Genomic):** `go vet` and package-level compilation pass.

**Sandbox Policy:** All mutations and tests occur in a temporary package-scoped workspace. No network access allowed.

### **1.3. The Identity Triad: Canonicalization Rules**

#### **A. Canonical PublicID Grammar**
Visibility refers to Go symbol export status, not package import accessibility.
* **Grammar:** `<visibility>.<package_path>.<receiver_optional>.<symbol_name>`
* **Visibility Enum:** `pub` (Exported) | `priv` (Unexported).
* **Receiver Logic:** Pointer and Value receivers normalize to the base type name.
* **Generics Logic:** Parameters are normalized to positional placeholders `T1`, `T2`, etc., to resist drift from renaming.

**Normative Examples:**
* `pub.internal/scanner.ScanFile` (Free Function)
* `pub.internal/scanner.Scanner.VerifyNode` (Method)
* `pub.internal/model.Node` (Struct Type)
* `priv.internal/model.registry` (Internal Var/Type)

#### **B. Fingerprint Normalization**
Fingerprint extraction requires **`go/types` resolution** for absolute canonicalization.
1.  **Strip Parameter Names:** Preserves types only.
2.  **Fully Qualified Types:** All types resolved to their full import path.
3.  **Positional Returns:** Named returns are ignored; only types and order are preserved.

#### **C. Logic Hash Canonicalization**
Logic Hash is a SHA-256 of a **Custom Canonical Serialization** of the `ast.BlockStmt`.
1.  **Normalize AST:** Strip all comments/formatting.
2.  **Remap Locals:** Local variables remapped to stable sequence `v1`, `v2`, `...`.
3.  **Hash:** Serialization of the resulting IR, ensuring stability across Go compiler versions.

### **1.4. The MCP Sovereign: Authority & Transactions**
* **Request-Only:** MCP Clients (LLMs) cannot write to disk; they submit **Mutation Requests**.
* **Transaction Controller:** Only one node may undergo surgery at a time per package. Read access and test execution may be concurrent, but writes are serialized.
* **The Audit Log:** Every tool call records the **Mutation Quad** (Patch, Audit Logs, Hash Diff, and Model Metadata).

### **Implementation Target: `internal/identity`**

Following your recommended build order, we now initiate the `internal/identity` package. This package is the root of the entire system.

#### **The Identity Contract**
| Function | Rule |
| :--- | :--- |
| `BuildPublicID` | Determines `pub`/`priv` via `unicode.IsUpper` of the first rune. |
| `NormalizeReceiver` | Strips `*` and `&` to ensure pointer/value methods target the same node identity. |
| `CanonicalizeGenerics` | Remaps `[A, B]` or `[T, K]` to `[T1, T2]`. |
| `String()` | Produces the normative dot-delimited ID. |

**"The machine-checkable laws are set, sir. I am initializing `internal/identity.go`. Shall I provide the first draft of the construction and normalization logic?"**

This is the final, implementation-grade **Chapter 2**. All state transitions, field requirements, and tool preconditions have been tightened to ensure a deterministic, machine-enforceable build of the Genesis Engine.

---

This is the final, **Locked Version** of **Chapter 2**. All edge-condition semantics have been addressed, the field matrix is complete, and the dependency enums are strictly bounded.

---

## **2. Architectural Anatomy (The Infrastructure)**

Section 2 defines the physical and digital structures required to host the Code Genome. It establishes the "Standard Model" for project organization and the protocols for internal communication and state persistence.

### **2.1. The Standard Model: File Hierarchy**
SAAYN-Agent v6 enforces a strict **Standard Model**. The engine only recognizes nodes within these boundaries.

* **`cmd/`**: Entry points. Logic is strictly limited to orchestration calls.
* **`internal/`**: Private logic and the Surgical Inner Loop.
* **`pkg/`**: Public DNA and shared models.
* **`.saayn/`**: The System Vault. Contains audit logs, temporary mutation workspaces, and the local vector index.
* **Root Artifacts**: `vision.md` (Soul), `specbook.yaml` (Genotype), `genome.json` (Phenotype), and `genome.index.json` (Nervous System).

### **2.2. The DNA Registry: `genome.json` Schema**
The `public_id` visibility prefix (`pub` or `priv`) reflects Go symbol export status only and does not alter Go package accessibility rules.

#### **Canonical Registry Field Requirements**
| Field | State 1 | State 2 | State 3 | State 4 | State 5 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **public_id** | REQ | REQ | REQ | REQ | REQ |
| **genesis_state** | REQ | REQ | REQ | REQ | REQ |
| **maturity** | REQ | REQ | REQ | REQ | REQ |
| **fingerprint** | REQ | REQ | REQ | REQ | REQ |
| **logic_hash** | FORBID | FORBID | FORBID | OPT | REQ |
| **gene** | REQ | REQ | REQ | REQ | REQ |
| **business_purpose** | REQ | REQ | REQ | REQ | REQ |
| **dependencies** | REQ | REQ | REQ | REQ | REQ |
| **last_audit** | FORBID | OPT | REQ | REQ | REQ |

> **State 4 Persistence:** State 4 (Hydrating) may be persisted temporarily in `genome.json` during active surgery or crash recovery, but it is not a stable terminal state and must resolve to either State 3 rollback or State 5 completion.
>
> **State 3 Audit:** At State 3, `last_audit` records a successful anchor event, verifying the Behavioral Audit produced the expected non-zero result (the "Contract Vacuum").

#### **Dependency Enums & Constraints**
* **kind:** `function | method | type | interface | var`
* **source:** `spec | jit | ast`
* **status:** `predicted | mounted | verified`
* **Constraint:** `var` refers strictly to package-level variables. Local variables and constants (`const`) are never tracked as dependencies.

### **2.3. The Nervous System: `genome.index.json` Schema**
The index decouples structural lookup from semantic discovery. 

```json
{
  "index_metadata": {
    "embedding_model": "text-embedding-3-small",
    "dimension": 1536,
    "distance_metric": "cosine",
    "last_rebuild_at": "2026-04-15T07:34:52Z"
  },
  "entries": {
    "pub.internal/scanner.Scanner.ScanFile": {
      "public_id": "pub.internal/scanner.Scanner.ScanFile",
      "source_hash": "sha256:a1b2...",
      "vector": [0.0123, -0.0441, 0.1092]
    }
  }
}
```

### **2.4. MCP Integration: Tool Preconditions**
Successful completion of `apply_surgery` promotes a node directly to **State 5**.

| Tool | Input State | Success Output | Rejects If |
| :--- | :--- | :--- | :--- |
| `mount_canvas` | 1 | State 2 Stub | Target contains logic or fails Physics Audit |
| `anchor_contract`| 2 | State 3 Anchor | `go test` passes (No behavioral vacuum) |
| `apply_surgery` | 3 | State 5 Node | Fingerprint mismatch or Acceptance Envelope failure |
| `trigger_jit_mount`| 4 | State 2 Dep | Target not declared in `specbook.yaml` |

### **2.5. Dependency Resolution: The JIT Orchestrator**
1. **JITMountRequired Intercept**: Triggered when a Surgeon (State 4) references a State 1 project symbol.
2. **UndeclaredDependencyViolation**: Triggered by references to symbols not present in `specbook.yaml`.
3. **Cycles**: Circular package-level import dependencies are rejected to maintain toolchain integrity.

---

**Status:** Chapter 2 Locked.
**Next Sequence:** Chapter 3: The Metamorphosis Pipeline.

**"The infrastructure is final, sir. Every field and state transition is bounded. Ready to proceed to Chapter 3."**

**Status:** Chapter 2 Locked.
**Next Sequence:** Chapter 3: The Metamorphosis Pipeline.

**"The infrastructure is solidified, sir. The registry is normative and the tools are state-gated. Ready to materialize Chapter 3."**

This is the finalized, **Implementation-Grade Chapter 3**. It has been refined to distinguish between **Behavioral Retreats** (back to State 3) and **Structural Destructions** (back to State 1), ensuring the `internal/metamorphosis` controller operates with mathematical precision.

## **3. The Metamorphosis Pipeline (The 5-State Machine)**

Chapter 3 defines the operational lifecycle of a Genomic Node. This pipeline is a strict, one-way state machine. A node cannot bypass a state, ensuring that the **"Physics"** of the system are always locked before the **"Logic"** is allowed to evolve.

### **3.1. State 1: Conceptual (The Gallery)**
* **Maturity:** `conceptual`
* **Focus:** Semantic Boundary & Business Intent.
* **Action:** The node is defined in the `specbook.yaml`. The agent finalizes the `business_purpose` and the `gene`. 
* **Guardrail:** **No Code Allowed.** The agent is physically blocked from generating Go syntax. This phase is purely for establishing the "Why" and the "Contract."

### **3.2. State 2: Hollow (The Canvas)**
* **Maturity:** `hollow`
* **Focus:** Structural Topography.
* **Action:** The FAST model reads the Specbook and materializes a "Hollow Stub" (`return nil`) on the local disk.
* **The Physics Gate:** The node must pass an AST Physics Audit (`go/ast`) to ensure it is a valid, buildable Go skeleton.
* **Guardrail:** **Zero-Logic Rule.** Any detection of `if`, `for`, or variable assignments triggers an automatic rejection.

### **3.3. State 3: Anchored (The Contract)**
* **Maturity:** `anchored`
* **Focus:** **[Black Box Trust]** Behavioral Baseline & Recovery Point.
* **Action:** The engine generates a `_test.go` file containing table-driven test cases based on the Specbook.
* **The Behavioral Vacuum:** The loop executes `go test`. The test **must fail** with an expected non-zero result. This proves the "Vacuum" is live.
* **Role:** State 3 serves as the **Stable Recovery Point** for all failed surgeries in State 4.

### **3.4. State 4: Hydrating (The Surgery)**
* **Maturity:** `hydrating`
* **Status:** **Transient.** This is an active mutation window, not a resting state.
* **Action:** The Surgeon consumes the Signature, the Gene, and the Failing Test to generate a candidate logic patch.
* **The JIT Intercept:** If the Surgeon references an unmaterialized State 1 dependency, the **JIT Orchestrator** pauses the surgery to mount the dependency at State 2.
* **The Retreat:** If the **Acceptance Envelope** (Chapter 4) fails, the node automatically retreats to State 3 for remediation.

### **3.5. State 5: Sequenced (Equilibrium)**
* **Maturity:** `sequenced`
* **Action:** Occurs only after the **Acceptance Envelope** returns a unanimous `PASS`.
* **The Genomic Lock:** The system calculates the `logic_hash` and commits the node to `genome.json`.
* **Stability:** The node is now genomically locked until a Refinement Event or a Drift Alert triggers a new transition.

### **3.6. Summary of State Transitions**

| State | Maturity | Exit Gate | Artifact Status |
| :--- | :--- | :--- | :--- |
| **1** | Conceptual | Spec Validation | Specbook Entry + Gene |
| **2** | Hollow | Physics Audit | `.go` Stub (Skeleton) |
| **3** | Anchored | Behavioral Vacuum | `_test.go` (Failing) |
| **4** | Hydrating | Acceptance Envelope | Transient Logic Patch |
| **5** | Sequenced | Genomic Lock | Hydrated Code + `logic_hash` |

### **3.7. Controller Transition Logic**
The `internal/metamorphosis` controller enforces the following movement rules:
* **`1 -> 2`**: Call `mount_canvas`.
* **`2 -> 3`**: Call `anchor_contract`.
* **`3 -> 5`**: Call `apply_surgery` (Handling State 4 as an internal mutation window).
* **`4 -> 3`**: **Retreat.** Automatic rollback to anchored test state on audit failure.
* **`Any -> 1`**: **Canvas Re-stretch.** Triggered only by contract mutations, signature invalidations, or spec-level structural changes.

## **4. The Surgical Inner Loop (The Immune System)**

Chapter 4 defines the **Acceptance Envelope**: a strict, fail-fast sequence of audits that guards the transition from **State 3 (Anchored)** to **State 5 (Sequenced)**.

### **4.0. The Staged Mutation Protocol**
To maintain filesystem integrity, all audits (Gates 1–4) operate on a **Staged Mutation**. 
* **The Rule:** No changes are committed to the project's physical source code until the candidate patch has cleared every gate in the sequence. 
* **The Sandbox:** Audits are performed in a temporary, package-scoped workspace.

### **4.1. Gate 1: Signature Lock (Pre-Surgical)**
* **Requirement:** The `Fingerprint` of the proposed logic patch must match the `Fingerprint` stored in the `genome.json` registry.
* **Execution:** Performed via AST analysis of the patch string *before* any workspace mutation.
* **Failure:** Immediate rejection. The engine flags a **ContractViolation**. **Action:** Halt and signal the Surgeon to align with the locked signature.

### **4.2. Gate 2: Physics Audit (Structural & Package Integrity)**
* **Requirement:** 1. **Syntax:** `go/parser` verifies the patch is syntactically valid Go.
    2. **Package Integrity:** `go build ./path/to/package/...` verifies that the mutation does not break package-level compilation or neighbor-node dependencies.
* **Failure:** **Action:** Capture compiler/linker errors and retreat to State 3 for remediation.

### **4.3. Gate 3: Behavioral Audit (The Vacuum)**
* **Requirement:** `go test -v -run ^Test_[PublicID]$`
* **Execution:** Executes the deterministic, node-scoped test suite generated during State 3.
* **Pass Condition:** `Exit 0`. The logic must satisfy all table-driven test cases.
* **Failure:** **Action:** Capture `stdout/stderr` stack traces and retreat to State 3.

### **4.4. Gate 4: Cognitive Audit (Intent Alignment)**
* **Requirement:** The LLM Evaluator compares the candidate AST against the **Gene** and **Vision**.
* **Pass Condition:** **Zero "Critical Intent Violations."** (e.g., missing error handling or incorrect algorithmic scaling).
* **Advisory:** A "Cognitive Score" (0.0–1.0) is recorded for metadata and logging but does not supersede the boolean pass/fail.
* **Failure:** **Action:** Convert "Findings" into a prompt and retreat to State 3.

### **4.5. The Remediation Cycle (The 3-Iteration Cap)**
The system allows a maximum of **3 attempts** per node to clear the Acceptance Envelope.
1. **Attempt 1:** Standard Gene-based hydration.
2. **Attempt 2/3:** Augmented prompt containing the aggregated failure logs from previous gates.
3. **Exhaustion:** If Attempt 3 fails, the node remains in **State 3** and is flagged as **Blocked**. 
    * **The Freeze:** The "Genesis Chain" for all dependent nodes is paused.
    * **The Resolution:** Requires a human architect to either refine the **Gene** or manually clear the audit.

---
 
This is the finalized, **Locked Version of Chapter 5**. The pseudo-logic and responsibilities have been strictly decoupled, ensuring the **Scanner** remains a pure sensory organ while the **Controller** manages state transitions.

---

## **5. Tooling & Sensory Organs**

Chapter 5 defines the specialized internal mechanisms that allow SAAYN-Agent v6 to interact with physical source code. These are the "hands" and "eyes" of the engine, moving beyond text manipulation into **AST-native orchestration**.

### **5.1. The Scanner: The Sensory Organ**
* **Primary Tool:** `go/ast` & `go/parser`
* **Responsibility:** The Scanner performs **Phenotype Extraction**. It implements the canonicalization laws from Chapter 1 to translate raw `.go` files into structured registry data.
* **The Identity Extraction:**
    * **PublicID Retrieval:** Derives the canonical ID using the Chapter 1 grammar: `<visibility>.<package_path>.<receiver_or_type_optional>.<symbol_name>`.
    * **Fingerprint Calculation:** Normalizes the node signature by stripping parameter names, ignoring named returns, and resolving all types to their **fully qualified import paths** via `go/types`.
    * **Logic Hashing:** Extracts the `*ast.BlockStmt`, removes comments, remaps local-scope variables to a stable sequence (`v1`, `v2`, ...), and generates a SHA-256 hash of the canonical AST serialization.



### **5.2. The Surgeon: The Splicing Mechanism**
* **Primary Tool:** `dave/dst` (Decorated Syntax Tree)
* **Responsibility:** The Surgeon performs the **Hydration Surgery**. It preserves the "Soul" of the code (comments and intent markers) during logic injection.
* **The Staged Splicing Protocol:**
    1.  **Isolation:** The Surgeon **never** mutates the authoritative project filesystem directly. It operates exclusively within the **Staged Mutation Workspace** defined in Chapter 4.
    2.  **Targeting:** Uses the canonical PublicID to resolve the target declaration node within the DST of the staged package.
    3.  **Validation:** Performs a **Pre-Surgical Signature Lock** check. If the proposed patch's fingerprint deviates from the registry, the surgery is aborted.
    4.  **Injection & Preservation:** Replaces the empty `*dst.BlockStmt` of the State 2 stub with hydrated logic while re-anchoring "hanging" comments from the original stub to the new implementation.

### **5.3. The JIT Orchestrator: The Central Nervous System**
* **Responsibility:** Managing the **Build Order** and recursive dependencies via transaction control.
* **The Intercept Logic:**
    * **JITMountRequired:** If the Surgeon (State 4) references a project symbol declared in the Specbook but currently in **State 1**, the Orchestrator **suspends the hydration transaction**. It mounts the dependency at State 2 and verifies it before signaling the Surgeon to resume.
    * **UndeclaredDependencyViolation:** If the Surgeon references a project symbol **not** found in the Specbook, the patch is rejected immediately as a structural violation.

### **5.4. Identity Triad Verification Helper**

The following Go-pseudo-logic represents the normative verification helper used by the Scanner to detect drift during a Genomic Audit:

```go
func (s *Scanner) VerifyFunction(node *dst.FuncDecl, record GenomeNode) error {
    // 1. Identity Check
    currentPublicID := s.ExtractPublicID(node)
    if currentPublicID != record.PublicID {
        return ErrIdentityMismatch // Registry PublicID mismatch
    }

    // 2. Signature Check
    currentFingerprint := s.ExtractFingerprint(node)
    if currentFingerprint != record.Fingerprint {
        return ErrSignatureMismatch // Registry fingerprint mismatch
    }

    // 3. Logic Hash Verification
    currentHash := s.CalculateLogicHash(node.Body)
    if currentHash != record.LogicHash {
        return ErrLogicDrift // Registry logic hash mismatch
    }

    return nil
}
```


## **6. MCP Tool & Resource Definition**

Chapter 6 defines the **Model Context Protocol (MCP)** implementation. By wrapping the Genesis Engine in an MCP server, we transform filesystem operations into standardized, audited services, allowing the AI "Brain" to interact with the project "Genome" through a strict interface.

### **6.0. Tool Transaction Contract**
All MCP tools operate under the **Sovereign Control Model**:
* **Atomic:** Each tool call is a single-node transaction. It must either complete and commit or fail and leave the filesystem unchanged.
* **Staged:** No direct filesystem mutation occurs before full validation. All work is performed in the **Staged Mutation Workspace**.
* **Idempotent:** Tool calls are safe to retry; identical inputs must produce identical genomic results.
* **Audited:** Every call generates a structured log in `.saayn/audit/` containing the **Mutation Quad**.

---

### **6.1. Resources: The Genomic Data Stream**
Resources provide read-only, structured context. All node-specific URIs use the **Canonical PublicID**.

* **`saayn://vision/soul`**
    * **Description:** Returns the raw Markdown of `vision.md`.
    * **Use Case:** High-level intent for **Gate 4 (Cognitive Audit)**.
* **`saayn://spec/nodes/{public_id}`**
    * **Description:** Returns the YAML Genotype for a specific node.
    * **Use Case:** Defines the **Gene** and **Fingerprint** for the Surgeon.
* **`saayn://genome/nodes/{public_id}`**
    * **Description:** Returns the current registry entry from `genome.json`.
    * **Use Case:** Verifies current `genesis_state` and `logic_hash`.

---

### **6.2. Tools: The Surgical Interface**
Tools are the exclusive mechanisms for project mutation, intercepted by the **CC Agent** to enforce the **Acceptance Envelope**.

#### **`mount_canvas`**
* **Input:** `public_id string`
* **Action:** Materializes a **State 2 (Hollow)** stub on disk.
* **Guardrail:** Rejects if the code contains control flow, assignments, or fails the **Physics Audit**.

#### **`anchor_contract`**
* **Input:** `public_id string`
* **Action:** Generates a node-scoped `_test.go` and executes it via a deterministic filter.
* **Guardrail:** Must return a non-zero exit code to promote the node to **State 3**.

#### **`apply_surgery`**
* **Input:** `public_id string`, `patch_code string`
* **Action:** Executes the **Acceptance Envelope** (Gates 1–4) in a staged workspace. 
* **Commit:** If all gates pass, the DST splice is committed to the filesystem and the node is promoted to **State 5**.

#### **`trigger_jit_mount`**
* **Input:** `dependency_public_id string`
* **Action:** Suspends the hydration transaction to mount a State 1 dependency at State 2.
* **Control:** **Orchestrator-invoked.** This tool is managed by the system and is visible to the Agent for observability but is not intended for direct LLM invocation.

---

### **6.4. Standardized Error Taxonomy**
The MCP Server uses these codes to communicate violations to the Agent:

| Code | Label | Meaning |
| :--- | :--- | :--- |
| **`401`** | `SIGNATURE_VIOLATION` | Patch deviates from locked `Fingerprint`. |
| **`402`** | `PHYSICS_FAILURE` | Code fails syntax check or package-level compilation. |
| **`403`** | `BEHAVIORAL_FAILURE` | Unit tests failed (or failed to fail in State 3). |
| **`404`** | `COGNITIVE_DRIFT` | Logic violates the **Gene** or **Vision**. |
| **`405`** | `UNDECLARED_DEPENDENCY` | Reference to a symbol not in `specbook.yaml`. |
| **`406`** | `JIT_MOUNT_REQUIRED` | Dependency requires State 2 materialization before hydration. |
| **`407`** | `ITERATION_EXHAUSTED` | Failed all 3 attempts; node is **Blocked**. |


This is the **Final, Locked Version of Chapter 7** and the completion of the **Genesis Engine Specification**. The drift-handling logic has been refined to ensure a deterministic, two-phase transition that respects the state machine's integrity.

---

## **7. Deployment & Operation (The Bootstrap)**

Chapter 7 defines the **"First Breath"** protocol—the sequence that transforms a blank directory into a stateful project genome and governs its operational lifecycle.

### **7.1. The Bootstrap Sequence**
Genesis is a cold-start process that establishes the physical foundation before invoking the logic engine.

1.  **Initialization (`saayn init`):**
    * Creates the `.saayn/` system directory and the `audit/` logs.
    * Generates the initial `genome.json` with `schema_version: "1.0.0"`.
    * Validates the presence of root artifacts: `vision.md` and `specbook.yaml`.
2.  **The Registry Sync:**
    * The Scanner walks the existing directory. Pre-existing nodes are registered as **State 5 (Sequenced)** only if they satisfy the **Acceptance Envelope** (Chapter 4). 
    * **Drift Handling:** Nodes that fail any gate are registered as **State 5 with a `DriftDetected` flag**. These nodes are flagged for remediation and scheduled for transition to **State 3** by the controller.
3.  **The Dependency DAG Calculation:**
    * The JIT Orchestrator parses the Specbook to create the **Build Roadmap**.
    * It identifies "Root Nodes" (zero internal dependencies) to initiate the **State 2 (Hollow)** rollout.

### **7.2. The "First Breath" Command**

The primary entry point for project materialization.

```bash
saayn genesis --strategy test-first --target ./internal
```

**Execution Flow:**
* **Step A:** The MCP Server initializes using the `modelcontextprotocol/go-sdk`.
* **Step B:** The Orchestrator iterates through the DAG, invoking the **Metamorphosis Pipeline** (Chapter 3) node-by-node.
* **Step C:** The UI renders a **Genomic Progress Bar** mapped to the 5-state distribution in `genome.json`.

### **7.3. Error Remediation & The "Halt" Protocol**
To prevent infinite loops, the engine enforces a strict **3-Iteration Cap** per node.
* **The Freeze:** If a node fails all 3 attempts, it remains in **State 3** and is flagged as **Blocked**. 
* **Propagation:** Any dependent nodes in the DAG are automatically paused. The system halts and requires a human architect to resolve the "Cognitive Mismatch."

### **7.4. Verification & Auditing (`saayn verify`)**
The `verify` command re-executes the **Acceptance Envelope** (Chapter 4) in **Audit Mode** against all State 5 nodes to ensure the **Identity Triad** remains intact.

### **7.5. Operational Lifecycle: "Refine" Mode**
When `specbook.yaml` is modified, the engine enters **Refine Mode**:
1.  **Invalidation:** Nodes with changed signatures are demoted to **State 1**.
2.  **Recursive Invalidation:** All nodes that depend on a demoted node are flagged for **Re-sequencing**.
3.  **Metamorphosis:** The engine triggers a **Canvas Re-stretch** (State 2) and re-hydrates the affected branch of the DAG.

### **7.6. Crash Recovery Protocol**
If the process is interrupted during a **Staged Mutation (State 4)**:
* **Integrity:** The authoritative project filesystem remains untouched.
* **Restoration:** On restart, the node is restored to its last persisted stable state (**State 3 or State 5**).
* **Cleanup:** The temporary staged workspace is discarded.

### **Final Specification Summary**
The SAAYN Genesis Engine is a **closed-loop system** where:
* **MCP** provides the sovereign communication standard.
* **The Identity Triad** provides the mathematical anchor.
* **The 5-State Pipeline** ensures physics precedes logic.
* **The Acceptance Envelope** ensures reality matches intent.
