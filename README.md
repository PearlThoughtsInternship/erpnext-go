# ERPNext Legacy Modernization — Python to Go

<p align="center">
  <img src="https://img.shields.io/badge/iteration-3%20in%20progress-blue" alt="Status">
  <img src="https://img.shields.io/badge/tests-68%20passing-brightgreen" alt="Tests">
  <img src="https://img.shields.io/badge/coverage-85%25+-green" alt="Coverage">
  <img src="https://img.shields.io/badge/go-1.21+-00ADD8" alt="Go Version">
</p>

A demonstration of modernizing ERPNext (Python/Frappe) to Go using the **Strangler Fig Pattern** with AI-assisted, iterative extraction.

---

## Quick Start

```bash
git clone git@github.com:PearlThoughtsInternship/erpnext-go.git
cd erpnext-go
go test -v ./...
```

---

## Current Progress

| Iteration | Module | Tests | Coverage | Status |
|-----------|--------|-------|----------|--------|
| 1 | Mode of Payment | 19 | 85.3% | ✅ Complete |
| 2 | Tax Calculator | 24 | 90.2% | ✅ Complete |
| 3 | **GL Entry Engine** | 25 | 49.1% | 🔄 In Progress |
| 4 | Account Master | — | — | 📋 Planned |
| 5 | Journal Entry | — | — | 📋 Planned |

**Total: 68 tests passing across 3 packages**

---

## "But The Accounts Module Has Dependencies!"

> *"I want to extract the Accounts module, but it depends on Stock, Selling, and other modules. Don't I need to migrate everything at once?"*

**No.** This is the core insight from Sam Newman's *Monolith to Microservices* and Michael Feathers' *Working Effectively with Legacy Code*:

### The Bounded Context Strategy

```
┌─────────────────────────────────────────────────────────────────────┐
│                     ERPNext Monolith (Python)                       │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐         │
│  │ Accounts │◄──│  Stock   │◄──│ Selling  │◄──│    HR    │         │
│  │          │──►│          │──►│          │──►│          │         │
│  └──────────┘   └──────────┘   └──────────┘   └──────────┘         │
│       ▲                                                             │
│       │ Dependencies flow everywhere                                │
└───────┼─────────────────────────────────────────────────────────────┘
        │
        │  Extract with INTERFACES at the boundary
        ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Go Bounded Context                             │
│                                                                     │
│   ┌──────────────────────────────────────────────────────────┐     │
│   │                    Ledger Package                         │     │
│   │                                                           │     │
│   │   ┌─────────────────┐     ┌─────────────────┐            │     │
│   │   │  GLEntry Model  │     │  MakeGLEntries  │            │     │
│   │   │  (pure Go)      │     │  (pure logic)   │            │     │
│   │   └─────────────────┘     └─────────────────┘            │     │
│   │              │                     │                      │     │
│   │              ▼                     ▼                      │     │
│   │   ┌─────────────────────────────────────────────┐        │     │
│   │   │         PORT INTERFACES                      │        │     │
│   │   │  AccountLookup    CompanySettings           │        │     │
│   │   │  GLEntryStore     BudgetValidator           │        │     │
│   │   └─────────────────────────────────────────────┘        │     │
│   └──────────────────────────────────────────────────────────┘     │
│                              │                                      │
│                              ▼                                      │
│   ┌──────────────────────────────────────────────────────────┐     │
│   │              ADAPTERS (Swappable)                         │     │
│   │                                                           │     │
│   │   🧪 Test: MockAccountLookup    🏭 Prod: PostgresAdapter │     │
│   │   🧪 Test: MockGLStore          🏭 Prod: MariaDBAdapter  │     │
│   └──────────────────────────────────────────────────────────┘     │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### How We Handle Dependencies

| Dependency Type | Strategy | Example |
|-----------------|----------|---------|
| **Data from other modules** | Interface + Mock | `AccountLookup.GetAccount()` returns mock data in tests |
| **Writes to other modules** | Interface + Stub | `GLEntryStore.Save()` validates behavior without DB |
| **Complex cross-module logic** | Anti-Corruption Layer | Translate Stock concepts to Accounts concepts |
| **Shared calculations** | Extract to shared package | `Flt()`, `Round()` utilities |

### The Test Double Hierarchy

From *xUnit Test Patterns* by Gerard Meszaros:

| Double | Purpose | We Use It For |
|--------|---------|---------------|
| **Mock** | Verify interactions | `GLEntryStore.Save()` was called correctly |
| **Stub** | Return canned answers | `AccountLookup.IsDisabled()` returns `false` |
| **Fake** | Working implementation | In-memory store for integration tests |
| **Spy** | Record calls for later | Verify GL entries posted in correct order |

### Real Example: GL Entry Engine

The GL Entry Engine depends on:
- Account master data → **`AccountLookup` interface**
- Company settings → **`CompanySettings` interface**
- Budget validation → **`BudgetValidator` interface**
- Payment ledger → **`PaymentLedgerStore` interface**

```go
// ports.go - Define what we NEED, not how to get it
type AccountLookup interface {
    GetAccount(name string) (*Account, error)
    IsDisabled(name string) (bool, error)
}

