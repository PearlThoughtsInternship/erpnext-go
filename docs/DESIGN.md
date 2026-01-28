# Design Decisions

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

### 1. Interface-Based Dependencies

**Decision:** All external dependencies accessed via interfaces defined in the domain layer.

**Context:**
- ERPNext uses global functions (`frappe.get_value()`, `frappe.db.sql()`)
- Global state makes unit testing impossible without full stack
- We need to test business logic in isolation

```mermaid
flowchart LR
    subgraph domain["💎 Domain Layer"]
        logic["Business Logic"]
        interface["Interface Definition"]
    end

    subgraph infra["🔧 Infrastructure"]
        postgres["PostgresAccountLookup"]
        mock["MockAccountLookup"]
    end

    logic --> interface
    interface -.->|"implemented by"| postgres
    interface -.->|"implemented by"| mock

    subgraph test["🧪 Tests"]
        unittest["Unit Tests"]
    end

    unittest --> mock

    style domain fill:#fff3cd,stroke:#856404
    style infra fill:#e2e3e5
    style test fill:#d4edda
```

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. Global functions (like Python) | Familiar to ERPNext devs | Untestable |
| B. Dependency injection via constructors | Testable, explicit | Verbose |
| **C. Interface-based injection** ✅ | Testable, swappable | Slight indirection |

**Consequences:**
- ✅ Unit tests run in milliseconds (no DB)
- ✅ Can swap PostgreSQL for SQLite in tests
- ✅ Legacy bridge is just another adapter
- ⚠️ More files/types to manage

---

### 2. Typed Sentinel Errors

**Decision:** Define error constants for programmatic error handling.

**Context:**
- Python uses `frappe.throw(message)` with string matching
- Callers can't reliably handle specific errors
- Go's `errors.Is()` enables type-safe error checks

```mermaid
flowchart TB
    subgraph python["🐍 Python (ERPNext)"]
        throw["frappe.throw(message)"]
        catch["try/except + string match"]
    end

    subgraph go["🔷 Go (Ours)"]
        sentinel["var ErrDuplicateCompany = errors.New(...)"]
        wrapper["ValidationError{Err, Details}"]
        check["errors.Is(err, ErrDuplicateCompany)"]
    end

    throw -->|"migrates to"| sentinel
    catch -->|"migrates to"| check
    sentinel --> wrapper
    wrapper --> check

    style python fill:#306998,color:#fff
    style go fill:#00ADD8,color:#fff
```

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. Return string errors | Simple | No type safety |
| B. Custom error types | Full control | Boilerplate |
| **C. Sentinel errors + wrapper** ✅ | Type-safe, details included | Moderate complexity |

**Consequences:**
- ✅ Callers can handle specific errors
- ✅ Details preserved for logging/display
- ✅ Works with `errors.Is()` and `errors.As()`
- ⚠️ Must remember to use Unwrap()

---

### 3. Table-Driven Tests

**Decision:** Use table-driven tests for comprehensive coverage.

```mermaid
flowchart LR
    subgraph table["Test Table"]
        case1["Case 1: empty accounts"]
        case2["Case 2: unique companies"]
        case3["Case 3: duplicate companies"]
        caseN["Case N: ..."]
    end

    subgraph runner["Test Runner"]
        loop["for _, tt := range tests"]
        subtest["t.Run(tt.name, ...)"]
        assert["Assert result"]
    end

    table --> loop
    loop --> subtest
    subtest --> assert

    style table fill:#cce5ff
    style runner fill:#d4edda
```

**Options Considered:**

| Option | Pros | Cons |
|--------|------|------|
| A. One function per test | Clear names | Verbose, hard to add cases |
| B. BDD framework (Ginkgo) | Expressive | External dependency |
| **C. Table-driven tests** ✅ | Go idiom, easy to extend | Slightly dense |

**Consequences:**
- ✅ Adding test cases is one line
- ✅ Subtests run in parallel if needed
- ✅ Clear failure messages with test names
- ⚠️ Table setup can be verbose for complex inputs

---

### 4. Domain Purity

**Decision:** Domain structs contain no infrastructure code.

```mermaid
flowchart TB
    subgraph bad["❌ Active Record (Frappe)"]
        doc["Document"]
        doc --> data["Data Fields"]
        doc --> validation["Validation"]
        doc --> persistence["Save/Load/Delete"]
        doc --> ui["UI Hooks"]
    end

    subgraph good["✅ Rich Domain Model (Go)"]
        subgraph domain["Domain"]
            entity["Entity"]
            entity --> data2["Data Fields"]
            entity --> rules["Business Rules"]
        end

        subgraph infra2["Infrastructure"]
            repo["Repository"]
            repo --> persistence2["Save/Load/Delete"]
        end

        entity -.-> repo
    end

    style bad fill:#f8d7da,stroke:#721c24
    style good fill:#d4edda,stroke:#155724
    style domain fill:#fff3cd
```

