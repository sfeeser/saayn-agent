# 🧬 The Greenfield Protocol: Spawning Life from Intent
- Code is ephemeral.
- Intent is eternal.

SAAYN-Agent v6 introduces GENESIS, a deterministic "Genesis Engine" that materializes full Go projects from a Markdown Vision and a YAML Specification. It doesn't just "write code"; it sequences a living genome through a recursive Surgical Inner Loop.

### 🚀 The Genesis Command
To bootstrap a new project, you provide the Soul (readme.md), the Skeleton (specbook.yaml), and a target destination.

```Bash
./saayn greenfield -f readme.md -s specbook.yaml --target ./my-new-app
```

### 🔄 The Surgical Inner Loop (The "Immune System")

Most AI agents fail because they "guess" their way through a file. SAAYN uses a multi-tier Remediation Loop to ensure that every node is buildable and logically sound before it ever touches your disk.

1. The Local Physics Audit (ast.go)
    - Before the AI is even allowed to review the logic
    - SAAYN performs a Structural Walk using go/ast.
    - Graph Integrity: Does this file break internal dependencies?
    - Interface Compliance: Does the struct actually satisfy the interface defined in the Specbook?
    - Syntax Lock: Is the code buildable?
    - Result: If the "Physics" fail, SAAYN auto-rejects the draft and forces the LLM to provide a syntactically valid version.

2. The Cognitive Audit (LLM Review)
    - Once the code is physically sound, SAAYN triggers a Cognitive Review. It compares the draft against your readme.md intent.
    - Intent Drift: "The Readme says this should be distributed, but you used a local global variable. Fix it."
    - Business Purpose: "The Specbook says this function must be idempotent. Prove it."

3. The Remediation Cycle
    - If the Cognitive Audit fails, SAAYN enters a Self-Correction Loop (Max 3 iterations).
    - It feeds the AST errors and the Review Findings back to the LLM to generate a surgical patch.

🖥️ Genesis in Action (Live CLI Trace)


```Plaintext
$ ./saayn greenfield -f vision.md -s spec.yaml --target ./task-bot

🧬 PHASE 0: CONTEXTUAL INGESTION
--------------------------------------------------------------------------------
📄 Vision:   Found 'Distributed Worker' intent in vision.md
📜 Physics:  12 Nodes identified in spec.yaml
🏗️ Build Order: [model] -> [registry] -> [worker] -> [main]

🌱 MATERIALIZING GENOME (INNER LOOP ACTIVE)
--------------------------------------------------------------------------------

[03/12] PROCESSING: internal/worker/worker.go

    🔬 DRAFTING: Initializing node 'saayn.Worker.Start'...
    
    ⚖️  AST AUDIT: Walking via ast.go...
       ├─ Syntax Check... ✅
       └─ Interface Check (Worker)... ✅

    🧠 COGNITIVE REVIEW: Analyzing against 'Distributed' intent...
       └─ 🚩 FINDING: 
          "Runner uses time.Sleep. Vision requires context-aware 
           cancellation for cloud-native compliance."

    🔧 REMEDIATION (Iteration 1/3):
       ├─ Resubmitting Feedback to DEEP LLM...
       ├─ Applying AST Patch...
       └─ Re-verifying Physics... ✅

    🧠 FINAL REVIEW:
       └─ ✅ Intent verified. select{} block implemented.

    💾 COMMIT: Writing to disk...
       └─ ✅ Logic Hash: d4e5f6 | Identity: worker.go [Indexed]

    🧠 Starting Semantic Enrichment Process...
    ⚙️   saayn.workerSomething[worker.go]        intneral/worker/worker.go       [21939968]
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
--------------------------------------------------------------------------------
✅ GENESIS COMPLETE: 12 Nodes Materialized. 1 Remediated.
📊 Build Status: 100% PASS
```

### 🏗️ Why Developers Use SAAYN Greenfield
🛡️ Hallucination-Proof
Because SAAYN walks the AST (Abstract Syntax Tree) locally, it is physically impossible for the agent to "hallucinate" a package that doesn't exist or a function signature that doesn't match the spec.

### 🧠 Short-Term Memory Preservation
SAAYN treats the CLI Log as its "Memory." By printing the Logic Hash and Review Findings for every node, the agent maintains a persistent state across fresh context windows. If the process is interrupted, you simply point SAAYN at the genome.json and it picks up exactly where the last hash left off.

