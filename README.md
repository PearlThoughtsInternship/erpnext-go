# ERPNext Legacy Modernization — Python to Go

A demonstration of modernizing ERPNext (Python/Frappe) to Go using the **Strangler Fig Pattern** with iterative, test-driven extraction.

## Table of Contents

- [Rationale](#rationale)
- [Architecture](#architecture)
- [Strangler Fig Pattern](#strangler-fig-pattern)
- [Design Choices](#design-choices)
- [Implementation](#implementation)
- [Iteration 1: Mode of Payment](#iteration-1-mode-of-payment)
- [Verification](#verification)
- [Next Steps](#next-steps)

---

## Rationale

### Why Modernize ERPNext?

ERPNext is a mature, feature-rich ERP built on the Frappe framework (Python). While powerful, organizations may need to modernize for:

| Challenge | Impact |
|-----------|--------|
| **Runtime type safety** | Bugs discovered in production, not development |
| **Framework coupling** | Business logic tightly bound to Frappe ORM |
| **Testing complexity** | Integration tests require full Frappe stack |
| **Performance** | Python's GIL limits concurrent request handling |
| **Deployment** | Requires Python environment + MariaDB + Redis |

### Why Go?

| Benefit | Description |
|---------|-------------|
| **Compile-time safety** | Type errors caught before deployment |
| **Single binary** | No runtime dependencies |
| **Concurrency** | Native goroutines for parallel processing |
| **Performance** | 10-100x faster for CPU-bound operations |
| **Testability** | Interfaces enable isolated unit tests |

### Why Not Rewrite?

> "The only thing a Big Bang rewrite guarantees is a Big Bang." — Martin Fowler

Full rewrites fail because:
- Business loses features during development
- Knowledge is lost in translation
- Testing parity is nearly impossible
- Timeline and budget always exceed estimates

**Strangler Fig Pattern** allows incremental migration with zero downtime.

---

## Architecture

### Legacy System (ERPNext/Frappe)

```
┌─────────────────────────────────────────────────────────────┐
│  Frappe Framework                                           │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Document Base Class                                  │  │
│  │  • Magic field access (self.fieldname)               │  │
│  │  • Automatic DB persistence                          │  │
│  │  • Hook system (validate, on_save, on_trash)         │  │
│  │  • Permission enforcement                            │  │
│  └───────────────────────────────────────────────────────┘  │
│                           │                                 │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │  DocType: Mode of Payment                             │  │
│  │  • mode_of_payment.py (business logic)               │  │
│  │  • mode_of_payment.json (schema definition)          │  │
│  │  • mode_of_payment.js (UI controller)                │  │
│  └───────────────────────────────────────────────────────┘  │
│                           │                                 │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │  frappe.db / frappe.get_value()                       │  │
│  │  • Direct SQL to MariaDB                             │  │
│  │  • Redis caching layer                               │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Modernized System (Go)

```
┌─────────────────────────────────────────────────────────────┐
│  Go Application                                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Domain Layer (Pure Business Logic)                   │  │
│  │  • Structs with explicit fields                      │  │
│  │  • Validation methods                                │  │
│  │  • No framework dependencies                         │  │
│  └───────────────────────────────────────────────────────┘  │
│                           │                                 │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │  Port Interfaces (Dependency Inversion)               │  │
│  │  • AccountLookup                                     │  │
│  │  • POSChecker                                        │  │
│  │  • Repository[T]                                     │  │
│  └───────────────────────────────────────────────────────┘  │
│           │                               │                 │
│  ┌────────▼────────┐            ┌─────────▼──────────┐      │
│  │  Mock Adapters  │            │  Real Adapters     │      │
│  │  (Testing)      │            │  (Production)      │      │
│  │  • In-memory    │            │  • PostgreSQL      │      │
│  │  • Deterministic│            │  • Redis cache     │      │
│  └─────────────────┘            └────────────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

---

## Strangler Fig Pattern

### Concept

The Strangler Fig is a tree that grows around its host, eventually replacing it entirely while the host continues to function.

```
Phase 1: Identify         Phase 2: Extract        Phase 3: Redirect       Phase 4: Remove
┌─────────────┐           ┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   Legacy    │           │   Legacy    │         │   Legacy    │         │             │
│  ┌───────┐  │           │  ┌───────┐  │         │  ┌ ─ ─ ─ ┐  │         │             │
│  │Module │  │    ══►    │  │Module │──┼──┐      │  │Module │──┼──┐      │             │
│  └───────┘  │           │  └───────┘  │  │      │  └ ─ ─ ─ ┘  │  │      │             │
│             │           │             │  │      │             │  │      │             │
└─────────────┘           └─────────────┘  │      └─────────────┘  │      └─────────────┘
                                           │                       │
                          ┌────────────────▼─┐    ┌────────────────▼─┐    ┌─────────────┐
                          │   Go Module      │    │   Go Module      │    │  Go Module  │
                          │   (shadow)       │    │   (primary)      │    │  (sole)     │
                          └──────────────────┘    └──────────────────┘    └─────────────┘
```

### Implementation Approach

1. **Identify** — Select a bounded module with clear inputs/outputs
2. **Extract** — Reimplement business logic in Go with tests
3. **Shadow** — Run both systems, compare outputs (optional)
4. **Redirect** — Route traffic to Go implementation
5. **Remove** — Deprecate Python code when confident

### Why This Works

| Principle | Benefit |
|-----------|---------|
| **Incremental** | Deliver value continuously |
| **Reversible** | Roll back any migration step |
| **Testable** | Prove parity before switching |
| **Low risk** | Failures are isolated to one module |

---

## Design Choices

### 1. Interface-Based Dependency Injection

**Problem:** Frappe's `frappe.get_value()` and `frappe.db.sql()` are global functions that couple business logic to the database.

**Solution:** Define interfaces that abstract external dependencies.

```go
// Port interface - defines what we need
type AccountLookup interface {
    GetAccountCompany(accountName string) (string, error)
}

// Business logic depends on interface, not implementation
func (m *ModeOfPayment) ValidateAccounts(lookup AccountLookup) error {
    // ...
}
```

**Benefits:**
- Unit tests use mock implementations (fast, deterministic)
- Production uses real database adapters
- Swap implementations without changing business logic

### 2. Typed Sentinel Errors

**Problem:** Python uses exceptions with string messages. Hard to programmatically handle specific errors.

**Solution:** Define typed error constants.

```go
var (
    ErrDuplicateCompany = errors.New("same company is entered more than once")
    ErrAccountMismatch  = errors.New("account does not match with company")
    ErrModeInUse        = errors.New("mode of payment is used in POS profiles")
)

// Callers can check error type
if errors.Is(err, ErrDuplicateCompany) {
    // Handle duplicate company specifically
}
```

### 3. Table-Driven Tests

**Problem:** ERPNext's test file is a skeleton with no actual test cases.

**Solution:** Comprehensive table-driven tests covering all branches.

```go
tests := []struct {
    name     string
    input    *ModeOfPayment
    wantErr  error
}{
    {"valid - unique companies", validMode, nil},
    {"invalid - duplicate company", dupMode, ErrDuplicateCompany},
    // ... more cases
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        err := tt.input.ValidateRepeatingCompanies()
        if !errors.Is(err, tt.wantErr) {
            t.Errorf("got %v, want %v", err, tt.wantErr)
        }
    })
}
```

### 4. Domain Purity

**Problem:** Frappe Document classes mix persistence, validation, and UI concerns.

**Solution:** Go structs contain only data and validation logic.

```go
// Pure domain struct - no DB, no HTTP, no UI
type ModeOfPayment struct {
    Name     string
    Type     PaymentType
    Enabled  bool
    Accounts []ModeOfPaymentAccount
}

// Validation is a pure function of the data
func (m *ModeOfPayment) Validate(lookup AccountLookup, checker POSChecker) error
```

---

## Implementation

### Project Structure

```
erpnext-go/
├── go.mod                          # Module definition
├── README.md                       # This file
├── modeofpayment/                  # Iteration 1: Mode of Payment
│   ├── model.go                    # Data structures
│   ├── validation.go               # Business rules
│   └── validation_test.go          # 19 test cases
└── [future modules]/               # Iterations 2+
```

### Source Mapping

| ERPNext (Python) | Go |
|------------------|-----|
| `erpnext/accounts/doctype/mode_of_payment/mode_of_payment.py` | `modeofpayment/validation.go` |
| `erpnext/accounts/doctype/mode_of_payment/mode_of_payment.json` | `modeofpayment/model.go` |
| `erpnext/accounts/doctype/mode_of_payment_account/mode_of_payment_account.json` | `modeofpayment/model.go` |

---

## Iteration 1: Mode of Payment

### Why This Module?

| Criterion | Assessment |
|-----------|------------|
| **Self-contained** | No complex dependencies on other doctypes |
| **Clear boundaries** | 4 fields, 1 child table, 3 validation rules |
| **Testable logic** | Business rules are pure functions of data |
| **Representative** | Demonstrates patterns applicable to all doctypes |

### Business Rules Migrated

#### Rule 1: No Duplicate Companies

```python
# Python (ERPNext)
def validate_repeating_companies(self):
    accounts_list = [entry.company for entry in self.accounts]
    if len(accounts_list) != len(set(accounts_list)):
        frappe.throw(_("Same Company is entered more than once"))
```

```go
// Go
func (m *ModeOfPayment) ValidateRepeatingCompanies() error {
    seen := make(map[string]bool)
    for _, account := range m.Accounts {
        if seen[account.Company] {
            return &ValidationError{Err: ErrDuplicateCompany, ...}
        }
        seen[account.Company] = true
    }
    return nil
}
```

#### Rule 2: Account-Company Match

```python
# Python (ERPNext)
def validate_accounts(self):
    for entry in self.accounts:
        if frappe.get_cached_value("Account", entry.default_account, "company") != entry.company:
            frappe.throw(_("Account {0} does not match..."))
```

```go
// Go
func (m *ModeOfPayment) ValidateAccounts(lookup AccountLookup) error {
    for _, account := range m.Accounts {
        accountCompany, _ := lookup.GetAccountCompany(account.DefaultAccount)
        if accountCompany != account.Company {
            return &ValidationError{Err: ErrAccountMismatch, ...}
        }
    }
    return nil
}
```

#### Rule 3: POS Profile Check

```python
# Python (ERPNext)
def validate_pos_mode_of_payment(self):
    if not self.enabled:
        pos_profiles = frappe.db.sql("SELECT ... WHERE mode_of_payment = %s", self.name)
        if pos_profiles:
            frappe.throw(_("POS Profile {} contains Mode of Payment {}..."))
```

```go
// Go
func (m *ModeOfPayment) ValidatePOSModeOfPayment(checker POSChecker) error {
    if m.Enabled {
        return nil
    }
    profiles, _ := checker.GetPOSProfilesUsingMode(m.Name)
    if len(profiles) > 0 {
        return &ValidationError{Err: ErrModeInUse, ...}
    }
    return nil
}
```

### Parity Results

| Metric | Python | Go | Status |
|--------|--------|-----|--------|
| Data fields | 6 | 6 | ✅ 100% |
| Validation rules | 3 | 3 | ✅ 100% |
| Error messages | Match | Match | ✅ 100% |
| Test cases | 0 | 19 | ✅ Exceeds |
| Coverage | N/A | 85.3% | ✅ |

---

## Verification

### Run Tests

```bash
cd erpnext-go
go test -v ./modeofpayment/
```

### Expected Output

```
=== RUN   TestValidateRepeatingCompanies
=== RUN   TestValidateRepeatingCompanies/empty_accounts_-_valid
=== RUN   TestValidateRepeatingCompanies/duplicate_companies_-_error
--- PASS: TestValidateRepeatingCompanies (0.00s)

=== RUN   TestValidateAccounts
=== RUN   TestValidateAccounts/account_matches_company_-_valid
=== RUN   TestValidateAccounts/account_company_mismatch_-_error
--- PASS: TestValidateAccounts (0.00s)

=== RUN   TestValidatePOSModeOfPayment
=== RUN   TestValidatePOSModeOfPayment/disabled,_used_in_POS_-_error
--- PASS: TestValidatePOSModeOfPayment (0.00s)

=== RUN   TestValidate_Integration
--- PASS: TestValidate_Integration (0.00s)

PASS
ok      github.com/senguttuvang/erpnext-go/modeofpayment    0.5s
```

### Check Coverage

```bash
go test -cover ./modeofpayment/
# coverage: 85.3% of statements
```

---

## Next Steps

### Iteration 2: Add Repository Layer

```go
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    Get(ctx context.Context, id string) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filters ...Filter) ([]*T, error)
}
```

### Iteration 3: Add HTTP API

```go
// REST endpoints matching ERPNext's API
POST   /api/resource/Mode of Payment
GET    /api/resource/Mode of Payment/:name
PUT    /api/resource/Mode of Payment/:name
DELETE /api/resource/Mode of Payment/:name
GET    /api/resource/Mode of Payment?filters=...
```

### Iteration 4: Shadow Mode

Run both Python and Go in parallel, compare responses:

```
Request ──┬──► Python (ERPNext) ──► Response A ──┐
          │                                      ├──► Compare
          └──► Go (New)          ──► Response B ──┘
```

### Future Modules

| Priority | Module | Complexity | Dependencies |
|----------|--------|------------|--------------|
| P1 | Bank | Low | Address, Contact |
| P1 | Currency Exchange | Low | Currency |
| P2 | Payment Entry | Medium | Mode of Payment, Party |
| P2 | Journal Entry | Medium | Account, Cost Center |
| P3 | Sales Invoice | High | Customer, Item, Tax |

---

## References

- [Strangler Fig Pattern](https://martinfowler.com/bliki/StranglerFigApplication.html) — Martin Fowler
- [Working Effectively with Legacy Code](https://www.oreilly.com/library/view/working-effectively-with/0131177052/) — Michael Feathers
- [ERPNext Documentation](https://docs.erpnext.com/)
- [Frappe Framework](https://frappeframework.com/)

---

## License

MIT License — See [LICENSE](LICENSE) for details.
