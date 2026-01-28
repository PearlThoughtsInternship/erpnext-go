# Implementation Guide

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

```mermaid
flowchart LR
    subgraph phase1["📖 Phase 1"]
        analyze["Analyze<br/>Python"]
    end

    subgraph phase2["🏗️ Phase 2"]
        model["Model<br/>Go"]
    end

    subgraph phase3["✅ Phase 3"]
        validate["Validate<br/>Go"]
    end

    subgraph phase4["🧪 Phase 4"]
        test["Test<br/>Parity"]
    end

    subgraph phase5["🚀 Phase 5"]
        deploy["Deploy<br/>Shadow"]
    end

    phase1 --> phase2 --> phase3 --> phase4 --> phase5

    analyze --> a1["Schema"]
    analyze --> a2["Rules"]
    analyze --> a3["Tests"]

    model --> m1["Structs"]
    model --> m2["Enums"]
    model --> m3["Types"]

    validate --> v1["Rules"]
    validate --> v2["Errors"]
    validate --> v3["Methods"]

    test --> t1["Unit"]
    test --> t2["Integration"]
    test --> t3["Parity"]

    deploy --> d1["Feature Flag"]
    deploy --> d2["Monitor"]

    style phase1 fill:#cce5ff
    style phase2 fill:#d4edda
    style phase3 fill:#fff3cd
    style phase4 fill:#d1ecf1
    style phase5 fill:#f8d7da
```

### Detailed Workflow

```mermaid
flowchart TB
    start([Start Migration]) --> analyze

    subgraph analyze["Phase 1: Analysis"]
        a1[Read Python Source] --> a2[Extract Schema from JSON]
        a2 --> a3[Document Business Rules]
        a3 --> a4[Identify Dependencies]
    end

    analyze --> model

    subgraph model["Phase 2: Modeling"]
        m1[Create Package Structure] --> m2[Define Go Structs]
        m2 --> m3[Define Type Enums]
        m3 --> m4[Define Interfaces]
    end

    model --> implement

    subgraph implement["Phase 3: Implementation"]
        i1[Define Error Types] --> i2[Implement Validations]
        i2 --> i3[Add Python Comments]
    end

    implement --> test

    subgraph test["Phase 4: Testing"]
        t1[Create Mock Implementations] --> t2[Write Table-Driven Tests]
        t2 --> t3[Verify Coverage > 80%]
    end

    test --> deploy

    subgraph deploy["Phase 5: Deployment"]
        d1[Shadow Mode] --> d2[Compare Outputs]
        d2 --> d3{Match?}
        d3 -->|Yes| d4[Enable Traffic]
        d3 -->|No| d5[Fix Differences]
        d5 --> d2
    end

    deploy --> finish([Migration Complete])

    style analyze fill:#cce5ff
    style model fill:#d4edda
    style implement fill:#fff3cd
    style test fill:#d1ecf1
    style deploy fill:#f5c6cb
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

```mermaid
flowchart TB
    subgraph erpnext["erpnext/"]
        subgraph accounts["accounts/"]
            subgraph doctype["doctype/"]
                subgraph mop["mode_of_payment/"]
                    json["mode_of_payment.json<br/>📋 Schema"]
                    py["mode_of_payment.py<br/>🐍 Business Logic"]
                    test["test_mode_of_payment.py<br/>🧪 Tests"]
                    js["mode_of_payment.js<br/>🎨 UI (skip)"]
                end
            end
        end
    end

    json -->|"extract fields"| gomodel["Go Structs"]
    py -->|"extract rules"| govalidation["Go Validation"]
    test -->|"migrate cases"| gotest["Go Tests"]

    style json fill:#cce5ff
    style py fill:#fff3cd
    style test fill:#d4edda
    style js fill:#e2e3e5
```

#### 1.2 Extract Schema from JSON

**Document the mapping:**

| JSON Field | Type | Go Type |
|------------|------|---------|
| `mode_of_payment` | Data (required, unique) | `string` |
| `type` | Select | `PaymentType` (enum) |
| `enabled` | Check | `bool` |
| `accounts` | Table | `[]ModeOfPaymentAccount` |

#### 1.3 Extract Business Rules from Python

```mermaid
flowchart LR
    subgraph python["Python validate()"]
        v1["validate_accounts()"]
        v2["validate_repeating_companies()"]
        v3["validate_pos_mode_of_payment()"]
    end

    subgraph rules["Business Rules"]
        r1["R1: Account belongs<br/>to correct company"]
        r2["R2: No duplicate<br/>companies"]
        r3["R3: Can't disable<br/>if in POS"]
    end

    v1 --> r1
    v2 --> r2
    v3 --> r3

    style python fill:#306998,color:#fff
    style rules fill:#fff3cd
