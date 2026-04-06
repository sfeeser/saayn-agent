
```text
.
├── cmd/saayn/              # CLI logic
├── docs/
│   └── QA_PLAN.md          # <--- YOUR ARE HERE
├── internal/
│   ├── genome/             # Core logic
│   └── testrepo/           # <--- THE BYSTANDER (The "Patient")
│       └── bystander.go
└── ...
```

---

# V5 Certification Plan: The Innocent Bystander file

**Objective:** Verify the end-to-end "Surgical Inner Loop" using a non-critical target (bystander.go) to test AST integrity and semantic accuracy.

### Phase 1: Identity & Discovery
1. **Confirm `verify-llm-targets` works**
   `./saayn verify-llm-targets`
   *Expect: 3/3 Green.*

2. **Confirm `init` captures the target**
   `./saayn init`
   *Expect: 176+ nodes (including `testrepo.Calculator`).*

3. **Confirm `enrich` understands purpose**
   `./saayn enrich`
   *Expect: Summary for `Calculator.Add` identifies it as an arithmetic operation.*

4. **Confirm `search-intent` finds logic**
   `./saayn search-intent "math for adding integers"`
   *Expect: `saayn.testrepo.Calculator.Add` ranks #1.*

### Phase 2: The Surgical Loop
5. **Confirm `draft` translates intent**
   `./saayn draft "Modify the Calculator to support subtraction instead of addition"`
   *Expect: `surgery.yaml` generated with target `saayn.testrepo.Calculator.Add`.*

6. **Confirm `graph` hydrates context**
   `./saayn graph`
   *Expect: `innocent_bystander.go` source code added to the plan context.*

7. **Confirm `plan` generates the patch**
   `./saayn plan`
   *Expect: Valid Go code generated in `patch.yaml` using `-` instead of `+`.*

8. **Confirm `apply` executes AST Splicing**
   `./saayn apply`
   *Expect: Physical file `internal/testrepo/bystander.go` updated. No syntax errors.*

9. **Confirm `gen-test` verifies mutation**
   `./saayn gen-test saayn.testrepo.Calculator.Add`
   *Expect: Test passes for subtraction logic.*

---