**Consequences:**
- ✅ Domain logic testable without database
- ✅ Same validation works in API, CLI, batch jobs
- ✅ Clear separation of concerns
- ⚠️ More layers (domain, application, infrastructure)

---

### 5. Incremental Migration

**Decision:** Migrate one bounded context at a time with feature flags.

```mermaid
gantt
    title Strangler Fig Migration Timeline
    dateFormat  YYYY-MM-DD
    section Mode of Payment
    Extract to Go       :done, mop1, 2026-01-01, 7d
    Shadow Mode         :done, mop2, after mop1, 3d
    1% Traffic          :done, mop3, after mop2, 3d
    100% Traffic        :done, mop4, after mop3, 2d
    section Tax Calculator
    Extract to Go       :done, tax1, 2026-01-10, 10d
    Shadow Mode         :done, tax2, after tax1, 3d
    Rollout             :done, tax3, after tax2, 5d
    section GL Engine
    Extract to Go       :active, gl1, 2026-01-20, 14d
    Shadow Mode         :gl2, after gl1, 5d
    Rollout             :gl3, after gl2, 7d
    section Account Master
    Extract             :acc1, after gl3, 10d
```

**Consequences:**
- ✅ Production never breaks
- ✅ Rollback is one config change
- ✅ Team learns patterns incrementally
- ⚠️ Two systems running in parallel temporarily

---

## Pattern Catalog

### Patterns Used

```mermaid
mindmap
  root((DDD Patterns))
    Tactical
      Entity
        ModeOfPayment
        GLEntry
      Value Object
        PaymentType
        PostingOptions
      Aggregate
        ModeOfPayment + Accounts
      Repository
        AccountLookup
        GLEntryStore
    Strategic
      Bounded Context
        Ledger Package
        TaxCalc Package
      Anti-Corruption Layer
        Legacy Bridge
      Context Map
        Shared Kernel
```

### Pattern Details

#### Repository Pattern

```mermaid
classDiagram
    class Repository~T~ {
        <<interface>>
        +Create(ctx, *T) error
        +Get(ctx, id string) *T, error
        +Update(ctx, *T) error
        +Delete(ctx, id string) error
        +List(ctx, ...Filter) []*T, error
    }

    class ModeOfPaymentRepository {
        <<interface>>
        +FindByType(ctx, PaymentType) []*ModeOfPayment, error
    }

    class PostgresMoPRepo {
        -db *sql.DB
        +Create(ctx, *ModeOfPayment) error
        +FindByType(ctx, PaymentType) []*ModeOfPayment, error
    }

    class MockMoPRepo {
        -data map[string]*ModeOfPayment
        +Create(ctx, *ModeOfPayment) error
        +FindByType(ctx, PaymentType) []*ModeOfPayment, error
    }

    Repository~T~ <|-- ModeOfPaymentRepository
    ModeOfPaymentRepository <|.. PostgresMoPRepo
    ModeOfPaymentRepository <|.. MockMoPRepo
```

#### Specification Pattern

```mermaid
flowchart LR
    subgraph specs["Specifications"]
        s1["ValidateAccounts"]
        s2["ValidateRepeatingCompanies"]
        s3["ValidatePOSModeOfPayment"]
    end

    subgraph composite["Composite Validate()"]
        v["Validate()"]
    end

    s1 --> v
    s2 --> v
    s3 --> v

    v -->|"all pass"| success["✓ nil"]
    v -->|"any fail"| error["✗ error"]

    style specs fill:#cce5ff
    style composite fill:#fff3cd
```

---

## Trade-off Analysis

### Complexity vs Testability

```mermaid
quadrantChart
    title Complexity vs Testability Trade-off
    x-axis Low Complexity --> High Complexity
    y-axis Low Testability --> High Testability
    quadrant-1 Ideal Zone
    quadrant-2 Over-engineered
    quadrant-3 Legacy Trap
    quadrant-4 Technical Debt

    Global Functions: [0.2, 0.15]
    Interface DI: [0.65, 0.85]
    Framework Magic: [0.8, 0.4]
    Simple Structs: [0.3, 0.5]
```

### Performance vs Maintainability

