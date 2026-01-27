# 🎯 Design Decisions

> Design choices, trade-offs, and rationale for ERPNext to Go migration

---

## Table of Contents

- [Design Philosophy](#design-philosophy)
- [Key Decisions](#key-decisions)
- [Pattern Catalog](#pattern-catalog)
- [Trade-off Analysis](#trade-off-analysis)
- [Anti-Patterns to Avoid](#anti-patterns-to-avoid)
- [Decision Records](#decision-records)

---

## Design Philosophy

### Core Tenets

| Tenet | Description | Enforcement |
|-------|-------------|-------------|
| 🎯 **Parity First** | Match ERPNext behavior exactly before optimizing | Tests verify identical outputs |
| 🧪 **Test-Driven** | No code without tests | CI blocks merges without coverage |
| 📦 **Small Batches** | One DocType at a time | Each iteration is deployable |
| 🔌 **Loose Coupling** | Interfaces everywhere | Mocks for unit tests |
| 🚫 **No Magic** | Explicit over implicit | No reflection, no codegen |

### The Strangler's Oath

> "We will not break production. Every change is reversible. Tests prove parity before cutover."

---

## Key Decisions

### 1. 🔌 Interface-Based Dependencies

**Decision:** All external dependencies accessed via interfaces defined in the domain layer.

**Context:**
- ERPNext uses global functions (`frappe.get_value()`, `frappe.db.sql()`)
- Global state makes unit testing impossible without full stack
- We need to test business logic in isolation

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. Global functions (like Python) | Familiar to ERPNext devs | Untestable |
| B. Dependency injection via constructors | Testable, explicit | Verbose |
| **C. Interface-based injection** ✅ | Testable, swappable | Slight indirection |

**Decision:** Option C — Interfaces define contracts, implementations injected

```go
// Domain defines what it needs (port)
type AccountLookup interface {
    GetAccountCompany(accountName string) (string, error)
}

// Infrastructure provides it (adapter)
type PostgresAccountLookup struct { db *sql.DB }
func (p *PostgresAccountLookup) GetAccountCompany(name string) (string, error) { ... }

// Tests use mocks
type MockAccountLookup struct { accounts map[string]string }
func (m *MockAccountLookup) GetAccountCompany(name string) (string, error) { ... }
```

**Consequences:**
- ✅ Unit tests run in milliseconds (no DB)
- ✅ Can swap PostgreSQL for SQLite in tests
- ✅ Legacy bridge is just another adapter
- ⚠️ More files/types to manage

---

### 2. 🚨 Typed Sentinel Errors

**Decision:** Define error constants for programmatic error handling.

**Context:**
- Python uses `frappe.throw(message)` with string matching
- Callers can't reliably handle specific errors
- Go's `errors.Is()` enables type-safe error checks

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. Return string errors | Simple | No type safety |
| B. Custom error types | Full control | Boilerplate |
| **C. Sentinel errors + wrapper** ✅ | Type-safe, details included | Moderate complexity |

**Decision:** Option C — Sentinel errors with `ValidationError` wrapper

```go
// Sentinel errors for type checking
var (
    ErrDuplicateCompany = errors.New("same company is entered more than once")
    ErrAccountMismatch  = errors.New("account does not match with company")
)

// Wrapper for additional context
type ValidationError struct {
    Err     error  // Sentinel for errors.Is()
    Details string // Human-readable context
}

func (e *ValidationError) Unwrap() error { return e.Err }

// Usage
if errors.Is(err, ErrDuplicateCompany) {
    // Handle duplicate company specifically
}
```

**Consequences:**
- ✅ Callers can handle specific errors
- ✅ Details preserved for logging/display
- ✅ Works with `errors.Is()` and `errors.As()`
- ⚠️ Must remember to use Unwrap()

---

### 3. 📊 Table-Driven Tests

**Decision:** Use table-driven tests for comprehensive coverage.

**Context:**
- ERPNext has minimal test coverage for Mode of Payment
- We need to prove parity with Python behavior
- Adding test cases should be trivial

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. One function per test | Clear names | Verbose, hard to add cases |
| B. BDD framework (Ginkgo) | Expressive | External dependency |
| **C. Table-driven tests** ✅ | Go idiom, easy to extend | Slightly dense |

**Decision:** Option C — Table-driven with subtests

```go
func TestValidateRepeatingCompanies(t *testing.T) {
    tests := []struct {
        name     string
        accounts []ModeOfPaymentAccount
        wantErr  error
    }{
        {"empty accounts - valid", []ModeOfPaymentAccount{}, nil},
        {"duplicate - error", []ModeOfPaymentAccount{{Company: "A"}, {Company: "A"}}, ErrDuplicateCompany},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &ModeOfPayment{Accounts: tt.accounts}
            err := m.ValidateRepeatingCompanies()
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

**Consequences:**
- ✅ Adding test cases is one line
- ✅ Subtests run in parallel if needed
- ✅ Clear failure messages with test names
- ⚠️ Table setup can be verbose for complex inputs

---

### 4. 📦 Domain Purity

**Decision:** Domain structs contain no infrastructure code.

**Context:**
- Frappe's Document class mixes data, validation, persistence, and UI
- This tight coupling prevents reuse and testing
- We want business logic portable across contexts

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. Active Record (like Frappe) | Familiar, less code | Untestable, coupled |
| B. Anemic domain model | Simple structs | Logic scattered |
| **C. Rich domain model** ✅ | Logic with data, no infra | More design effort |

**Decision:** Option C — Structs with validation methods, no DB code

```go
// ✅ Domain struct - pure data and business rules
type ModeOfPayment struct {
    Name     string
    Type     PaymentType
    Enabled  bool
    Accounts []ModeOfPaymentAccount
}

// ✅ Validation is a domain concern
func (m *ModeOfPayment) Validate(lookup AccountLookup, checker POSChecker) error {
    // Business rules only - no DB calls directly
}

// ❌ NO persistence in domain
// func (m *ModeOfPayment) Save() error { ... }  // FORBIDDEN
```

**Consequences:**
- ✅ Domain logic testable without database
- ✅ Same validation works in API, CLI, batch jobs
- ✅ Clear separation of concerns
- ⚠️ More layers (domain, application, infrastructure)

---

### 5. 🔄 Incremental Migration

**Decision:** Migrate one bounded context at a time with feature flags.

**Context:**
- Big-bang rewrites have high failure rate
- Business needs continuity during migration
- We need ability to rollback any change

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. Big bang rewrite | Clean slate | High risk, long timeline |
| B. Branch by abstraction | Gradual | Complex branching |
| **C. Strangler Fig** ✅ | Incremental, reversible | Dual maintenance |

**Decision:** Option C — Strangler Fig with feature flags

```
Migration Timeline:
═══════════════════════════════════════════════════════════════

Week 1-2: Mode of Payment
├── Extract to Go ────────────────────► Tests pass
├── Deploy shadow mode ───────────────► Compare outputs
├── Enable for 1% traffic ────────────► Monitor errors
├── Enable for 100% traffic ──────────► Full cutover
└── Remove Python code ───────────────► Clean up

Week 3-4: Bank
├── Extract to Go ...
```

**Consequences:**
- ✅ Production never breaks
- ✅ Rollback is one config change
- ✅ Team learns patterns incrementally
- ⚠️ Two systems running in parallel temporarily

---

## Pattern Catalog

### Patterns Used

| Pattern | Where | Purpose |
|---------|-------|---------|
| **Repository** | `AccountLookup`, `POSChecker` | Abstract data access |
| **Value Object** | `PaymentType` | Type-safe enums |
| **Entity** | `ModeOfPayment` | Identity + behavior |
| **Aggregate** | `ModeOfPayment` + `Accounts` | Transactional boundary |
| **Specification** | `Validate()` methods | Encapsulate business rules |
| **Anti-Corruption Layer** | Legacy bridge | Protect from ERPNext changes |

### Pattern Details

#### Repository Pattern

```go
// Generic repository interface
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    Get(ctx context.Context, id string) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filters ...Filter) ([]*T, error)
}

// Specific implementation
type ModeOfPaymentRepository interface {
    Repository[ModeOfPayment]
    FindByType(ctx context.Context, paymentType PaymentType) ([]*ModeOfPayment, error)
}
```

#### Value Object Pattern

```go
// PaymentType is immutable and compared by value
type PaymentType string

const (
    Cash    PaymentType = "Cash"
    Bank    PaymentType = "Bank"
    General PaymentType = "General"
    Phone   PaymentType = "Phone"
)

// Validation at creation
func ParsePaymentType(s string) (PaymentType, error) {
    switch s {
    case "Cash", "Bank", "General", "Phone":
        return PaymentType(s), nil
    default:
        return "", fmt.Errorf("invalid payment type: %s", s)
    }
}
```

#### Specification Pattern

```go
// Each validation is a specification
func (m *ModeOfPayment) ValidateRepeatingCompanies() error { ... }
func (m *ModeOfPayment) ValidateAccounts(lookup AccountLookup) error { ... }
func (m *ModeOfPayment) ValidatePOSModeOfPayment(checker POSChecker) error { ... }

// Composite specification
func (m *ModeOfPayment) Validate(lookup AccountLookup, checker POSChecker) error {
    if err := m.ValidateAccounts(lookup); err != nil {
        return err
    }
    if err := m.ValidateRepeatingCompanies(); err != nil {
        return err
    }
    if err := m.ValidatePOSModeOfPayment(checker); err != nil {
        return err
    }
    return nil
}
```

---

## Trade-off Analysis

### Complexity vs. Testability

```
                    High Testability
                          │
                          │    ✅ Our Choice
                          │    (Interface-based DI)
                          │         ●
                          │
                          │
   Low Complexity ────────┼──────── High Complexity
                          │
              ●           │
        Global functions  │
        (Like Frappe)     │
                          │
                    Low Testability
```

### Performance vs. Maintainability

```
                    High Performance
                          │
                          │         ●
                          │    Hand-optimized C
                          │
              ●           │
         Go (our choice)  │
                          │
   Low Maintainability ───┼──────── High Maintainability
                          │
                          │
                          │              ●
                          │         Python/Frappe
                          │
                    Low Performance
```

### Migration Speed vs. Risk

```
                    Fast Migration
                          │
              ●           │
        Big Bang Rewrite  │
        (High Risk)       │
                          │
   High Risk ─────────────┼──────── Low Risk
                          │
                          │              ●
                          │    ✅ Strangler Fig
                          │    (Our Choice)
                          │
                    Slow Migration
```

---

## Anti-Patterns to Avoid

### ❌ Don't Do This

| Anti-Pattern | Why It's Bad | What To Do Instead |
|--------------|--------------|---------------------|
| **Global state** | Untestable, race conditions | Dependency injection |
| **String errors** | Can't handle programmatically | Typed sentinel errors |
| **God objects** | Hard to test, change, understand | Single responsibility |
| **Anemic domain** | Logic scattered across layers | Rich domain model |
| **Premature optimization** | Complexity without proof | Profile first |
| **Copy-paste ERPNext** | Inherit bad patterns | Redesign for Go idioms |

### Code Examples

```go
// ❌ BAD: Global state
var db *sql.DB  // Package-level variable

func GetAccount(name string) (*Account, error) {
    return db.Query("SELECT ...")  // How to test?
}

// ✅ GOOD: Dependency injection
type AccountService struct {
    repo AccountRepository
}

func (s *AccountService) GetAccount(ctx context.Context, name string) (*Account, error) {
    return s.repo.Get(ctx, name)  // repo can be mocked
}
```

```go
// ❌ BAD: String error checking
if err.Error() == "same company is entered more than once" {
    // Fragile!
}

// ✅ GOOD: Typed error checking
if errors.Is(err, ErrDuplicateCompany) {
    // Robust!
}
```

```go
// ❌ BAD: Mixing concerns (Active Record)
type ModeOfPayment struct {
    Name string
    db   *sql.DB  // Infrastructure in domain!
}

func (m *ModeOfPayment) Save() error {
    return m.db.Exec("INSERT ...")
}

// ✅ GOOD: Separation of concerns
type ModeOfPayment struct {
    Name string
    // No DB reference
}

type ModeOfPaymentRepository interface {
    Save(ctx context.Context, m *ModeOfPayment) error
}
```

---

## Decision Records

### ADR-001: Use Go for Modernization

**Status:** Accepted
**Date:** 2026-01-27

**Context:**
ERPNext is written in Python/Frappe. We need a modernization target language.

**Decision:**
Use Go for the modernized services.

**Consequences:**
- Team needs Go training
- Better performance and deployment story
- Strong typing catches bugs at compile time

---

### ADR-002: PostgreSQL over MariaDB

**Status:** Accepted
**Date:** 2026-01-27

**Context:**
ERPNext uses MariaDB. Go has better PostgreSQL tooling.

**Decision:**
Use PostgreSQL for Go services, sync data during migration.

**Consequences:**
- Need data sync mechanism
- Better Go ecosystem support
- Can use PostgreSQL-specific features

---

### ADR-003: Interface-Based Testing

**Status:** Accepted
**Date:** 2026-01-27

**Context:**
Need to test business logic without database.

**Decision:**
Define interfaces for all external dependencies.

**Consequences:**
- More interface definitions
- Mocks enable fast unit tests
- Production implementations are pluggable

---

## References

- [Domain-Driven Design](https://www.domainlanguage.com/ddd/) — Eric Evans
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) — Robert C. Martin
- [Effective Go](https://golang.org/doc/effective_go) — Go Team
- [Architecture Decision Records](https://adr.github.io/) — ADR GitHub
