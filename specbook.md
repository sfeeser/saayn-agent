# SAAYN Agent SpecBook (v1.7)
**Project Name:** `saayn-agent` | **Binary Name:** `saayn`  
**Motto:** The UPC Barcode System for AI-Native Codebases.

## Chapter 0: The Six Laws of SAAYN
1.  **Zero Collateral Damage:** Never modify a single byte outside of explicitly targeted `SAAYN:CHUNK` boundaries.
2.  **The Registry is Law:** If a chunk isn't in `chunk-registry.json`, it is invisible to the agent.
3.  **Director Mode:** Humans provide intent; the CLI performs extraction, prompting, and replacement.
4.  **Sub-Second Local Execution:** Local operations (parsing, hashing, staging) must be near-instantaneous.
5.  **Zero-Markdown Protocol:** The agent strictly enforces and validates "raw code only" output.
6.  **Transactional Atomicity:** All-or-nothing global state. Partial applies are strictly forbidden.

## Chapter 1: Marker Grammar & UUID Semantics
**The Marker Grammar:**
* `<comment_prefix> SAAYN:CHUNK_START:<uuid>`
* `<comment_prefix> SAAYN:CHUNK_END:<uuid>`

**UUID Format & Generation:**
* **Format:** `<slug>-v<version>-<8_hex_chars>` (e.g., `db-init-v1-a3c7d2f8`).
* **Generation:** Slug is sanitized from intent; suffix is cryptographically random.
* **Continuity Rule:** A UUID represents a unique **logical unit of responsibility**, not a file position. UUIDs must never be "recycled" for unrelated logic.

**Invariants:**
* **Immutability:** Once assigned to a logic block, the UUID never changes.
* **Global Uniqueness:** A UUID must appear exactly twice (START/END) in the entire repository.
* **Single Mapping:** A UUID maps to exactly one `file_path`.
* **Integrity:** Files must be UTF-8. Binary files are prohibited. No nested chunks.

## Chapter 2: The Registry Data Model
**Location:** `chunk-registry.json`

**Schema:**
* `chunks` (Array of Objects):
    * `uuid`: Primary Key.
    * `file_path`: Relative path.
    * `language_hint`: (e.g., `go`, `html`). Triggers the **Language Adapter**.
    * `business_purpose`: The "Why" (AI Context).
    * `content_hash`: SHA-256 of the code body.
    * `marker_hash`: SHA-256 of the exact START/END lines.
    * `version`: Auto-incrementing integer.
    * `line_span`: `{ "start": int, "end": int, "confidence": "low" }` (Advisory only).

## Chapter 3: The Durability & Transaction Model
To ensure global atomicity, `saayn` uses a **Durable Rollback Journal** with mandatory `fsync()`.

**Journal Requirements:**
* Located at `.saayn/journal/<operation_id>.json`.
* Contains `(original_path, backup_path, expected_hash)` for every target.
* Must be flushed to disk via `fsync()` before any file move occurs.

**Backups:**
* Originals are moved to `.saayn/backup/<operation_id>/`.
* Backups are only deleted upon a successful **Finalize** signal.

## Chapter 4: The Language Adapter Contract
Every adapter must implement:
* **`CommentPrefix()`**: Syntax for the specific language.
* **`SyntaxCheck(code)`**: **Level 1 (Mandatory)** - Parse-only validation to ensure the LLM returned syntactically valid code.
* **`Format(code)`**: (Optional) Canonical formatting (e.g., `go fmt`).

## Chapter 5: Zero-Markdown & Cache Semantics
**Extraction Rule:**
Valid LLM output must contain **zero** backticks, be non-empty, and pass the **SyntaxCheck**. Any conversational filler or markdown fences result in a **Hard Fail**.

**Cache Key Invariants:**
To ensure idempotency, the cache key must include:
`{ "uuid", "content_hash", "intent_hash", "coder_model", "prompt_version" }`

## Chapter 6: Final Correct Execution Flow

1.  **Acquire Lock:** Create `.saayn.lock`.
2.  **Recovery Check:** Detect existing journals; restore from backups if found.
3.  **Plan:** Identify target UUIDs with human-readable justifications.
4.  **Extract & Verify:** Read original files; verify `content_hash` matches registry.
5.  **Generate:** Coder model produces new code blocks.
6.  **Validate:** Run **SyntaxCheck** + **Format** + **SAAYN_TEST_CMD**.
7.  **Pre-Commit Revalidation:** Re-read original files; abort if `content_hash` changed since Step 4.
8.  **Stage:** Write all `.tmp` files and `registry.tmp`. Call `fsync()` on all.
9.  **Journal:** Write and `fsync()` the `saayn_journal.json`.
10. **Backup:** Move originals to `.saayn/backup/<op_id>/`.
11. **Apply:** Rename `.tmp` files to originals.
12. **Post-Apply Verification:** Re-read applied files; recompute hashes; abort/recover if mismatch.
13. **Commit Registry:** Rename `registry.tmp` to `chunk-registry.json`.
14. **Cleanup:** Delete journal and backups.
15. **Release Lock.**

## Chapter 7: Observability & Verification
* **`saayn verify`**: Returns a JSON object with statuses: `SYNC`, `MODIFIED`, `MISSING`, `DUPLICATE`, or `CORRUPTED`.
* **`saayn reconcile`**: Updates registry only after explicit human confirmation of drift.
* **`saayn undo`**: Reverts state using the `SAAYN_OP:<operation_id>` git tag and registry rollback.

## Chapter 8: Sovereign Licensing
**License:** Functional Source License (FSL-1.1-Apache-2.0).
* **Individual/Non-Competing Use:** 100% Free.
* **Big Tech Guardrail:** Commercial competition restricted for 2 years.
* **Conversion:** Becomes Apache 2.0 after 2 years.

## Chapter 9: Commands
* **`saayn init`**: Setup repo.
* **`saayn plan`**: Preview intended changes.
* **`saayn edit`**: Execute full transaction.
* **`saayn create`**: New chunk creation with explicit placement (`--mode <append|after:uuid|before:uuid>`).
* **`saayn verify`**: Detect drift.
* **`saayn reconcile`**: Update registry after human confirmation of manual edits.