```mermaid
quadrantChart
    title Performance vs Maintainability Trade-off
    x-axis Low Maintainability --> High Maintainability
    y-axis Low Performance --> High Performance

    Hand-optimized C: [0.2, 0.9]
    Go with interfaces: [0.7, 0.75]
    Python Frappe: [0.85, 0.3]
    Java Enterprise: [0.6, 0.5]
```

### Migration Speed vs Risk

```mermaid
quadrantChart
    title Migration Speed vs Risk Trade-off
    x-axis High Risk --> Low Risk
    y-axis Slow Migration --> Fast Migration

    Big Bang Rewrite: [0.15, 0.85]
    Strangler Fig: [0.8, 0.4]
    Branch by Abstraction: [0.6, 0.5]
    Parallel Run Forever: [0.9, 0.1]
```

---

## Anti-Patterns to Avoid

### What Not to Do

```mermaid
flowchart TB
    subgraph antipatterns["❌ Anti-Patterns"]
        global["Global State"]
        stringerr["String Errors"]
        godobj["God Objects"]
        anemic["Anemic Domain"]
        premature["Premature Optimization"]
        copypaste["Copy-Paste ERPNext"]
    end

    subgraph solutions["✅ Solutions"]
        di["Dependency Injection"]
        typed["Typed Errors"]
        srp["Single Responsibility"]
        rich["Rich Domain Model"]
        profile["Profile First"]
        redesign["Redesign for Go"]
    end

    global -->|"replace with"| di
    stringerr -->|"replace with"| typed
    godobj -->|"replace with"| srp
    anemic -->|"replace with"| rich
    premature -->|"replace with"| profile
    copypaste -->|"replace with"| redesign

    style antipatterns fill:#f8d7da,stroke:#721c24
    style solutions fill:#d4edda,stroke:#155724
```

### Code Comparison

```mermaid
flowchart LR
    subgraph bad["❌ BAD"]
        badcode["var db *sql.DB<br/><br/>func GetAccount(name) {<br/>  return db.Query(...)<br/>}"]
    end

    subgraph good["✅ GOOD"]
        goodcode["type AccountService struct {<br/>  repo AccountRepository<br/>}<br/><br/>func (s *AccountService) GetAccount(ctx, name) {<br/>  return s.repo.Get(ctx, name)<br/>}"]
    end

    bad -->|"refactor to"| good

    style bad fill:#f8d7da
    style good fill:#d4edda
```

---

## Decision Records

### ADR-001: Use Go for Modernization

```mermaid
flowchart TB
    subgraph context["Context"]
        c1["ERPNext written in Python"]
        c2["Need modernization target"]
        c3["Team has mixed experience"]
    end

    subgraph options["Options Evaluated"]
        o1["Python (stay)"]
        o2["Go"]
        o3["Rust"]
        o4["Java/Kotlin"]
    end

    subgraph decision["Decision: Go"]
        d1["Strong typing"]
        d2["Fast compilation"]
        d3["Single binary deploy"]
        d4["Good concurrency"]
    end

    context --> options
    options --> decision
    o2 -->|"selected"| decision

    style decision fill:#d4edda,stroke:#155724
```

**Status:** Accepted
**Date:** 2026-01-27

### ADR-002: PostgreSQL over MariaDB

```mermaid
flowchart LR
    subgraph legacy["Legacy"]
        mariadb[("MariaDB")]
    end

    subgraph modern["Modern"]
        postgres[("PostgreSQL")]
    end

    subgraph rationale["Rationale"]
        r1["Better Go drivers"]
        r2["JSON support"]
        r3["Enum types"]
        r4["Advanced indexing"]
    end

    legacy -->|"migrate to"| modern
    rationale --> modern

    style legacy fill:#f8d7da
    style modern fill:#d4edda
```

**Status:** Accepted
**Date:** 2026-01-27

### ADR-003: Interface-Based Testing

```mermaid
flowchart TB
    subgraph problem["Problem"]
        p1["Need to test without DB"]
        p2["External dependencies everywhere"]
        p3["Frappe uses globals"]
    end

    subgraph solution["Solution"]
        s1["Define interfaces for deps"]
        s2["Inject implementations"]
        s3["Use mocks in tests"]
    end

    subgraph benefit["Benefits"]
        b1["Tests run in milliseconds"]
        b2["No test database needed"]
        b3["Isolated unit tests"]
    end

    problem --> solution
    solution --> benefit

    style solution fill:#d4edda
```

**Status:** Accepted
**Date:** 2026-01-27

---

## References

- [Domain-Driven Design](https://www.domainlanguage.com/ddd/) — Eric Evans
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) — Robert C. Martin
- [Effective Go](https://golang.org/doc/effective_go) — Go Team
- [Architecture Decision Records](https://adr.github.io/) — ADR GitHub
