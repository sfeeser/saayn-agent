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


## **1. Executive Summary: The Genesis Manifesto**

### **1.1. Core Philosophy: Intent vs. Ephemerality**
In the traditional software lifecycle, **Source Code** is often mistaken for the truth. In reality, code is merely a transient byproduct—an ephemeral artifact prone to drift, decay, and technical debt. The true "Source of Truth" is **Intent**.

The **Genesis Engine** operates on a revolutionary hierarchy:
1.  **The Soul (Vision):** The high-level human purpose of the system.
2.  **The Skeleton (Specbook):** The mathematical and structural boundaries.
3.  **The Genome (Actual State):** The physical manifestation in Go code.

By prioritizing **Intent** over **Syntax**, SAAYN-Agent v6 ensures that the codebase is not "written" but "sequenced." If the intent remains eternal, the code can be re-materialized or repaired at any time to match it.

### **1.2. The Deterministic Guarantee: Why "Black Box Trust" Replaces "Vibe Coding"**
Current AI-assisted development suffers from **"Vibe Coding"**—the reliance on an LLM's statistical probability to produce working code. SAAYN-Agent eliminates this "black box" through **Deterministic Physics**.

We do not trust the AI's "vibe." We trust the **Go Toolchain**. 
The **Black Box Trust** model mandates that before logic is even hydrated, a behavioral contract (The Unit Test) must be anchored. The AI is then trapped in a remediation loop where the only exit is a `PASS` signal from the local compiler and test runner. 

### **1.3. The Identity Triad: Defining the Genomic Fingerprint**
To maintain 100% precision during high-fidelity code surgery, SAAYN-Agent v6 discards line-number-based editing. Instead, it identifies every code node through an indestructible **Identity Triad**:

1.  **PublicID:** A globally unique, human-readable identifier (e.g., `pkg.internal.scanner.ScanFile`). This acts as the "address" of the gene.
2.  **Fingerprint:** The structural signature of the node (parameters and return types). This enforces the **Signature Lock**, preventing the AI from changing the contract without authorization.
3.  **Logic Hash:** A cryptographic signature of the AST body, excluding whitespace and comments. This is the **Drift Detector**; if the hash changes, the system knows the logic has mutated.

### **1.4. The MCP Sovereign**
Under this manifesto, SAAYN-Agent is not just a tool; it is a **Sovereign Environment**. By using the **Model Context Protocol (MCP)**, the agent exposes its resources (Specs, Vision, Genome) as standardized protocols. This allows any external "Brain" (LLM) to interface with the local "Body" (Filesystem) through a secure, audited, and stateful umbilical cord.

> **"We are not building a generator; we are building a specialized surgeon that respects the laws of physics and the sanctity of intent."**

## **2. Architectural Anatomy (The Infrastructure)**

Section 2 defines the physical and digital structures required to host the Code Genome. It establishes the "Standard Model" for project organization and the protocols for internal communication and state persistence.

## 💾 The Genome Schema

```json
{
  "project_metadata": {
    "project_name": "SAAYN-Agent",
    "version": "1.0.0",
    "logic_engine": "saayn-v5-ast",
    "total_nodes": 175,
    "last_sequence_hash": "sha256:7f8e9a..."
  },
  "nodes": {
    "pkg.internal.scanner.ScanFile": {
      "uuid": "node-scn-001",
      "public_id": "scanner.ScanFile[scanner.go]",
      "uip": "internal/scanner.ScanFile",
      
      // --- GENESIS STATE MACHINE ---
      // genesis_state: 1 | 2 | 3 | 4
      "genesis_state": 4, 
      
      // maturity: "conceptual" | "hollow" | "anchored" | "sequenced"
      "maturity": "sequenced", 
      
      // --- THE IDENTITY TRIAD ---
      "fingerprint": "func(filePath string) ([]model.Node, error)",
      "logic_hash": "sha256:a1b2c3d4e5f6...", // Set only when state == 4
      "linked_spec_uuid": "spec-pkg-scn-001",

      // --- GENOMIC CONTENT ---
      "gene": "Recursively parses Go files into AST nodes. Excludes whitespace/comments from hash. Must handle symlinks.",
      "business_purpose": "The primary sensory organ for identifying code identity and logic drift.",

      // --- DEPENDENCY GRAPH ---
      // NOTE: Evolution of this field:
      // STATE 1-2: Predicted/Spec-bound (Used to build the Roadmap/DAG).
      // STATE 3:   Dynamic (Updated via JIT Canvas Mounting if Surgeon adds calls).
      // STATE 4:   AST-Derived (The Final Truth). Must be extracted via go/ast 
      //            to verify it matchesauthorized boundaries.
      "dependencies": [
        "pkg.internal.model.Node",
        "pkg.internal.astutil.CalculateHash"
      ],
      
      // --- RECONCILIATION DATA ---
      "last_audit": {
        "timestamp": "2026-04-15T07:34:52Z",
        "physics_pass": true,
        "behavioral_pass": true,
        "cognitive_score": 0.98,
        "last_remediation_count": 1
      }
    }
  }
}
```

