# Exercise 3: Repository Pattern

**Difficulty:** 🟡 Intermediate
**Time:** 1-2 hours
**Goal:** Implement a repository interface for storing GL entries

---

## 🔗 How This Connects to the Big Picture

> *You're learning Ports & Adapters — the pattern that makes incremental migration possible.*

### Architectural Connection (Goal 1: Legacy Modernization)

From Robert C. Martin's **Clean Architecture** and the [main README](../../README.md#but-the-accounts-module-has-dependencies):

The **Repository Pattern** is how we extract the Accounts module WITHOUT migrating Stock, Selling, and HR. We define **interfaces** (ports) at the boundary:

```go
// This is the KEY PATTERN from the main README
type AccountLookup interface {
    GetAccount(name string) (*Account, error)
    IsDisabled(name string) (bool, error)
}
```

Now the GL Entry Engine works with **any implementation**:
- `MockAccountLookup` for tests (instant, no DB needed)
- `PostgresAccountLookup` for production
- `ERPNextAPIAccountLookup` during migration (calls Python)

**This is how the Strangler Fig Pattern works in practice.**

### Code Intelligence Connection (Goal 2: AI Assistants)

Your Code Intelligence Platform needs **pluggable storage**. Today it might store embeddings in LanceDB, tomorrow in Weaviate. The Repository Pattern makes this possible:

```
┌─────────────────────────────────────────────────────────────────────────┐
│  YOUR CODE INTELLIGENCE TOOL ARCHITECTURE                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────┐      ┌──────────────────┐      ┌─────────────────┐    │
│  │ CLI/MCP     │ ───► │  IndexService    │ ───► │ VectorStore     │    │
│  │ Interface   │      │  (business logic)│      │ (interface)     │    │
│  └─────────────┘      └──────────────────┘      └────────┬────────┘    │
│                                                          │              │
│                              ┌────────────────────────────┼─────────┐   │
│                              │                            │         │   │
│                              ▼                            ▼         ▼   │
│                       ┌───────────┐              ┌──────────┐ ┌──────┐ │
│                       │ LanceDB   │              │ Weaviate │ │ Mock │ │
│                       │ (local)   │              │ (cloud)  │ │(test)│ │
│                       └───────────┘              └──────────┘ └──────┘ │
│                                                                          │
│  The IndexService doesn't know which storage it's using.                │
│  That's the power of the Repository Pattern.                            │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**This exercise teaches the pattern** you'll use in your own tool.

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

---

## 🧠 Code Intelligence Insight

When your tool analyzes ERPNext, you'll find data access scattered everywhere:

```python
# ERPNext pattern - direct database access
def get_gl_entries(voucher_no):
    return frappe.db.sql("""
        SELECT * FROM `tabGL Entry`
        WHERE voucher_no = %s
    """, voucher_no, as_dict=1)
```

This is **hard to analyze** because:
1. SQL is embedded in strings (not AST-parseable)
2. Table relationships are implicit
3. No clear interface boundary

**Your tool must extract data access patterns** even when they're messy. Understanding the Repository Pattern helps you:
1. Recognize when code *should* use repositories (but doesn't)
2. Suggest refactoring to cleaner patterns
3. Map data flows through the codebase

**In your tool**, you'll likely build:
```typescript
// Your Code Intelligence tool (TypeScript example)
interface SymbolStore {
    save(symbol: Symbol): Promise<void>;
    findByName(name: string): Promise<Symbol[]>;
    findByFile(filePath: string): Promise<Symbol[]>;
}

// Mock for testing
class InMemorySymbolStore implements SymbolStore { ... }

// Real storage
class LanceDBSymbolStore implements SymbolStore { ... }
```

**Same pattern, different domain.** Learn it here, apply it there.