```

#### 1.4 Identify External Dependencies

```mermaid
flowchart LR
    subgraph python["Python Calls"]
        f1["frappe.get_cached_value()"]
        f2["frappe.db.sql()"]
    end

    subgraph interfaces["Go Interfaces"]
        i1["AccountLookup"]
        i2["POSChecker"]
    end

    f1 -->|"becomes"| i1
    f2 -->|"becomes"| i2

    style python fill:#306998,color:#fff
    style interfaces fill:#00ADD8,color:#fff
```

---

### Step 2: Create Go Models

**Goal:** Define data structures matching the Python schema.

#### 2.1 Package Structure

```mermaid
flowchart TB
    subgraph pkg["modeofpayment/"]
        model["model.go<br/>Structs & Types"]
        validation["validation.go<br/>Business Rules"]
        ports["ports.go<br/>Interfaces"]
        test["validation_test.go<br/>Tests"]
    end

    model --> validation
    ports --> validation
    validation --> test

    style model fill:#cce5ff
    style validation fill:#fff3cd
    style ports fill:#d4edda
    style test fill:#d1ecf1
```

#### 2.2 Type Definitions

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
```

---

### Step 3: Implement Validation

**Goal:** Port Python business rules to Go.

#### 3.1 Validation Flow

```mermaid
sequenceDiagram
    autonumber
    participant Caller
    participant MoP as ModeOfPayment
    participant V1 as ValidateAccounts
    participant V2 as ValidateRepeatingCompanies
    participant V3 as ValidatePOSModeOfPayment
    participant Lookup as AccountLookup

    Caller->>MoP: Validate(lookup, checker)

    MoP->>V1: ValidateAccounts(lookup)
    V1->>Lookup: GetAccountCompany(account)
    Lookup-->>V1: company
    V1-->>MoP: ✓ OK

    MoP->>V2: ValidateRepeatingCompanies()
    V2-->>MoP: ✓ OK

    MoP->>V3: ValidatePOSModeOfPayment(checker)
    V3-->>MoP: ✓ OK

    MoP-->>Caller: nil (success)
```

#### 3.2 Error Hierarchy

```mermaid
flowchart TB
    subgraph errors["Error Types"]
        sentinel["Sentinel Errors"]
        wrapper["ValidationError"]
    end

    subgraph sentinels["Sentinels"]
        e1["ErrDuplicateCompany"]
        e2["ErrAccountMismatch"]
        e3["ErrModeInUse"]
    end

    sentinel --> sentinels

    wrapper --> err["Err (sentinel)"]
    wrapper --> details["Details (context)"]

    sentinels --> err

    usage["errors.Is(err, ErrDuplicateCompany)"]
    wrapper --> usage

    style errors fill:#f8d7da
    style sentinels fill:#fff3cd
    style usage fill:#d4edda
```

---

### Step 4: Write Tests

**Goal:** Prove Go implementation matches Python behavior.

#### 4.1 Test Structure

```mermaid
flowchart TB
    subgraph test["validation_test.go"]
        subgraph mocks["Mock Implementations"]
            m1["mockAccountLookup"]
            m2["mockPOSChecker"]
        end

        subgraph tests["Table-Driven Tests"]
            t1["TestValidateRepeatingCompanies"]
            t2["TestValidateAccounts"]
            t3["TestValidatePOSModeOfPayment"]
            t4["TestValidate (orchestrator)"]
        end
    end

    mocks --> tests

    style mocks fill:#cce5ff
    style tests fill:#d4edda
```

#### 4.2 Test Coverage Goals

