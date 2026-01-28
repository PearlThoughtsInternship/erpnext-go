# Exercise 2: Implement Balance Validation

**Difficulty:** 🟡 Intermediate
**Time:** 1-2 hours
**Goal:** Write functions that validate GL entries balance correctly

---

## 🔗 How This Connects to the Big Picture

> *You're learning Characterization Testing — capturing existing behavior before changing code.*

### Architectural Connection (Goal 1: Legacy Modernization)

From Michael Feathers' **Working Effectively with Legacy Code**: Before you can change legacy code, you must **characterize** its current behavior. This validation logic (`debit == credit`) is a **business rule** that must work identically in Python and Go.

By implementing this validation yourself, you understand:
- What the golden rule of accounting actually means in code
- How ERPNext enforces business rules (the `frappe.throw()` pattern)
- What a **Characterization Test** captures

**See the [main README](../../README.md#how-we-handle-dependencies)** for how test doubles enable isolated testing.

### Code Intelligence Connection (Goal 2: AI Assistants)

The **#1 thing** enterprises need from legacy modernization is: **"What business rules are hidden in this code?"**

This validation logic is a **business rule**. It's not documented anywhere — it's embedded in Python code. Your Code Intelligence Platform must:

```
┌─────────────────────────────────────────────────────────────────────────┐
│  THE BUSINESS RULE EXTRACTION PROBLEM                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ERPNext has THOUSANDS of business rules like this:                     │
│                                                                          │
│  • "Debits must equal Credits" (this exercise)                          │
│  • "Can't submit invoice if credit limit exceeded"                      │
│  • "Stock can't go negative without allow_negative_stock setting"       │
│  • "Journal Entry must have at least 2 lines"                           │
│                                                                          │
│  Your Tool's Job:                                                        │
│  1. IDENTIFY functions that validate/check/raise errors                 │
│  2. EXTRACT the conditions (if abs(total_debit - total_credit) > 0.001) │
│  3. SUMMARIZE as human-readable rules for AI assistants                 │
│                                                                          │
│  "The system enforces that total debits equals total credits,           │
│   with a tolerance of 0.001 for floating point errors."                 │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**By implementing this rule yourself**, you understand how business logic is structured in enterprise code.

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

---

## 🧠 Code Intelligence Insight

This Python code is **exactly what your tool will analyze**. Notice the patterns:

| Pattern | What Your Tool Should Extract |
|---------|-------------------------------|
| `sum(flt(d.debit) for d in gl_map)` | Aggregation over a collection |
| `abs(total_debit - total_credit) > 0.001` | Validation rule with tolerance |
| `frappe.throw("...")` | Error message = business rule description |

**Your indexer should recognize:**
1. `frappe.throw()` = validation failure point
2. The condition before `throw()` = the business rule
3. The message = human-readable description

When you build your Code Intelligence tool, you'll write code that:
```python
# Your future tool (pseudocode)
for function in parsed_ast:
    if contains_call(function, "frappe.throw"):
        condition = extract_preceding_if_condition(function)
        message = extract_throw_message(function)
        rules.append(BusinessRule(condition, message))
```

**This exercise teaches you what those rules look like** so you can build extractors that find them.
