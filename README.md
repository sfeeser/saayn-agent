# SAAYN: Code Genome System (CGS)
> **"The AI proposes changes. CGS enforces correctness."**

SAAYN (CGS Version) is a deterministic, AST-based mutation and repair engine. It treats your Go codebase as a living genome, providing a cryptographic ledger (the Registry) that allows AI agents to perform surgical edits without breaking the build.

## 🧬 Core Philosophy
Standard AI coding tools treat code as **Text**. SAAYN treats code as **Structure**.
* **No Markers:** No more `// CHUNK_START`. We use AST fingerprints.
* **No Drift:** Logic hashes ignore whitespace and comments.
* **No Breakage:** Every edit is validated by `go build` before it touches your disk.

---

## 🚀 Quick Start

### 1. Installation
Clone the repo and build the binary:
```bash
git clone https://github.com/your-org/saayn-agent
cd saayn-agent
go build -o saayn main.go
```

### 2. The Digital Census (Init)
Point SAAYN at your Go project to index the DNA. This creates `genome.json`.
```bash
./saayn init --path /path/to/your/project
```

### 3. The Surgical Edit (Test)
To test a mutation, you can manually trigger the Surgeon. (Note: In production, this is handled by the AI Agent).
```bash
# This targets a specific UUID in your genome.json 
# and replaces its body with new logic.
./saayn edit --uuid "abc-123" --body "return a + b + 10"
```

---

## 🏗 System Architecture

The system operates in a 4-stage loop:

1.  **SCAN**: `internal/scanner` walks the AST to find every function (The Twigs).
2.  **IDENTIFY**: `internal/genome` creates a stable **Identity Triad** (Public ID, Fingerprint, Logic Hash).
3.  **GRAFT**: `internal/surgeon` performs a subtree swap in the AST memory model.
4.  **GATE**: `internal/validator` runs `go build`. If the code is invalid, the surgery is rolled back.



---

## 🛠 Directory Map
* `/cmd/saayn`: The CLI entry point.
* `/internal/scanner`: The "Eyes" (AST Parsing).
* `/internal/genome`: The "Brain" (Identity and Hashing).
* `/internal/surgeon`: The "Hands" (AST Mutation).
* `/internal/validator`: The "Immune System" (Compiler Gates).
* `/pkg/model`: The "DNA" (Schema).

---

## ⚖️ Guarantees
* **Indestructible Identity:** Renaming a file does not break the registry.
* **Atomic Commits:** Files are only updated if they pass `go vet` and `go build`.
* **Format Preservation:** All edits are automatically passed through `gofmt`.