### **2.1. The Standard Model: File Hierarchy**
SAAYN-Agent v6 enforces a strict **Standard Model** for project structure. This hierarchy ensures that the Genesis Engine can predictably locate nodes and that the "Sensory Organs" (Scanner) can map the filesystem with 100% accuracy.



* **`cmd/`**: The Command Tier. Contains the entry points for the CLI. Logic here is strictly limited to bootstrap and orchestration calls.
* **`internal/`**: The Core Anatomy. Private logic, state management, and the Surgical Inner Loop. No code outside this project may import from this directory.
* **`pkg/`**: The Public DNA. Shared models and library code that define the "Connectome" between packages.
* **Root Artifacts**:
    * `specbook.yaml`: The Genotype (Desired State).
    * `genome.json`: The Phenotype (Actual State).
    * `genome.index.json`: The Nervous System (Semantic State).

### **2.2. The DNA Registry: `genome.json` Schema**
The DNA Registry is the definitive record of the **Actual State**. It tracks the metamorphosis of every node from a conceptual requirement to a sequenced reality.

**Key Registry Logic:**
* **Stateful Tracking**: Every node must carry a `genesis_state` (1–4) and a matching `maturity` label.
* **The Gene Field**: A node-specific, immutable behavioral instruction that acts as the prompt for State 3 (Hydration).
* **The Identity Triad**: Mandatory inclusion of `public_id`, `fingerprint`, and `logic_hash` for every sequenced node.

### **2.3. The Nervous System: `genome.index.json`**
The Nervous System provides the **Semantic Map** of the genome. While `genome.json` understands the physics of the code, `genome.index.json` understands the *intent*.

* **Vector Embeddings**: Uses 1536-dimensional vectors to represent the "Business Purpose" and "Gene" of every node.
* **Proximity Logic**: Enables natural language discovery (e.g., "Find where we handle Stripe callbacks") by calculating cosine similarity between human intent and the code's semantic signature.
* **Local Storage**: To ensure privacy and speed, the index is stored and queried locally, requiring an LLM only for the initial vector generation.

### **2.4. MCP Integration: The SAAYN-MCP Server**
The **Model Context Protocol (MCP)** acts as the universal docking port. By wrapping the infrastructure in an MCP server, we decouple the "Cognitive Brain" from the "Local Body."

**MCP Resource Mapping:**
* `saayn://vision`: Streams the `readme.md` (Soul).
* `saayn://spec/nodes/{id}`: Streams the YAML contract (Genotype).
* `saayn://genome/nodes/{id}`: Streams the current AST state (Phenotype).

**MCP Tool Mapping:**
* `get_node_dna`: Allows the agent to "Lazy Load" a specific AST node.
* `mount_canvas`: Triggers the FAST model to generate State 2 stubs.
* `apply_surgery`: Submits a patch for State 3 hydration, gated by the Physics Audit.

### **2.5. Dependency Resolution: The JIT Orchestrator**
The architecture includes a specialized **Just-In-Time (JIT) Orchestrator** that prevents "Dependency Hallucination."

If a node in State 3 (Hydration) attempts to call an external node that is currently in State 1 (Conceptual), the Orchestrator **freezes** the current sequence, initiates a **State 2 (Hollow)** materialization for the dependency, and only then resumes the original surgery. This ensures the compiler always sees a valid interface.


## **3. The Metamorphosis Pipeline (The 5-State Machine)**

Chapter 3 defines the operational lifecycle of a Genomic Node. This pipeline is a strict, one-way state machine. A node cannot bypass a state, ensuring that the **"Physics"** of the system are always locked before the **"Logic"** is allowed to evolve.

### **3.1. State 1: Conceptual (The Gallery)**
* **Agent:** DEEP Model (Architect)
* **Focus:** Semantic Boundary & Business Intent
* **Action:** The node is defined in the `specbook.yaml`. The agent debates and finalizes the `business_purpose` and the `gene`. 
* **Guardrail:** **No Code Allowed.** The agent is physically blocked from generating Go syntax. This phase is purely for establishing the "Why" and the "Contract."

### **3.2. State 2: Hollow (The Canvas)**
* **Agent:** FAST Model (Structural Engineer)
* **Focus:** Structural Topography
* **Action:** The agent reads the Specbook and materializes a "Hollow Stub" on the local disk. It writes the package declaration, imports, empty structs, and zero-return functions (`return nil`).
* **Guardrail:** **Zero-Logic Rule.** Any detection of `if`, `for`, or variable assignments triggers an automatic rejection. 
* **The Physics Gate:** The node must pass an AST Physics Audit (`go/ast`) to ensure it is a valid, buildable Go skeleton.

