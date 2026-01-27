# 🔧 Implementation Guide

> Step-by-step guide to migrating ERPNext DocTypes to Go

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Migration Workflow](#migration-workflow)
- [Step-by-Step Guide](#step-by-step-guide)
- [Code Patterns](#code-patterns)
- [Testing Strategy](#testing-strategy)
- [Deployment Guide](#deployment-guide)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Development Environment

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.21+ | Runtime |
| Git | 2.x | Version control |
| PostgreSQL | 15+ | Database (optional for iteration 1) |
| ERPNext | v14+ | Source system reference |

### Setup Commands

```bash
# Verify Go installation
go version
# go version go1.21.0 darwin/arm64

# Clone repository
git clone git@github.com:senguttuvang/erpnext-go.git
cd erpnext-go

# Run tests to verify setup
go test -v ./...

# Check coverage
go test -cover ./...
```

---

## Migration Workflow

### Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                    DocType Migration Workflow                       │
└─────────────────────────────────────────────────────────────────────┘

  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
  │ Analyze │───►│  Model  │───►│Validate │───►│  Test   │───►│ Deploy  │
  │ Python  │    │   Go    │    │   Go    │    │  Parity │    │ Shadow  │
  └─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘
       │              │              │              │              │
       ▼              ▼              ▼              ▼              ▼
  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
  │• Schema │    │• Structs│    │• Rules  │    │• Unit   │    │• Feature│
  │• Rules  │    │• Enums  │    │• Errors │    │• Integr │    │  Flag   │
  │• Tests  │    │• Types  │    │• Methods│    │• Parity │    │• Monitor│
  └─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘
```

### Checklist

```
□ Phase 1: Analysis
  □ Read Python source code
  □ Identify fields from JSON schema
  □ Document business rules
  □ Note external dependencies

□ Phase 2: Model
  □ Create Go structs
  □ Define type enums
  □ Add field comments

□ Phase 3: Validation
  □ Implement each business rule
  □ Define error types
  □ Create port interfaces

□ Phase 4: Testing
  □ Write table-driven tests
  □ Cover all branches
  □ Verify error messages match

□ Phase 5: Deployment
  □ Add HTTP handlers (if needed)
  □ Configure feature flag
  □ Deploy to shadow mode
  □ Monitor and compare
```

---

## Step-by-Step Guide

### Step 1: Analyze Python Source

**Goal:** Understand what we're migrating.

#### 1.1 Locate the Files

```
erpnext/
└── erpnext/
    └── accounts/
        └── doctype/
            └── mode_of_payment/
                ├── mode_of_payment.json    # Schema
                ├── mode_of_payment.py      # Business logic
                ├── test_mode_of_payment.py # Tests (if any)
                └── mode_of_payment.js      # UI (skip)
```

#### 1.2 Extract Schema from JSON

```json
// mode_of_payment.json
{
  "fields": [
    {"fieldname": "mode_of_payment", "fieldtype": "Data", "reqd": 1, "unique": 1},
    {"fieldname": "type", "fieldtype": "Select", "options": "Cash\nBank\nGeneral\nPhone"},
    {"fieldname": "enabled", "fieldtype": "Check", "default": "1"},
    {"fieldname": "accounts", "fieldtype": "Table", "options": "Mode of Payment Account"}
  ]
}
```

**Document the mapping:**

| JSON Field | Type | Go Type |
|------------|------|---------|
| `mode_of_payment` | Data (required, unique) | `string` |
| `type` | Select | `PaymentType` (enum) |
| `enabled` | Check | `bool` |
| `accounts` | Table | `[]ModeOfPaymentAccount` |

#### 1.3 Extract Business Rules from Python

```python
# mode_of_payment.py
class ModeofPayment(Document):
    def validate(self):
        self.validate_accounts()
        self.validate_repeating_companies()
        self.validate_pos_mode_of_payment()
```

**Document the rules:**

| Rule | Method | Description |
|------|--------|-------------|
| R1 | `validate_repeating_companies()` | No duplicate companies |
| R2 | `validate_accounts()` | Account belongs to correct company |
| R3 | `validate_pos_mode_of_payment()` | Can't disable if in POS |

#### 1.4 Identify External Dependencies

```python
# These need interfaces in Go:
frappe.get_cached_value("Account", entry.default_account, "company")  # → AccountLookup
frappe.db.sql("SELECT ... FROM tabSales Invoice Payment ...")         # → POSChecker
```

---

### Step 2: Create Go Models

**Goal:** Define data structures matching the Python schema.

#### 2.1 Create Package Structure

```bash
mkdir -p modeofpayment
touch modeofpayment/model.go
touch modeofpayment/validation.go
touch modeofpayment/validation_test.go
```

#### 2.2 Define Structs (model.go)

```go
// Package modeofpayment implements the Mode of Payment doctype from ERPNext.
// Migrated from: erpnext/accounts/doctype/mode_of_payment/mode_of_payment.py
package modeofpayment

// PaymentType represents the type of payment method.
// Maps to: type DF.Literal["Cash", "Bank", "General", "Phone"]
type PaymentType string

const (
    Cash    PaymentType = "Cash"
    Bank    PaymentType = "Bank"
    General PaymentType = "General"
    Phone   PaymentType = "Phone"
)

// ModeOfPaymentAccount represents a child table row.
// Maps to: Mode of Payment Account doctype
type ModeOfPaymentAccount struct {
    Company        string // Link to Company
    DefaultAccount string // Link to Account
}

// ModeOfPayment represents a payment method master record.
// Maps to: Mode of Payment doctype
type ModeOfPayment struct {
    Name     string                 // Primary key (mode_of_payment field)
    Type     PaymentType            // Cash, Bank, General, Phone
    Enabled  bool                   // Active status
    Accounts []ModeOfPaymentAccount // Child table
}
```

#### 2.3 Define Interfaces (model.go)

```go
// AccountLookup abstracts database queries for account information.
// Production: queries Account doctype
// Testing: returns mock data
type AccountLookup interface {
    GetAccountCompany(accountName string) (string, error)
}

// POSChecker abstracts database queries for POS profile information.
// Production: queries Sales Invoice Payment table
// Testing: returns mock data
type POSChecker interface {
    GetPOSProfilesUsingMode(modeName string) ([]string, error)
}
```

---

### Step 3: Implement Validation

**Goal:** Port Python business rules to Go.

#### 3.1 Define Errors (validation.go)

```go
package modeofpayment

import (
    "errors"
    "fmt"
)

// Sentinel errors matching ERPNext's frappe.throw() messages
var (
    ErrDuplicateCompany = errors.New("same company is entered more than once")
    ErrAccountMismatch  = errors.New("account does not match with company")
    ErrModeInUse        = errors.New("mode of payment is used in POS profiles")
)

// ValidationError wraps sentinel errors with additional context
type ValidationError struct {
    Err     error
    Details string
}

func (e *ValidationError) Error() string {
    if e.Details != "" {
        return fmt.Sprintf("%s: %s", e.Err.Error(), e.Details)
    }
    return e.Err.Error()
}

func (e *ValidationError) Unwrap() error {
    return e.Err
}
```

#### 3.2 Implement Rules (validation.go)

**Rule 1: No Duplicate Companies**

```go
// ValidateRepeatingCompanies checks that no company appears multiple times.
//
// Python equivalent:
//   def validate_repeating_companies(self):
//       accounts_list = [entry.company for entry in self.accounts]
//       if len(accounts_list) != len(set(accounts_list)):
//           frappe.throw(_("Same Company is entered more than once"))
func (m *ModeOfPayment) ValidateRepeatingCompanies() error {
    seen := make(map[string]bool)
    for _, account := range m.Accounts {
        if seen[account.Company] {
            return &ValidationError{
                Err:     ErrDuplicateCompany,
                Details: fmt.Sprintf("company '%s' appears multiple times", account.Company),
            }
        }
        seen[account.Company] = true
    }
    return nil
}
```

**Rule 2: Account-Company Match**

```go
// ValidateAccounts verifies that each account's parent company matches.
//
// Python equivalent:
//   def validate_accounts(self):
//       for entry in self.accounts:
//           if frappe.get_cached_value("Account", entry.default_account, "company") != entry.company:
//               frappe.throw(_("Account {0} does not match..."))
func (m *ModeOfPayment) ValidateAccounts(lookup AccountLookup) error {
    for _, account := range m.Accounts {
        if account.DefaultAccount == "" {
            continue
        }

        accountCompany, err := lookup.GetAccountCompany(account.DefaultAccount)
        if err != nil {
            return fmt.Errorf("failed to lookup account %s: %w", account.DefaultAccount, err)
        }

        if accountCompany != account.Company {
            return &ValidationError{
                Err: ErrAccountMismatch,
                Details: fmt.Sprintf("account '%s' belongs to '%s', not '%s'",
                    account.DefaultAccount, accountCompany, account.Company),
            }
        }
    }
    return nil
}
```

**Rule 3: POS Profile Check**

```go
// ValidatePOSModeOfPayment prevents disabling a mode used in POS profiles.
//
// Python equivalent:
//   def validate_pos_mode_of_payment(self):
//       if not self.enabled:
//           pos_profiles = frappe.db.sql("SELECT ... WHERE mode_of_payment = %s")
//           if pos_profiles:
//               frappe.throw(_("POS Profile {} contains Mode of Payment {}..."))
func (m *ModeOfPayment) ValidatePOSModeOfPayment(checker POSChecker) error {
    if m.Enabled {
        return nil
    }

    profiles, err := checker.GetPOSProfilesUsingMode(m.Name)
    if err != nil {
        return fmt.Errorf("failed to check POS profiles: %w", err)
    }

    if len(profiles) > 0 {
        return &ValidationError{
            Err:     ErrModeInUse,
            Details: fmt.Sprintf("POS Profiles using this mode: %v", profiles),
        }
    }
    return nil
}
```

**Orchestrator**

```go
// Validate runs all validation checks.
// Matches ERPNext's validate() method.
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

### Step 4: Write Tests

**Goal:** Prove Go implementation matches Python behavior.

#### 4.1 Create Mock Implementations

```go
// validation_test.go
package modeofpayment

import "errors"

// mockAccountLookup simulates database queries
type mockAccountLookup struct {
    accounts map[string]string // account name → company name
}

func (m *mockAccountLookup) GetAccountCompany(name string) (string, error) {
    company, ok := m.accounts[name]
    if !ok {
        return "", errors.New("account not found")
    }
    return company, nil
}

// mockPOSChecker simulates POS profile queries
type mockPOSChecker struct {
    profilesByMode map[string][]string // mode name → profile names
}

func (m *mockPOSChecker) GetPOSProfilesUsingMode(name string) ([]string, error) {
    profiles := m.profilesByMode[name]
    return profiles, nil
}
```

#### 4.2 Write Table-Driven Tests

```go
func TestValidateRepeatingCompanies(t *testing.T) {
    tests := []struct {
        name     string
        accounts []ModeOfPaymentAccount
        wantErr  error
    }{
        {
            name:     "empty accounts - valid",
            accounts: []ModeOfPaymentAccount{},
            wantErr:  nil,
        },
        {
            name: "unique companies - valid",
            accounts: []ModeOfPaymentAccount{
                {Company: "Company A"},
                {Company: "Company B"},
            },
            wantErr: nil,
        },
        {
            name: "duplicate companies - error",
            accounts: []ModeOfPaymentAccount{
                {Company: "Company A"},
                {Company: "Company A"},
            },
            wantErr: ErrDuplicateCompany,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &ModeOfPayment{Accounts: tt.accounts}
            err := m.ValidateRepeatingCompanies()

            if tt.wantErr == nil {
                if err != nil {
                    t.Errorf("expected no error, got: %v", err)
                }
            } else {
                if !errors.Is(err, tt.wantErr) {
                    t.Errorf("expected error %v, got: %v", tt.wantErr, err)
                }
            }
        })
    }
}
```

#### 4.3 Run Tests

```bash
# Run with verbose output
go test -v ./modeofpayment/

# Check coverage
go test -cover ./modeofpayment/

# Generate coverage report
go test -coverprofile=coverage.out ./modeofpayment/
go tool cover -html=coverage.out -o coverage.html
```

---

## Code Patterns

### Pattern: Field Mapping

| Python (Frappe) | Go |
|-----------------|-----|
| `DF.Data` | `string` |
| `DF.Int` | `int` |
| `DF.Float` | `float64` |
| `DF.Check` | `bool` |
| `DF.Select` | Custom type (enum) |
| `DF.Link` | `string` (stores name/ID) |
| `DF.Table[T]` | `[]T` |
| `DF.Date` | `time.Time` |
| `DF.Datetime` | `time.Time` |
| `DF.Currency` | `decimal.Decimal` |

### Pattern: Error Translation

| Python | Go |
|--------|-----|
| `frappe.throw(message)` | `return &ValidationError{Err: ..., Details: ...}` |
| `frappe.throw(msg, title=...)` | Include title in Details |
| `frappe.ValidationError` | `ErrValidation` sentinel |

### Pattern: Database Calls

| Python | Go Interface |
|--------|--------------|
| `frappe.get_value("DocType", name, field)` | `interface { GetField(name) (value, error) }` |
| `frappe.get_cached_value(...)` | Same interface (caching is implementation detail) |
| `frappe.db.sql(query, values)` | Custom interface method for the query |
| `frappe.get_doc("DocType", name)` | `interface { Get(name) (*Entity, error) }` |

---

## Testing Strategy

### Test Pyramid

```
                    ┌─────────────────┐
                    │   End-to-End    │  Few (manual/automated)
                    │    Tests        │  - Full system with ERPNext
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   Integration   │  Some (with real DB)
                    │     Tests       │  - Repository tests
                    └────────┬────────┘  - API tests
                             │
           ┌─────────────────▼─────────────────┐
           │           Unit Tests              │  Many (fast, isolated)
           │  - Validation logic               │  - Run on every commit
           │  - Business rules                 │  - <1 second total
           └───────────────────────────────────┘
```

### Coverage Targets

| Layer | Target | Rationale |
|-------|--------|-----------|
| Domain (validation) | 90%+ | Core business logic |
| Application (use cases) | 80%+ | Orchestration logic |
| Infrastructure (repositories) | 70%+ | Integration points |
| HTTP handlers | 60%+ | Mostly delegation |

---

## Deployment Guide

### Shadow Mode Deployment

```yaml
# Feature flag configuration
features:
  mode_of_payment_go:
    enabled: true
    shadow_mode: true      # Run both, compare results
    rollout_percentage: 0  # 0% = shadow only
```

### Gradual Rollout

```
Week 1: Shadow Mode (0% traffic to Go)
├── Deploy Go service
├── Mirror all requests to Go
├── Log differences
└── Fix any discrepancies

Week 2: Canary (1% traffic to Go)
├── Route 1% of traffic to Go
├── Monitor error rates
├── Compare response times
└── Validate data consistency

Week 3: Ramp Up (10% → 50% → 100%)
├── Gradually increase traffic
├── Monitor at each step
└── Rollback if issues

Week 4: Cleanup
├── Disable Python endpoint
├── Remove feature flag
└── Archive Python code
```

---

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| Test fails with "nil pointer" | Uninitialized slice | Use `[]T{}` not `nil` |
| Error types don't match | Missing `Unwrap()` | Implement `Unwrap()` method |
| Coverage drops | Untested error paths | Add error case tests |
| Integration test fails | Mock data out of date | Sync with ERPNext schema |

### Debug Commands

```bash
# Verbose test output
go test -v ./...

# Test single function
go test -v -run TestValidateRepeatingCompanies ./modeofpayment/

# Check for race conditions
go test -race ./...

# Profile memory usage
go test -memprofile=mem.out ./...
go tool pprof mem.out
```

---

## References

- [Testing in Go](https://golang.org/pkg/testing/) — Go Documentation
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests) — Go Wiki
- [Effective Go](https://golang.org/doc/effective_go) — Go Team
