# Exercise 3: Repository Pattern

**Difficulty:** 🟡 Intermediate
**Time:** 1-2 hours
**Goal:** Implement a repository interface for storing GL entries

---

## Background

The **Repository Pattern** abstracts data storage. Instead of your business logic knowing about databases, it uses an interface:

```go
// Business logic uses the interface
type GLEntryStore interface {
    Save(entry *GLEntry) error
    GetByVoucher(voucherType, voucherNo string) ([]GLEntry, error)
}

// Different implementations for different contexts
type PostgresStore struct { db *sql.DB }     // Production
type InMemoryStore struct { entries []GLEntry }  // Testing
```

This is a core principle of **Hexagonal Architecture**.

---

## Your Task

Implement an in-memory repository that stores GL entries.

1. `Save(entry *GLEntry) error` - Add entry to storage
2. `GetByVoucher(voucherType, voucherNo string) ([]GLEntry, error)` - Find entries
3. `GetAll() []GLEntry` - Return all stored entries
4. `Clear()` - Remove all entries

---

## Files

- `repository.go` - Skeleton to complete
- `repository_test.go` - Tests to pass

---

## Success Criteria

```
=== RUN   TestInMemoryStore_Save
--- PASS: TestInMemoryStore_Save (0.00s)
=== RUN   TestInMemoryStore_GetByVoucher
--- PASS: TestInMemoryStore_GetByVoucher (0.00s)
=== RUN   TestInMemoryStore_Clear
--- PASS: TestInMemoryStore_Clear (0.00s)
PASS
```

---

## Hints

<details>
<summary>Hint 1: Using a slice for storage</summary>

```go
type InMemoryStore struct {
    entries []GLEntry
}

func (s *InMemoryStore) Save(entry *GLEntry) error {
    s.entries = append(s.entries, *entry)
    return nil
}
```

</details>

<details>
<summary>Hint 2: Filtering by voucher</summary>

```go
func (s *InMemoryStore) GetByVoucher(voucherType, voucherNo string) ([]GLEntry, error) {
    var result []GLEntry
    for _, e := range s.entries {
        if e.VoucherType == voucherType && e.VoucherNo == voucherNo {
            result = append(result, e)
        }
    }
    return result, nil
}
```

</details>

---

## What You'll Learn

- Interface definition in Go
- Implicit interface implementation
- Repository pattern for data access
- How to make code testable
