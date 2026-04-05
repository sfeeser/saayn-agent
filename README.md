# SAAYN-AGENT: The Code Genome System
> **"Specifications Are All You Need."**

## 🏗 System Architecture

SAAYN V5 is a **Semantic Surgery Engine** designed to treat your codebase as a living organism rather than a collection of text files.

### 🧠 The Intent Bridge
The core of the architecture is an agentic **Local Vector Index** that maps human language to code reality. When you ask a question in "plain English," SAAYN performs a high-speed mathematical lookup against your locally stored embeddings.

* **Zero-Knowledge Discovery:** You don't need to know the function names. Searching for *"How do we handle expired sessions?"* will rank the relevant security logic #1, even if the function is named `pkg.ValidateTick()`.
* **Privacy-First Embeddings:** Your semantic map is processed and stored **locally** in `genome.index.json`. Your architectural "intent" stays on your machine, not in a cloud database.
* **The Identity Triad:** Every piece of code is tracked by its **Public ID**, **Structural Fingerprint**, and **Logic Hash**. If you change a single character, the "Nervous System" detects the drift instantly.

### 🛠 The Workflow Loop
1.  **Ingest:** `init` and `enrich` build the initial "Map of Intent."
2.  **Locate:** `search-intent` and `trace` find the exact "Twigs" that need pruning.
3.  **Operate:** `plan` and `apply` execute multi-file AST splicing with 100% syntax safety.


---

## 🚀 Quick Start

### Installation
Clone the repo and build the binary:
```bash
git clone https://github.com/your-org/saayn-agent
cd saayn-agent
go build -o saayn main.go
```

### 🚀 The SAAYN Workflow (In Execution Order)

### 🧬 Initialize the Code Genome
* **`saayn verify-llm-targets`** – Validate connectivity for FAST, DEEP, and embedding tiers.
* **`saayn init`** – Map all functions, structs, and variables into a `genome.json` registry.
* **`saayn enrich`** – Generate semantic "business purpose" summaries and logic hashes.

### 🔍 Verify Your Context
* **`saayn search-intent`** – Perform semantic natural language searches across the genome.
* **`saayn trace`** – Structurally identify call stacks and dependency trees for any target.
* **`saayn verify`** – Audit the codebase for "logic drift" against registered genome hashes.

### ⚡ Execute the Surgery
* **`saayn draft`** – Translate natural language requests into structured multi-file plans.
* **`saayn graph`** – Analyze the "blast radius" and hydrate source code for deep context.
* **`saayn plan`** – Orchestrate LLMs to generate and refine strict code patch sets.
* **`saayn apply`** – Execute batch AST splicing to write patches to the local filesystem.
* **`saayn gen-test`** – Automatically generate and run Go tests to verify mutated logic.

---

### 🛠️ Maintenance & Utilities
* **`saayn help`** – Provides detailed documentation and flags for any specific command.
* **`saayn completion`** – Generates shell autocompletion scripts for a smoother CLI experience.

## ⚖️ Guarantees
* **Indestructible Identity:** Renaming a file does not break the registry.
* **Atomic Commits:** Files are only updated if they pass `go vet` and `go build`.
* **Format Preservation:** All edits are automatically passed through `gofmt`.

### cobra style help, just type "saayn"
```
./saayn 
A deterministic, AST-based mutation and repair engine 
that treats code as a living genome.

Usage:
  saayn [command]

Available Commands:
  apply              Apply a generated V5 surgery patch set to the local filesystem using batch AST splicing
  completion         Generate the autocompletion script for the specified shell
  draft              Initialize a code surgery plan based on a natural language intent
  enrich             Uses FAST Cognition to automatically document the business purpose of code
  gen-test           Generates and verifies a Go test for the last mutated function
  graph              Analyze the blast radius of a planned surgery and hydrate context source code
  help               Help about any command
  init               Initialize the Code Genome for the project
  plan               Generate and auto-refine a strict multi-file code patch set from a surgery plan
  search-intent      Search the semantic genome index using a natural language intent
  trace              Structurally grep the codebase to find all callers of a specific target
  verify             Detects logic drift between live code and the genome.json
  verify-llm-targets Verifies end-to-end LLM connectivity for both FAST and DEEP cognitive tiers
```

### 1. initialize the code genome.json file with the ID of major nodes

```text
./saayn init                                                                                                                                                                       
🧬 Initializing Code Genome at: .                                                                                                                                                                                   
✅ Success! Indexed 175 nodes into genome.json                                                                                                                                                                      
💡 Next step: Run './saayn enrich' to generate semantic summaries.
```

### 2. Now populate the genome json file with the business purpose of each node
saayn-agent contacts your AI and sequences your code's genome
```
./saayn enrich
🧠 Starting Semantic Enrichment Process...
⚙️   saayn.printActiveNodeHeader[enrich.go]        cmd/saayn/enrich.go       [21939968]
 ├─ 🧬 logic changed
 ├─ 🔄 purpose reset
 ├─ 🔍 analyzing [175/175]
 └─ ✅ purpose updated
    This function serves to display a standardized header for an "active" SAAYN
    code node, providing immediate visual context about its type, identity,
    file location, and a truncated logic hash. It helps the autonomous agent or
    an observer track and understand which specific code element SAAYN is
    currently processing or analyzing within the code genome.
📡 Synchronizing Semantic Index...
📊 Enrichment Summary:
  - Updated: 1
  - Skipped (Already Enriched): 174
  - Skipped (Missing/Drifted):  0
  - Failed:  0
🧠 Semantic Index Summary:
  - Created: 0
  - Updated: 1
  - Deleted: 0
  - Skipped: 174

💾 Genome memory successfully updated.
💾 Semantic index successfully updated.
```

### 3. Now use saayn to walk to the code genome:
The response you see below will generate in less than a second.  
saayn-agent does inferencing locally and it is blazing fast
```text
./saayn search-intent "find the code that formats the summary text with indentation"

🔎 Search Intent: find the code that formats the summary text with indentation

Top Matches
1. saayn.printWrappedIndented               cmd/saayn/search_intent.go     function   0.71
2. saayn.wrapText                           cmd/saayn/enrich.go            function   0.68
3. saayn.printIndexSummary                  cmd/saayn/enrich.go            function   0.66
4. saayn.printEnrichmentSummary             cmd/saayn/enrich.go            function   0.65
5. saayn.printSearchResults                 cmd/saayn/search_intent.go     function   0.63
```