### **3.3. State 3: Anchored (The Contract)**
* **Agent:** Local CC Agent (The Judge)
* **Focus:** **[Black Box Trust]** Behavioral Baseline
* **Action:** The engine generates a `_test.go` file containing table-driven test cases based on the `input/output` contracts defined in State 1.
* **Guardrail:** **The Behavioral Vacuum.** The loop executes `go test`. The test **must fail** with a "not implemented" or "zero-value mismatch" error. This failure anchors the reality: we now have a mathematical proof of the requirement.

### **3.4. State 4: Hydrating (The Surgery)**
* **Agent:** DEEP Model (Surgeon)
* **Focus:** Algorithmic Implementation
* **Action:** The Surgeon is provided with the Hollow Signature, the Gene, and the failing Test Case. It writes the actual execution logic.
* **The JIT Constraint:** If the Surgeon attempts to call a dependency that is still in State 1, the **JIT Canvas Mounting** protocol pauses the loop, forces the dependency into State 2, and then resumes.
* **The Signature Lock:** The Surgeon is forbidden from changing the function signature. If a signature change is required, the node must be demoted back to State 1 for a "Re-stretch."

### **3.5. State 5: Sequenced (Equilibrium)**
* **Agent:** Local CC Agent (The Record-Keeper)
* **Focus:** Genomic Locking
* **Action:** Once the code passes the **Behavioral Audit** (tests pass) and the **Cognitive Audit** (logic matches the Gene), the agent performs the final sequence.
* **The Identity Triad:** The system calculates the whitespace-agnostic `logic_hash` and locks the `public_id`.
* **Persistence:** The node is marked as `State 5` in `genome.json`. It is now a permanent part of the project's "Actual State."

### **3.6. Summary of State Transitions**

| State | Maturity | Primary Gate | Artifact Generated |
| :--- | :--- | :--- | :--- |
| **1** | Conceptual | Spec Validation | Specbook Entry + Gene |
| **2** | Hollow | AST Physics Audit | `.go` Stub (Skeleton) |
| **3** | Anchored | Test Failure Gate | `_test.go` (The Anchor) |
| **4** | Hydrating | Test Success Gate | Functional Go Logic |
| **5** | Sequenced | Hash Lock Gate | Final `genome.json` Entry |

## **4. The Surgical Inner Loop (The Immune System)**

Chapter 4 defines the active verification engine that powers the Genesis transitions. This "Immune System" is a set of autonomous audits that protect the project from structural decay, behavioral drift, and "hallucinated" logic. If a node cannot pass through these four gates, it is physically impossible for it to be committed to the disk.

### **4.1. Gate 1: The Physics Audit (Structural Integrity)**
* **Tooling:** `go/ast` and `dave/dst`
* **Purpose:** To ensure the code adheres to the laws of the Go language and the project's internal anatomy.
* **Execution:** Before any logic is reviewed, the local SAAYN agent performs a **Structural Walk**. It verifies:
    * **Syntax Lock:** Does the code parse? (No missing brackets or illegal characters).
    * **Interface Compliance:** If the node is a struct, does it satisfy the interfaces defined in the **Hollow Canvas**?
    * **Import Hygiene:** Are all imports utilized and authorized by the Specbook?
* **Failure Protocol:** If the Physics Audit fails, the draft is instantly rejected. The agent sends the raw compiler error back to the Surgeon for immediate re-materialization.

### **4.2. Gate 2: The Behavioral Audit (The Vacuum Gate)**
* **Tooling:** `go test -v`
* **Purpose:** To transform the **Black Box Trust** into a physical reality.
* **Execution:** This gate operates in two distinct modes:
    1.  **Phase 3 (Anchor):** The audit *must fail*. We confirm that the test runner successfully identifies the missing logic. This proves the test is "live."
    2.  **Phase 4 (Hydration):** The audit *must pass*. The Surgeon's logic is executed against the table-driven test cases. 
* **Failure Protocol:** If `go test` returns a non-zero exit code during hydration, the **stdout/stderr** logs are captured and fed into the Remediation Cycle as high-signal feedback.

### **4.3. Gate 3: The Cognitive Audit (Intent Alignment)**
* **Tooling:** Vector Similarity (`genome.index.json`) + LLM Evaluator
* **Purpose:** To ensure the code does what the **Gene** says it should do, not just what the compiler allows.
* **Execution:** The system compares the materialized logic against the **Gene** (State 1 intent) and the **Vision** (README).
    * **Drift Detection:** "The Gene requires exponential backoff, but the implementation uses a linear sleep. **REJECTED.**"
    * **Semantic Scoring:** The node is assigned a "Cognitive Score." A score below 0.85 triggers an automatic remediation, even if the tests pass.
