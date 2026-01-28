# Exercise 1: Complete the GLEntry Struct

**Difficulty:** 🟢 Beginner
**Time:** 30 minutes
**Goal:** Learn the GL Entry data structure by completing missing fields

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

## Next Exercise

After completing this, move to:
```bash
cd ../02-validation
```