// engine.go - Business logic uses interfaces
func (e *Engine) MakeGLEntries(glMap []GLEntry, opts PostingOptions) error {
    // Validate disabled accounts - works with ANY implementation
    if err := e.validateDisabledAccounts(glMap); err != nil {
        return err
    }
    // ...
}

// engine_test.go - Tests use mocks
func TestValidateDisabledAccounts(t *testing.T) {
    engine := &Engine{
        Accounts: &mockAccountLookup{...}, // No real DB needed
    }
    // Test runs in milliseconds
}
```

### The Strangler Fig In Action

```
Phase 1: Shadow Mode
┌──────────────────────────────────────────────────────────────┐
│                        Request                                │
│                           │                                   │
│            ┌──────────────┴──────────────┐                   │
│            ▼                              ▼                   │
│     ┌──────────────┐              ┌──────────────┐           │
│     │   ERPNext    │              │     Go       │           │
│     │   (Python)   │              │   (Shadow)   │           │
│     └──────┬───────┘              └──────┬───────┘           │
│            │                              │                   │
│            ▼                              ▼                   │
│     ┌──────────────┐              ┌──────────────┐           │
│     │   Response   │   Compare    │   Response   │           │
│     │   (Primary)  │◄────────────►│   (Logged)   │           │
│     └──────────────┘              └──────────────┘           │
└──────────────────────────────────────────────────────────────┘

Phase 2: Traffic Switch
┌──────────────────────────────────────────────────────────────┐
│                        Request                                │
│                           │                                   │
│                           ▼                                   │
│                    ┌──────────────┐                          │
│                    │     Go       │ ◄── Primary now          │
│                    │   (Primary)  │                          │
│                    └──────────────┘                          │
│                                                               │
│     ERPNext (Python) still available for rollback            │
└──────────────────────────────────────────────────────────────┘
```

---

## Project Structure

```
erpnext-go/
├── modeofpayment/     # ✅ Iteration 1: Payment methods
├── taxcalc/           # ✅ Iteration 2: Tax calculations
├── ledger/            # 🔄 Iteration 3: GL Entry Engine
└── docs/
    ├── ARCHITECTURE.md      # System design + diagrams
    ├── DESIGN.md            # Design decisions
    ├── IMPLEMENTATION.md    # Step-by-step guide
    └── AI_ENGINEERING.md    # AI-assisted modernization
```

---

## Documentation

| Document | What You'll Learn |
|----------|-------------------|
| **[Architecture](docs/ARCHITECTURE.md)** | Hexagonal architecture, bounded contexts, DDD patterns |
| **[Design Decisions](docs/DESIGN.md)** | Why interfaces? Why typed errors? Trade-offs explained |
| **[Implementation Guide](docs/IMPLEMENTATION.md)** | Step-by-step migration process |
| **[AI Engineering](docs/AI_ENGINEERING.md)** | How AI accelerates legacy modernization |
| **[Parity Report](../PARITY_REPORT.md)** | Field-by-field Python → Go comparison |

---

## Key References

| Book | Author | Key Pattern |
|------|--------|-------------|
| [Monolith to Microservices](https://www.oreilly.com/library/view/monolith-to-microservices/9781492047834/) | Sam Newman | **Strangler Fig**, Branch by Abstraction |
| [Working Effectively with Legacy Code](https://www.oreilly.com/library/view/working-effectively-with/0131177052/) | Michael Feathers | **Seams**, Characterization Tests |
| [Domain-Driven Design](https://www.oreilly.com/library/view/domain-driven-design-tackling/0321125215/) | Eric Evans | **Bounded Contexts**, Anti-Corruption Layer |
| [Clean Architecture](https://www.oreilly.com/library/view/clean-architecture-a/9780134494272/) | Robert C. Martin | **Ports & Adapters**, Dependency Inversion |

---

## Why This Approach Works

| Concern | Solution |
|---------|----------|
| "Modules have dependencies" | Interfaces abstract dependencies; mocks provide test isolation |
| "Can't test without full stack" | Pure domain logic + injected dependencies = instant unit tests |
| "Migration takes forever" | Extract one bounded context at a time; value delivered incrementally |
| "How do we know it's correct?" | Shadow mode compares Python vs Go outputs before switching |
| "What if we need to rollback?" | Feature flags control routing; legacy remains operational |

---

## Contributing

1. Pick a module from the [iteration roadmap](docs/IMPLEMENTATION.md#iteration-roadmap)
2. Read the [AI Engineering guide](docs/AI_ENGINEERING.md) for workflow
3. Follow existing patterns in `modeofpayment/` and `taxcalc/`
4. Aim for 85%+ test coverage

---

<p align="center">
  <sub>Built with ❤️ for legacy modernization | <a href="docs/AI_ENGINEERING.md">AI-Assisted</a></sub>
</p>
