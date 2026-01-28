# Exercise 1: Complete the GLEntry Struct

**Difficulty:** 🟢 Beginner
**Time:** 30 minutes
**Goal:** Learn the GL Entry data structure by completing missing fields

---

## 🔗 How This Connects to the Big Picture

> *You're learning Domain-Driven Design — how to model business concepts in code.*

### Architectural Connection (Goal 1: Legacy Modernization)

From Eric Evans' **Domain-Driven Design**: A **Bounded Context** has a clear domain model. `GLEntry` is the core entity of the Accounts bounded context.

When you understand the GLEntry struct, you understand:
- What data flows through the accounting system
- What fields are required vs optional
- How ERPNext models financial transactions

**See the [main README](../../README.md#but-the-accounts-module-has-dependencies)** for how this fits into the Bounded Context Strategy.

### Code Intelligence Connection (Goal 2: AI Assistants)

The **Code Intelligence Platform** you're building needs to understand enterprise domain concepts like `GLEntry`. Here's why this exercise matters:

```
┌─────────────────────────────────────────────────────────────────────────┐
│  WHAT YOUR TOOL WILL DO (eventually)                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Developer asks: "How does ERPNext create GL entries?"                  │
│                                                                          │
│  Your Tool's Job:                                                        │
│  1. FIND the relevant Python code (gl_entry.py, accounts_controller.py) │
│  2. EXTRACT the domain model (GLEntry fields, relationships)            │
│  3. PROVIDE context to AI assistant                                      │
│                                                                          │
│  AI Assistant outputs:                                                   │
│  "GLEntry has 20+ fields. Key ones are Account, Debit, Credit,          │
│   VoucherType. The balance rule is enforced in make_gl_entries()..."    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**By understanding GLEntry yourself**, you'll build better tools that extract and explain this knowledge to AI assistants.

---

## Background

In accounting, a **General Ledger Entry (GL Entry)** is the fundamental record of a financial transaction. Every time money moves in ERPNext, GL entries are created.

For example, a Sales Invoice creates these GL entries:

| Account | Debit | Credit |
|---------|-------|--------|
| Debtors | ₹11,800 | - |
| Sales | - | ₹10,000 |
| CGST | - | ₹900 |
| SGST | - | ₹900 |

**Notice:** Total Debit = Total Credit (₹11,800 = ₹11,800) - This is the golden rule!

---

## Your Task

Open `glentry.go` and complete the missing struct fields.

Each `// TODO` comment tells you:
1. The field name to add
2. The data type
3. What it represents

---

## Instructions

1. **Read the TODO comments** in `glentry.go`
2. **Add the missing fields** with correct types
3. **Run tests** to verify:

```bash
go test -v
```

---

## Success Criteria

When you run `go test -v`, you should see:

```
=== RUN   TestGLEntryFields
--- PASS: TestGLEntryFields (0.00s)
=== RUN   TestGLEntryDefaults
--- PASS: TestGLEntryDefaults (0.00s)
PASS
```

---

## Hints

<details>
<summary>Hint 1: Data Types</summary>

Common field types in this struct:
- `string` - for text (accounts, names)
- `float64` - for monetary amounts
- `time.Time` - for dates
- `bool` - for flags (is_cancelled, etc.)

</details>

<details>
<summary>Hint 2: Field Names</summary>

Go uses CamelCase for exported fields:
- Python: `posting_date` → Go: `PostingDate`
- Python: `is_cancelled` → Go: `IsCancelled`

</details>

<details>
<summary>Hint 3: Check the Real Model</summary>

Look at `/ledger/model.go` for the complete struct.
But try to figure it out yourself first!

</details>

---

## What You'll Learn

- How GL entries are structured
- Go struct syntax
- The relationship between Python and Go data models
- Running Go tests

---

## 🧠 Code Intelligence Insight

When your tool indexes ERPNext, it will extract struct definitions like this:

```python
# ERPNext Python (what your tool will parse)
class GLEntry(Document):
    account = ...
    debit = ...
    credit = ...
    posting_date = ...
```

Your tool's **AST parser** will need to:
1. Identify this as a **domain entity** (not just any class)
2. Extract the **field names and types**
3. Understand the **relationships** (GLEntry → Account)

**This exercise gives you the domain knowledge** to build that parser correctly. You'll know what fields matter, what they mean, and how they relate.

---

## Next Exercise

After completing this, move to:
```bash
cd ../02-validation
```
