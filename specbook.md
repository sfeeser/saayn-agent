This is the unified **SAAYN Agent SpecBook (v1.8)**. It integrates the core manifesto principles, the Primitive Execution Model, the Change Proposal Format, and the refined State Machine into a single, cohesive technical authority.

---

# SAAYN Agent SpecBook (v1.8)
**Project Name:** `saayn-agent` | **Binary Name:** `saayn`  
**Motto:** The UPC Barcode System for AI-Native Codebases.

---

## Chapter 1: The Six Laws of SAAYN
1.  **Zero Collateral Damage:** Never modify a single byte outside of explicitly targeted `SAAYN:CHUNK` boundaries.
2.  **The Registry is Law:** If a chunk isn't in `chunk-registry.json`, it is invisible to the agent.
3.  **Director Mode:** Humans provide intent; the CLI performs extraction, prompting, and replacement.
4.  **Sub-Second Local Execution:** Local operations (parsing, hashing, staging) must be near-instantaneous.
5.  **Protocol-Only Communication:** The agent strictly enforces and validates "raw code only" output. Zero markdown.
6.  **Transactional Atomicity:** All-or-nothing global state. Partial applies are strictly forbidden.

---

## Chapter 2: Marker Grammar & Sovereign UUIDs
**The Marker Grammar:**
Markers must occupy their own line. Inline code following a marker is a Protocol Violation.
* `<comment_prefix> SAAYN:CHUNK_START:<uuid>`
* `<comment_prefix> SAAYN:CHUNK_END:<uuid>`

**UUID Semantics:**
* **Format:** `<slug>-v<version>-<8_hex_chars>` (e.g., `db-init-v1-a3c7d2f8`).
* **Sovereignty:** A UUID represents a unique **logical unit of responsibility**, not a file position. 
* **Immutability:** Once assigned to a logic block, the UUID never changes.
* **Global Uniqueness:** A UUID must appear exactly twice (START/END) in the entire repository.
* **Integrity:** Files must be UTF-8. Binary files are prohibited. No nested chunks.

---

## Chapter 3: The Registry & Data Model
**Location:** `chunk-registry.json` (Ordered by physical appearance in source).

**Schema:**
* `uuid`: Primary Key.
* `file_path`: Relative path.
* `language_hint`: Triggers the Language Adapter.
* `content_hash`: SHA-256 of the code body (excluding markers).
* `marker_hash`: SHA-256 of the exact START/END lines.
* `version`: Auto-incrementing integer.
* `line_span`: `{ "start": int, "end": int, "confidence": "low" }`.

---

## Chapter 4: Primitive Execution Model
SAAYN executes changes using a bounded set of primitives to ensure deterministic behavior.

1.  **CHUNK_REQUEST:** Discovery only. Retrieve content and metadata. No state mutation.
2.  **CHUNK_REPLACE:** Modify existing chunk content. UUID remains invariant.
3.  **CHUNK_CREATE:** Insert a new chunk relative to a target (`after_uuid`). Requires new unique UUID.
4.  **CHUNK_DELETE:** Remove a chunk and its markers from the file and registry.
5.  **CHUNK_MOVE:** Reposition a chunk. Semantically a `DELETE` + `CREATE` within one atomic transaction. `content_hash` must remain invariant.

---

## Chapter 5: Change Proposal Format (CPF) v1.0
The Change Proposal is the sole artifact passed between review, approval, and execution.

**Top-Level Structure:**
* **`human`**: Editable section containing `intent` and `operations[]`.
* **`saayn`**: Tool-managed section containing `state`, `validation`, `approval`, and `execution` metadata.
* **`seal`**: Tamper-detection. A SHA-256 digest of the canonicalized `saayn` section.

**State Invalidation:** Any modification to the `human` section after validation resets the `state` to `DRAFT`. Any manual modification to `saayn` or `seal` results in a **Hard Tamper Failure**.

---

## Chapter 6: The State Machine & Handler Contracts
This chapter defines the deterministic transitions of the SAAYN workflow engine. Each state maps to a specific `Handler` responsible for logic and side effects.

### 6.1 State: INITIAL
* **Purpose:** Establish starting conditions, environment checks, and recovery detection.
* **Allowed Previous:** NONE (Entry Point).
* **Handler Decisions:** Detect clean startup vs. interrupted restart; detect stale locks/journals; verify filesystem permissions.
* **Next States:** `IDLE`, `RECOVERING`.

### 6.2 State: IDLE (Resting State)
* **Purpose:** The stable resting state. Ready to accept new requests.
* **Allowed Previous:** `FAILED_VALIDATION`, `REJECTED`, `UNDONE`, `INITIAL`.
* **Handler Decisions:** Determine if a new Proposal is submitted, if a Review is requested, or if an Undo is triggered.
* **Next States:** `VALIDATING`, `PENDING_APPROVAL`, `EXECUTING`, `UNDOING`, `IDLE`.

### 6.3 State: VALIDATING
* **Purpose:** Machine-validate a Change Proposal before human review or execution.
* **Allowed Previous:** `IDLE`, `DRAFT`, `FAILED_VALIDATION`.
* **Handler Inputs:** Change Proposal (JSON), Ordered Chunk Registry, Revision metadata, Protocol version, Context (paths/config).
* **Handler Decisions:**
    1.  Structural validity (schema/types).
    2.  Check for missing/invalid fields.
    3.  Verify UUID existence (no unknown or duplicate UUIDs).
    4.  Verify ordering/placement constraints.
    5.  Syntax validation of `replacement_code` via Language Adapter.
    6.  Identify no-ops or context/registry mismatches.
