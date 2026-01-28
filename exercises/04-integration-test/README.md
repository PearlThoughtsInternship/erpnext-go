# Exercise 4: Write an Integration Test

**Difficulty:** 🔴 Advanced
**Time:** 2-3 hours
**Goal:** Write a realistic integration test that proves Python/Go parity

---

## 🔗 How This Connects to the Big Picture

> *You're learning Parity Verification — the Strangler Fig's proof of correctness.*

### Architectural Connection (Goal 1: Legacy Modernization)

From Sam Newman's **Monolith to Microservices**: The Strangler Fig Pattern requires **Shadow Mode** — running Python and Go in parallel and comparing outputs.

Look at the [main README](../../README.md#the-strangler-fig-in-action) for the full diagram:
1. **Phase 1 (Shadow)**: Both systems run, Go is compared against Python
2. **Phase 2 (Traffic Switch)**: Go becomes primary, Python is rollback-ready

**Parity tests are how you build confidence to make that switch.** This exercise teaches you to write them.

### Code Intelligence Connection (Goal 2: AI Assistants)

This exercise is **the most important one** for your Code Intelligence work. Here's why:

```
┌─────────────────────────────────────────────────────────────────────────┐
│  THE PARITY VERIFICATION CHALLENGE                                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Your Tool generates Go code from Python context:                       │
│                                                                          │
│     Python (ERPNext)                    Go (Generated)                  │
│     ════════════════                    ════════════════                │
│     on_submit() creates                 OnSubmit() creates              │
│     GL Entries with:                    GL Entries with:                │
│     - Debit: Receivables $1000          - Debit: Receivables $1000     │
│     - Credit: Sales $1000               - Credit: Sales $1000          │
│                                                                          │
│  HOW DO YOU PROVE THEY'RE EQUIVALENT?                                   │
│                                                                          │
│  Answer: PARITY TESTS                                                   │
│                                                                          │
│  1. Run same scenario on both systems                                   │
│  2. Compare outputs (GL entries, stock movements, etc.)                 │
│  3. If outputs match → Parity achieved ✓                               │
│                                                                          │
│  This exercise teaches you HOW TO WRITE THOSE TESTS.                    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Your tool's effectiveness is measured by parity test pass rate.** Learn to write them well.

---

## Background

Integration tests verify that components work together correctly. In our migration, they prove that Go behaves identically to Python.

A good integration test:
1. Uses realistic ERPNext data
2. Tests a complete business flow
3. Verifies specific outputs match Python
4. Documents the expected behavior

---

## Your Task

Write an integration test for a **Payment Entry** scenario:

When a customer pays an invoice:
- Cash/Bank account receives money (Debit)
- Customer's receivable account decreases (Credit)

**Example Payment of ₹11,800:**

| Account | Debit | Credit |
|---------|-------|--------|
| Bank - ACME | ₹11,800 | - |
| Debtors - ACME | - | ₹11,800 |

---

## Instructions

1. Look at `integration_test_template.go` for the structure
2. Complete the `TestPaymentEntryFlow` function
3. Add test cases for:
   - Basic payment
   - Partial payment
   - Payment with write-off

---

## Success Criteria

```
=== RUN   TestPaymentEntryFlow
=== RUN   TestPaymentEntryFlow/full_payment
--- PASS: TestPaymentEntryFlow/full_payment (0.00s)
=== RUN   TestPaymentEntryFlow/partial_payment
--- PASS: TestPaymentEntryFlow/partial_payment (0.00s)
PASS
```

---

## Reference

Look at the actual integration test for guidance:
- `/ledger/integration_test.go` - `TestRealisticPaymentEntryGLEntries`

---

## What You'll Learn

- Writing comprehensive integration tests
- Parity verification approach
- Table-driven test structure
- ERPNext accounting patterns

---

## 🧠 Code Intelligence Insight

Your Code Intelligence Platform's **success metric** is: *"Does the generated code pass parity tests?"*

Here's how your complete pipeline will work:

```
┌─────────────────────────────────────────────────────────────────────────┐
│  END-TO-END CODE INTELLIGENCE PIPELINE                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  PHASE 1: INDEX                                                         │
│  ───────────────                                                         │
│  Your Tool parses ERPNext Python → extracts symbols, relationships      │
│  Stores in vector DB for semantic search                                │
│                                                                          │
│  PHASE 2: RETRIEVE (when developer asks a question)                     │
│  ─────────────────                                                       │
│  "How does Payment Entry create GL entries?"                            │
│  → Search vector DB → Find relevant chunks                              │
│  → Build focused context (2000 tokens, not 20000)                       │
│                                                                          │
│  PHASE 3: GENERATE                                                       │
│  ────────────────                                                        │
│  Send context to LLM → "Generate equivalent Go code"                    │
│  LLM outputs Go structs, handlers, validation                           │
│                                                                          │
│  PHASE 4: VERIFY ← THIS EXERCISE                                        │
│  ───────────────                                                         │
│  Run parity tests against generated code                                │
│  ✓ Same GL entries created?                                             │
│  ✓ Same validation errors?                                              │
│  ✓ Same edge case handling?                                             │
│                                                                          │
│  EVIDENCE: "85% of generated Payment Entry code passes parity tests"   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Without parity tests, you can't prove your tool works.** This exercise teaches you the verification approach that will define your internship success metrics.

---

## 📊 Success Metrics (Your Evidence)

At the end of Week 4, you'll present:

| Metric | Target | How You'll Measure |
|--------|--------|-------------------|
| **Parity Test Pass Rate** | > 80% | Tests you write in this exercise |
| **Token Reduction** | > 70% | Your tool vs vanilla Claude Code |
| **Context Precision** | > 75% | Relevant chunks / total chunks |
| **Developer Satisfaction** | > 4/5 | Peer feedback on generated code |

**This exercise is where you build the verification layer** that proves your tool's value.
