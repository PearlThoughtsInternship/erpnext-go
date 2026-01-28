# Exercise 2: Implement Balance Validation

**Difficulty:** 🟡 Intermediate
**Time:** 1-2 hours
**Goal:** Write functions that validate GL entries balance correctly

---

## Background

The **golden rule of accounting**: Debits must equal Credits.

Every transaction in a business creates at least two GL entries that balance:
- Buy inventory → Debit Inventory, Credit Cash
- Make sale → Debit Receivables, Credit Sales
- Pay salary → Debit Salary Expense, Credit Cash

If entries don't balance, something is wrong!

---

## Your Task

Implement validation functions in `validation.go`:

1. `TotalDebit(entries []GLEntry) float64` - Sum all debit amounts
2. `TotalCredit(entries []GLEntry) float64` - Sum all credit amounts
3. `IsBalanced(entries []GLEntry) bool` - Check if debits = credits
4. `Difference(entries []GLEntry) float64` - Calculate debit - credit

---

## Instructions

```bash
# Read the skeleton code
cat validation.go

# Run tests (they'll fail initially)
go test -v

# Implement the functions one by one
# Run tests after each implementation
```

---

## Success Criteria

```
=== RUN   TestTotalDebit
--- PASS: TestTotalDebit (0.00s)
=== RUN   TestTotalCredit
--- PASS: TestTotalCredit (0.00s)
=== RUN   TestIsBalanced
--- PASS: TestIsBalanced (0.00s)
=== RUN   TestDifference
--- PASS: TestDifference (0.00s)
PASS
```

---

## Hints

<details>
<summary>Hint 1: Looping in Go</summary>

```go
func TotalDebit(entries []GLEntry) float64 {
    var total float64
    for _, entry := range entries {
        total += entry.Debit
    }
    return total
}
```

</details>

<details>
<summary>Hint 2: Floating Point Comparison</summary>

Due to floating point precision, don't compare floats directly:

```go
// BAD
return totalDebit == totalCredit

// GOOD - use a small tolerance
const epsilon = 0.0001
return math.Abs(totalDebit - totalCredit) < epsilon
```

</details>

<details>
<summary>Hint 3: Empty Slice Edge Case</summary>

What should happen with an empty slice?
- `TotalDebit([])` → 0
- `TotalCredit([])` → 0
- `IsBalanced([])` → true (0 = 0)

</details>

---

## What You'll Learn

- Go slice iteration
- Function implementation
- Floating point comparison
- Edge case handling
- Writing idiomatic Go

---

## Real-World Context

In ERPNext Python:

```python
def check_if_in_list(gl_map):
    total_debit = sum(flt(d.debit) for d in gl_map)
    total_credit = sum(flt(d.credit) for d in gl_map)

    if abs(total_debit - total_credit) > 0.001:
        frappe.throw("Debit and Credit not equal")
```

Your Go implementation should produce the same results!