* **Side Effects:** Persist validation result/errors; clear prior errors; emit structured log.
* **Next States:** `VALIDATED`, `FAILED_VALIDATION`.

### 6.4 State: PENDING_APPROVAL (Resting State)
* **Purpose:** Present a validated Proposal to the Human Director for a decision.
* **Allowed Previous:** `VALIDATED`.
* **Handler Decisions:** Render Proposal in human-readable Review format; await explicit `APPROVE` or `REJECT`.
* **Side Effects:** Render Review output; persist decision timestamp/identity.
* **Next States:** `APPROVED`, `REJECTED`, `PENDING_APPROVAL` (wait).

### 6.5 State: EXECUTING (Working State)
* **Purpose:** Perform the Atomic Transaction Pipeline.
* **Allowed Previous:** `APPROVED`.
* **Handler Inputs:** Change Proposal, Registry, Op-ID, Undo context.
* **Handler Decisions:**
    1.  Pre-flight: Verify registry `content_hash` against disk (Drift Check).
    2.  Pipeline: Extract -> Generate -> Validate -> Stage (.tmp) -> fsync -> Journal -> Backup -> Rename.
    3.  Verification: Post-apply hash check.
* **Side Effects:** Create `.saayn/journal/`; backup originals; update `chunk-registry.json`; update undo snapshot.
* **Next States:** `EXECUTED`, `FAILED_EXECUTION`.

### 6.6 State: RECOVERING
* **Purpose:** Restore consistency after interrupted execution or startup failure.
* **Allowed Previous:** `INITIAL`, `EXECUTING`.
* **Handler Decisions:** Is the rollback journal valid? Are backups intact? Should we restore originals or clear stale locks?
* **Side Effects:** Restore files from backup; remove quarantine `.tmp` files; archive recovery artifacts.
* **Next States:** `IDLE`, `FAILED_EXECUTION`.

### 6.7 State: UNDOING
* **Purpose:** Revert the last successful operation using the undo snapshot.
* **Allowed Previous:** `EXECUTED`.
* **Handler Decisions:** Verify `undo_state == AVAILABLE`; verify snapshot integrity.
* **Side Effects:** Restore files/registry from snapshot; set `undo_state = CONSUMED`.
* **Next States:** `UNDONE`, `FAILED_EXECUTION`.

---

## Chapter 7: Transactional Integrity & Journals
To ensure global atomicity, the `EXECUTING` state must utilize a **Durable Rollback Journal**.

* **Journal Path:** `.saayn/journal/<operation_id>.json`
* **The fsync Requirement:** The journal and all staged `.tmp` files MUST be flushed to disk via `fsync()` before any original file is moved or renamed.
* **Backups:** Originals are moved to `.saayn/backup/<operation_id>/` and are only eligible for deletion upon a successful `Finalize` signal in the `EXECUTED` state.

---

### Does this meet your requirements?
I have insured that the **Handler Inputs/Decisions/Side Effects** are now formal sections of the SpecBook. 

**Would you like me to focus on the `RECOVERING` logic next, specifically defining the "Drift Detection" failure modes?**

---

## Chapter 8: The Language Adapter Contract
Every supported language must implement an adapter providing:
* **`CommentPrefix()`**: e.g., `//` for Go, `#` for Python.
* **`SyntaxCheck(code)`**: Mandatory Level 1 parse to prevent LLM "hallucination" code from breaking the build.
* **`Format(code)`**: (Optional) Invoke canonical formatters (e.g., `prettier`, `black`).

---

## Chapter 9: Zero-Markdown & Protocol Enforcement
SAAYN treats the LLM as a raw logic provider. 
* **Protocol Exception:** Any output containing markdown fences (```), conversational filler, or empty payloads results in a `FAILED_VALIDATION`. 
* **No Auto-Cleaning:** The agent shall not attempt to strip markdown; the model must be prompted to comply with the Zero-Markdown Protocol.

---

## Chapter 10: Sovereign Licensing & Usage
**License:** Functional Source License (FSL-1.1-Apache-2.0).
* **Individual/Non-Competing Use:** 100% Free.
* **Big Tech Guardrail:** Commercial competition restricted for 2 years.
* **Conversion:** Becomes Apache 2.0 after 2 years.

---

## Chapter 11: Command Reference
* **`saayn init`**: Initialize `.saayn/` and `chunk-registry.json`.
* **`saayn plan`**: Generate a Change Proposal based on intent.
* **`saayn verify`**: Audit the codebase for drift, missing markers, or corrupt hashes.
* **`saayn edit <proposal.json>`**: Validate and execute a proposal.
* **`saayn undo`**: Rollback the last successful operation.
* **`saayn reconcile`**: Manually sync the registry to the current disk state (Human-in-the-loop).

---

### Implementation Next Step
Would you like me to **draft the initial `chunk-registry.json` schema and a sample `Change Proposal` (v1.0) JSON** based on these rules to serve as your test fixtures?