```mermaid
pie title Test Coverage Targets
    "Domain Logic" : 90
    "Application" : 80
    "Infrastructure" : 70
    "HTTP Handlers" : 60
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

```mermaid
flowchart LR
    subgraph python["Python"]
        throw["frappe.throw(msg)"]
        throwTitle["frappe.throw(msg, title)"]
    end

    subgraph go["Go"]
        valErr["ValidationError{Err, Details}"]
    end

    throw --> valErr
    throwTitle --> valErr

    style python fill:#306998,color:#fff
    style go fill:#00ADD8,color:#fff
```

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

```mermaid
flowchart TB
    subgraph pyramid["Test Pyramid"]
        e2e["🔝 E2E Tests<br/>Few (manual/automated)<br/>Full system with ERPNext"]

        integration["📊 Integration Tests<br/>Some (with real DB)<br/>Repository & API tests"]

        unit["🧪 Unit Tests<br/>Many (fast, isolated)<br/>Run on every commit"]
    end

    e2e --> integration --> unit

    style e2e fill:#f8d7da,stroke:#721c24
    style integration fill:#fff3cd,stroke:#856404
    style unit fill:#d4edda,stroke:#155724
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

### Shadow Mode Architecture

```mermaid
flowchart TB
    request["📨 Request"] --> gateway["🔀 API Gateway"]

    gateway --> python["🐍 Python<br/>(Primary)"]
    gateway -.->|shadow| go["🔷 Go<br/>(Shadow)"]

    python --> response["📤 Response"]

    go --> comparator["📊 Comparator"]
    response --> comparator

    comparator --> decision{Match?}

    decision -->|Yes| log_match["✅ Log: MATCH"]
    decision -->|No| log_diff["⚠️ Log: DIFF"]

    subgraph config["Configuration"]
        ff["Feature Flag:<br/>shadow_mode: true<br/>rollout: 0%"]
    end

    gateway -.-> config

    style python fill:#306998,color:#fff
    style go fill:#00ADD8,color:#fff
    style comparator fill:#fff3cd
```

### Gradual Rollout Timeline

```mermaid
gantt
    title Deployment Phases
    dateFormat  YYYY-MM-DD
    section Shadow Mode
    Deploy Go service      :s1, 2026-01-28, 2d
    Mirror all requests    :s2, after s1, 3d
    Fix discrepancies      :s3, after s2, 2d

    section Canary (1%)
    Route 1% traffic       :c1, after s3, 3d
    Monitor error rates    :c2, after s3, 3d
    Validate consistency   :c3, after s3, 3d

    section Ramp Up
    10% traffic            :r1, after c1, 2d
    50% traffic            :r2, after r1, 2d
    100% traffic           :r3, after r2, 2d

    section Cleanup
    Disable Python         :x1, after r3, 1d
    Remove feature flag    :x2, after x1, 1d
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

### Debug Decision Tree

```mermaid
flowchart TB
    start["Test Failing?"] --> type{Error Type?}

    type -->|"nil pointer"| nil["Check slice/map<br/>initialization"]
    type -->|"errors.Is fails"| unwrap["Add Unwrap()<br/>method"]
    type -->|"wrong value"| value["Check Python<br/>behavior"]
    type -->|"timeout"| timeout["Check mock<br/>implementation"]

    nil --> fix1["Use []T{}<br/>not nil"]
    unwrap --> fix2["func (e *Err) Unwrap()<br/>{ return e.Err }"]
    value --> fix3["Run Python<br/>comparison"]
    timeout --> fix4["Add return<br/>to mock"]

    style start fill:#f8d7da
    style fix1 fill:#d4edda
    style fix2 fill:#d4edda
    style fix3 fill:#d4edda
    style fix4 fill:#d4edda
```

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

## Iteration Roadmap

```mermaid
timeline
    title ERPNext Go Migration Roadmap
    section Phase 1 (Complete)
        2026-01 : Mode of Payment
               : 19 tests, 85% coverage
    section Phase 2 (Complete)
        2026-01 : Tax Calculator
               : 24 tests, 90% coverage
    section Phase 3 (In Progress)
        2026-01 : GL Entry Engine
               : 25 tests, 49% coverage
    section Phase 4 (Planned)
        2026-02 : Account Master
               : Fiscal Year
    section Phase 5 (Planned)
        2026-02 : Journal Entry
               : Manual adjustments
```

---

## References

- [Testing in Go](https://golang.org/pkg/testing/) — Go Documentation
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests) — Go Wiki
- [Effective Go](https://golang.org/doc/effective_go) — Go Team
