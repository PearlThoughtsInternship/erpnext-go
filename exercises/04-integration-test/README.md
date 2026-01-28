# Exercise 4: Write an Integration Test

**Difficulty:** 🔴 Advanced
**Time:** 2-3 hours
**Goal:** Write a realistic integration test that proves Python/Go parity

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