### ⚖️ Deterministic Evolution
The Identity Triad (PublicID, Fingerprint, Logic Hash) ensures that your project is born with an audit trail. You can run saayn verify one second after genesis and see a perfect 1:1 match between your Vision and your Reality.

### 🛠️ Getting Started
0. Define your Soul: Write a readme.md describing your app.
1. Define your Physics: Write a specbook.yaml defining your packages.
2. Spawn: Run saayn genesis.
3. Stop coding in the dark. Start sequencing your genome.


### The 4-State Genome State Machine

| Genome State | What It Does (Action & Focus) | Who Is In Charge? | Agent Guardrails (Strict Enforcement) |
| :--- | :--- | :--- | :--- |
| **1. Conceptual** | Defines the "Gallery." Debates and locks in the `business_purpose` and architectural boundaries of the package. | Actor/Critic (DEEP) | **No Code Allowed:** The agent must reject any payload containing Go code. <br><br>**Dependency Lock:** If the LLM requests a dependency not in the Gallery, the agent pauses all states and forces the new dependency through State 1 first. |
| **2. Hollow (Canvas)** | Generates the structural topography. Writes interfaces, empty structs, and zero-return functions to satisfy the compiler. | Skeleton Builder (FAST) | **Zero-Logic Rule:** The agent must reject the AST payload if it detects any logic loops (`for`, `if`), variable assignments, or non-empty struct fields.<br><br>**Physics Gate:** Must compile perfectly before advancing to State 3. |
| **3. Hydrating** | The surgical phase. Analyzes the hollow signatures and writes the actual execution logic, algorithms, and data routing. | Surgeon (DEEP) | **State 1 Isolation:** The agent must physically block the LLM from calling any package that is still in State 1.<br><br>**Signature Lock:** If the Surgeon tries to change a function signature established in State 2, the agent rejects the patch and forces a Canvas re-stretch. |
| **4. Sequenced** | The logic is written, the AST passes, and the code fulfills the semantic intent of the `specbook.yaml`. | CC Agent (Local orchestrator) | **Hash Locking:** The agent calculates and commits the final Identity Triad to `genome.json`.<br><br>**Drift Detection:** If a human manually edits the file later, the hash breaks, and the agent automatically drops the node back to State 3. |

---

### Why the "Signature Lock" Guardrail is the Masterstroke

If you look at State 3, the **Signature Lock** is the most critical guardrail your agent enforces. 

When the DEEP model is writing the logic inside `StartWorker()`, it might get frustrated and try to change the function signature to make its life easier (e.g., adding `context.Context` when it wasn't in the skeleton). 

If the CC Agent allows that, it breaks the AST for every *other* file that is currently relying on the original hollow signature. By enforcing the Signature Lock, the CC Agent acts like a strict project manager: *"You cannot change the contract without scheduling a meeting with the rest of the architecture."* <FollowUp label="Want to map out the 'Signature Lock' override protocol?" query="If the DEEP model genuinely needs to change a signature in State 3 because the FAST model made a mistake in State 2, how does the CC Agent safely handle that without breaking the AST?" />


If you made it this far, then here is the real vision:

```
Tony Stark (leaning back in his chair, swirling a drink, holographic displays lighting up):

        "SAAYN, you up?"

SAAYN:  "For you, sir… always."

Tony:   "Alright, listen up. I’ve got a vision burning a hole in my head and a specbook 
         that’s tighter than the Mark 42’s flight stabilizers. Drop everything. Ingest 
         the readme.md like it’s the arc reactor blueprint. Parse the specbook.yaml, then 
         fire up that SAAYN Genesis Engine.  I want a full Go project materialized 
         — clean, deterministic, hallucination-proof. Run the Surgical Inner Loop. Audit 
         the physics with AST verified Physics, do the cognitive review, and if anything 
         even thinks about drifting from my intent… remediate it. Three iterations max.
         
         Make it elegant. Make it fast. Make it me.
         
         Thrill me, SAAYN."
         
         (He stands up, grabs his jacket, and heads for the door)
         
         "And don’t wait up. Daddy’s got a board meeting… or a cheeseburger. Whatever comes first."
```