* **Failure Protocol:** The Evaluator provides a "Finding" (e.g., "Missing context-aware cancellation") which is used to refine the next iteration of the code.

### **4.4. Gate 4: The Signature Lock (The Contract Protector)**
* **Tooling:** Fingerprint Matching
* **Purpose:** To prevent "Recursive Collateral Damage."
* **Execution:** The engine compares the **Fingerprint** of the new patch against the Fingerprint locked in the **Hollow State**. 
    * **The Law:** The Surgeon is allowed to change *how* the function works, but never *what* the function looks like to the rest of the world.
* **Failure Protocol:** If a Signature Mismatch is detected, the CC Agent blocks the splice. It forces a **Canvas Re-stretch**, requiring the architect to acknowledge that changing this signature will affect every dependent node in the DAG.

### **4.5. The Remediation Cycle (Self-Healing)**
If any gate fails, the node enters the **Remediation Cycle**. 
1.  **Feedback Synthesis:** The CC Agent aggregates Physics errors, Test failures, and Cognitive findings.
2.  **Context Injection:** The Surgeon is handed its previous (failed) attempt and the aggregate feedback.
3.  **Iteration Cap:** The loop allows for a maximum of **3 iterations**. If Equilibrium is not reached by the third attempt, the system "Freezes" the node and alerts the Human Architect for manual intervention.

## **5. Tooling & Sensory Organs**

Chapter 5 details the specialized internal mechanisms that allow SAAYN-Agent v6 to interact with physical source code. These are the "hands" and "eyes" of the Genesis Engine, moving beyond simple text manipulation into the realm of **AST-native orchestration**.

### **5.1. The Scanner: The Sensory Organ**
* **Primary Tool:** `go/ast` & `go/parser`
* **Responsibility:** The Scanner is responsible for **Phenotype Extraction**. It reads the local filesystem and translates raw `.go` files into the structured data required for the `genome.json`.
* **The Identity Extraction:**
    * **PublicID Retrieval:** It walks the AST to find function declarations and receiver types to build the `Package.Receiver.Func` string.
    * **Fingerprint Calculation:** It normalizes parameter names to types only (e.g., `func(string) error`) to create a collision-proof structural signature.
    * **Logic Hashing:** It extracts the `*ast.BlockStmt` (the function body), strips whitespace and comments, and generates a SHA-256 hash. This allows SAAYN to detect a "Logic Mutation" even if a human only changed the indentation.

### **5.2. The Surgeon: The Splicing Mechanism**
* **Primary Tool:** `dave/dst` (Decorated Syntax Tree)
* **Responsibility:** The Surgeon performs the **Hydration Surgery**. While standard AST is "lossy" (it discards comments), DST preserves the "Soul" of the code (human-written documentation and intent markers).
* **The Splicing Protocol:**
    1.  **Targeting:** The Surgeon uses the `PublicID` to find the exact coordinates of the Hollow stub within the DST.
    2.  **Validation:** It performs a pre-surgical check to ensure the existing `Fingerprint` matches the `genome.json`.
    3.  **Injection:** It replaces the empty `*dst.BlockStmt` of the stub with the hydrated logic generated by the DEEP model.
    4.  **Preservation:** It re-anchors any "hanging" comments from the original stub to the new implementation logic, ensuring the "Intent" remains attached to the "Action."

### **5.3. The JIT Orchestrator: The Central Nervous System**
* **Responsibility:** Managing the **Build Order** and recursive dependencies.
* **The Freeze-and-Mount Logic:**
    * During State 4 (Hydration), if the Surgeon encounters a call to an undefined node, the Orchestrator interrupts the process.
    * It places the current node into a **Wait State**.
    * It triggers the **FAST Model** to "Mount" a State 2 Hollow Canvas for the missing dependency.
    * Once the dependency passes the Physics Audit, the Orchestrator signals the Surgeon to resume. This ensures the Surgeon never writes code against an interface that hasn't been "physically" verified.

### **5.4. The Identity Triad Verification Logic**

The following Go-pseudo-logic represents how the Scanner verifies the **Identity Triad** during a Genomic Audit:

```go
func (s *Scanner) VerifyNode(node *dst.FuncDecl) (bool, error) {
    currentFingerprint := s.ExtractFingerprint(node)
    currentHash := s.CalculateLogicHash(node.Body)

    // Gate 1: Signature Lock
    if currentFingerprint != genome.Fingerprint {
        return false, ErrSignatureMismatch // Triggers State 2 Re-stretch
    }

    // Gate 2: Drift Detection
    if currentHash != genome.LogicHash {
        return false, ErrLogicDrift // Triggers State 3 Re-hydration
    }

    return true, nil
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
