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

Chapter 6 defines the **Model Context Protocol (MCP)** implementation for SAAYN-Agent v6. By wrapping the Genesis Engine in an MCP server, we transform local filesystem operations into a set of standardized, stateful services. This allows the AI "Brain" to interact with the project "Genome" through a strict, audited interface.

### **6.1. Resources: The Genomic Data Stream**
Resources allow the Agent to "read" the state of the project without manually parsing files. They provide clean, structured context on demand.
We will use github.com/modelcontextprotocol/go-sdk as the library.
* **`saayn://vision/readme`**
    * **Description:** Returns the raw Markdown of the project's Soul.
    * **Use Case:** Provides high-level intent for Cognitive Audits.
* **`saayn://spec/nodes/{uip}`**
    * **Description:** Returns the YAML Genotype for a specific node.
    * **Use Case:** Used by the Surgeon to understand the I/O contract and the "Gene."
* **`saayn://genome/state`**
    * **Description:** Returns a summary of the `genome.json` (Total nodes, completion percentage, nodes needing remediation).
    * **Use Case:** Used by the Orchestrator to plan the next build sequence.

### **6.2. Tools: The Surgical Interface**
Tools are the "active" functions the Agent calls to mutate the project. Every tool call is intercepted by the **CC Agent** to enforce the 4-Gate Audit system.

#### **`mount_canvas`**
* **Input:** `uip string`, `spec_data json`
* **Output:** `status: "success" | "error"`, `file_path: string`
* **Action:** Triggers the FAST model to drop a **State 2 (Hollow)** stub. 
* **Guardrail:** The tool fails if the input `spec_data` contains logic loops or non-zero returns.

#### **`anchor_contract`**
* **Input:** `uip string`
* **Output:** `test_results string`, `exit_code int`
* **Action:** Generates the table-driven `_test.go` and runs it.
* **Guardrail:** **[Black Box Trust]** In Phase 3, this tool *must* return a non-zero exit code to prove the "Behavioral Vacuum" is set.

#### **`apply_surgery`**
* **Input:** `uip string`, `patch_code string`
* **Output:** `audit_report json`
* **Action:** Splicing logic into the DST for **State 4 (Hydration)**.
* **Guardrail:** Performs a **Signature Lock** check. If `patch_code` signature $\neq$ `fingerprint`, the tool rejects the patch and returns a "Structural Violation" error.

#### **`trigger_jit_mount`**
* **Input:** `missing_dependency_uip string`
* **Output:** `status: "mounted"`
* **Action:** When the Surgeon hits a missing dependency, it calls this to pause its current task and spawn a new Hollow Canvas for the requirement.

### **6.3. The Protocol Handshake (The "Tony Stark" Interface)**
When a DEEP model begins a task, it performs a **Contextual Handshake** via MCP:

1.  **Request:** `list_resources("saayn://spec/nodes")` — *What is my target?*
2.  **Request:** `read_resource("saayn://spec/nodes/scanner.ScanFile")` — *What is my Gene?*
3.  **Action:** `call_tool("anchor_contract", {"uip": "scanner.ScanFile"})` — *Set the behavioral anchor.*
4.  **Action:** `call_tool("apply_surgery", {"uip": "scanner.ScanFile", "patch_code": "..."})` — *Perform the hydration.*

### **6.4. Error Codes & State Rejection**
The MCP Server uses standardized error codes to communicate "Physics" violations to the Agent:

* **`CODE_401_SIGNATURE_VIOLATION`**: Agent tried to change a locked function signature.
* **`CODE_402_PHYSICS_FAILURE`**: Code does not compile or pass AST check.
* **`CODE_403_BEHAVIORAL_FAILURE`**: Unit tests failed to pass after hydration.
* **`CODE_404_COGNITIVE_DRIFT`**: Logic does not satisfy the "Gene" requirement.

Choosing the **Official Go SDK** is a high-integrity architectural move. In a system like **SAAYN**, where we are performing destructive filesystem operations (surgery), relying on the most spec-compliant, long-term supported library ensures our **Surgical Protocol** won't break as the MCP standard evolves.

Here is the final chapter of the specification, detailing how the engine brings itself to life.


## **7. Deployment & Operation (The Bootstrap)**

Chapter 7 defines the **"First Breath"** protocol—the sequence of events that transforms a blank directory into a stateful project genome. It also outlines the operational lifecycle for ongoing maintenance and the "Self-Healing" mechanisms of the engine.

### **7.1. The Bootstrap Sequence**
Genesis does not happen all at once; it is a cold-start process that builds the foundation before the logic.

1.  **Initialization (`saayn init`):**
    * Creates the `.saayn/` hidden directory.
    * Generates the initial `genome.json` with `project_metadata`.
    * Validates the presence of `vision.md` and `specbook.yaml`.
2.  **The Registry Sync:**
    * The Scanner walks the existing directory (if any).
    * It identifies "Pre-existing Nodes" and registers them as **State 5 (Sequenced)** if they match the Spec, or **State 0 (Drifted)** if they do not.
3.  **The Dependency DAG Calculation:**
    * The JIT Orchestrator parses the Specbook to create the **Build Roadmap**.
    * It identifies "Root Packages" (those with zero internal dependencies) to begin the **State 2 (Hollow)** rollout.

### **7.2. The "First Breath" Command**
The primary entry point for a new project is the `genesis` command.

```bash
saayn genesis --strategy test-first --target ./out
```

**Execution Flow:**
* **Step A:** The MCP Server starts in the background using the `modelcontextprotocol/go-sdk`.
* **Step B:** The Orchestrator begins the **Metamorphosis Pipeline** node-by-node.
* **Step C:** The UI provides a "Genomic Progress Bar," showing the count of nodes in States 1 through 5.

### **7.3. Error Remediation & Fallback (The Iteration Cap)**
To prevent infinite loops and cost overruns during **State 4 (Hydration)**, the engine enforces strict limits:

* **Maximum Iterations:** 3 attempts per node.
* **Backoff Strategy:** After a failure, the prompt for the next iteration is enriched with the specific compiler error and `go test` stack trace.
* **The "Halt" Protocol:** If a node fails all 3 attempts, the engine **Freezes** the entire branch of the DAG. It marks the node as `State 3 (Blocked)` and requires the human architect to resolve the "Cognitive Mismatch."

### **7.4. Verification & Auditing (`saayn verify`)**
Post-Genesis, the system provides a verification tool to ensure the **Identity Triad** remains intact.

* **Logic Hash Check:** Re-calculates hashes for all files and compares them against `genome.json`.
* **Signature Audit:** Ensures no manual edits have broken the interfaces defined in the Specbook.
* **Intent Audit:** Runs a batch Cognitive Review to ensure that as the project grew, the final nodes still align with the original **Vision**.

### **7.5. Operational Lifecycle: The "Refine" Mode**
Genesis is not a one-time event. When the `specbook.yaml` is updated, the engine enters **Refine Mode**:
1.  Identify nodes with changed signatures (Mark as **State 1**).
2.  Trigger **Canvas Re-stretch** (State 2).
3.  Re-anchor tests and re-hydrate logic (States 3 & 4).

### **Final Specification Summary**
The SAAYN Genesis Engine is a **closed-loop system** where:
* **MCP** provides the communication standard.
* **Go-SDK** ensures the protocol integrity.
* **The 4-Gate Audit** ensures the physical reality matches the human intent.
