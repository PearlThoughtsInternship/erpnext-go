# Architecture

> System architecture for ERPNext legacy modernization to Go

---

## Table of Contents

- [Overview](#overview)
- [System Context](#system-context)
- [Container Diagram](#container-diagram)
- [Component Diagram](#component-diagram)
- [Data Flow](#data-flow)
- [Domain Model](#domain-model)
- [Layer Architecture](#layer-architecture)
- [Integration Patterns](#integration-patterns)

---

## Overview

This document describes the architecture of the Go-based modernization layer that will progressively replace ERPNext's Python/Frappe backend using the Strangler Fig pattern.

### Architectural Principles

| Principle | Description |
|-----------|-------------|
| 🎯 **Domain-Driven** | Business logic isolated from infrastructure |
| 🔌 **Port/Adapter** | Interfaces define boundaries, implementations are swappable |
| 🧪 **Test-First** | Every module has comprehensive tests before deployment |
| 📦 **Module-Per-DocType** | Each ERPNext DocType becomes a Go package |
| 🔄 **Incremental** | Migrate one bounded context at a time |

---

## System Context

### C4 Context Diagram

```mermaid
C4Context
    title System Context - ERPNext Modernization

    Person(user, "Business User", "Accountant, Sales Rep, HR Manager")
    Person(admin, "System Admin", "IT Administrator")

    System_Boundary(erpnext, "ERPNext Ecosystem") {
        System(legacy, "ERPNext Legacy", "Python/Frappe monolith")
        System(modern, "Go Services", "Modernized bounded contexts")
    }

    System_Ext(bank, "Banking APIs", "Payment gateways")
    System_Ext(tax, "Tax Authorities", "GST/VAT systems")

    Rel(user, legacy, "Uses", "HTTP/REST")
    Rel(user, modern, "Uses", "HTTP/REST")
    Rel(admin, legacy, "Manages")
    Rel(admin, modern, "Manages")
    Rel(legacy, bank, "Integrates")
    Rel(modern, tax, "Reports to")
```

### Current State (ERPNext Monolith)

```mermaid
flowchart TB
    subgraph users["👥 USERS"]
        browser["🌐 Web Browser"]
        mobile["📱 Mobile App"]
        api["🔌 API Clients"]
    end

    users --> erpnext

    subgraph erpnext["🐍 ERPNext (Frappe Framework)"]
        subgraph modules["Business Modules"]
            accounts["📊 Accounts"]
            stock["📦 Stock"]
            selling["💰 Selling"]
            hr["👥 HR"]
        end

        accounts <--> stock
        stock <--> selling
        accounts <--> selling
        selling <--> hr

        modules --> orm

        subgraph orm["Frappe ORM Layer"]
            document["Document Class"]
            doctype["DocType System"]
        end
    end

    orm --> storage

    subgraph storage["💾 Storage Layer"]
        mariadb[("MariaDB<br/>Database")]
        redis[("Redis<br/>Cache")]
        files[("File<br/>Storage")]
    end

    style erpnext fill:#306998,color:#fff
    style modules fill:#4B8BBE,color:#fff
    style orm fill:#FFD43B,color:#000
    style storage fill:#f8f9fa
```

### Target State (Hybrid with Go Services)

```mermaid
flowchart TB
    users["👥 Users"] --> gateway

    subgraph gateway["🔀 API Gateway / Router"]
        router["Feature Flag Router"]
    end

    gateway --> legacy
    gateway --> modern

    subgraph legacy["🐍 LEGACY (Python)"]
        py_accounts["Accounts Module"]
        py_stock["Stock Module"]
        py_selling["Selling Module"]
        py_hr["HR Module"]
    end

    subgraph modern["🔷 MODERN (Go)"]
        go_mop["✅ Mode of Payment"]
        go_tax["✅ Tax Calculator"]
        go_ledger["🔄 GL Engine"]
        go_bank["📋 Bank (planned)"]
    end

    legacy <-->|"Migration<br/>in progress"| modern

    legacy --> mariadb[("MariaDB")]
    modern --> postgres[("PostgreSQL")]

    mariadb <-.->|"Data Sync"| postgres

    style legacy fill:#306998,color:#fff
    style modern fill:#00ADD8,color:#fff
    style gateway fill:#6c757d,color:#fff
```

---

## Container Diagram

### Go Service Architecture

```mermaid
flowchart TB
    subgraph goapp["🔷 Go Application"]
        subgraph http["🌐 HTTP Layer"]
            router["Router<br/>(chi/mux)"]
            middleware["Middleware<br/>(auth, log)"]
            handlers["Handlers<br/>(REST)"]
        end

        http --> app

        subgraph app["⚙️ Application Layer"]
            commands["Commands<br/>(use cases)"]
            queries["Queries<br/>(read ops)"]
            events["Events<br/>(pub/sub)"]
        end

        app --> domain

        subgraph domain["💎 Domain Layer"]
            mop["modeofpayment<br/>├─ model.go<br/>├─ validation.go<br/>└─ repository"]
            ledger["ledger<br/>├─ model.go<br/>├─ engine.go<br/>└─ ports.go"]
            taxcalc["taxcalc<br/>├─ calculator.go<br/>└─ model.go"]
        end

        domain --> infra

        subgraph infra["🔧 Infrastructure Layer"]
            postgres["PostgreSQL<br/>Repository"]
            redis["Redis<br/>Cache"]
            legacy_bridge["Legacy<br/>Bridge"]
        end
    end

    infra --> external

    subgraph external["External Systems"]
        db[("PostgreSQL")]
        cache[("Redis")]
        erpnext["ERPNext API"]
    end

    style http fill:#cce5ff,stroke:#004085
    style app fill:#d4edda,stroke:#155724
    style domain fill:#fff3cd,stroke:#856404
    style infra fill:#e2e3e5,stroke:#383d41
```

---

## Component Diagram

### Mode of Payment Package

```mermaid
classDiagram
    class ModeOfPayment {
        +Name string
        +Type PaymentType
        +Enabled bool
        +Accounts []ModeOfPaymentAccount
        +Validate(AccountLookup, POSChecker) error
        +ValidateRepeatingCompanies() error
        +ValidateAccounts(AccountLookup) error
        +ValidatePOSModeOfPayment(POSChecker) error
    }

    class ModeOfPaymentAccount {
        +Company string
        +DefaultAccount string
    }

    class PaymentType {
        <<enumeration>>
        Cash
        Bank
        General
        Phone
    }

    class AccountLookup {
        <<interface>>
        +GetAccountCompany(name string) string, error
    }

    class POSChecker {
        <<interface>>
        +GetPOSProfilesUsingMode(name string) []string, error
    }

    class ValidationError {
        +Err error
        +Details string
        +Error() string
        +Unwrap() error
    }

    ModeOfPayment "1" *-- "*" ModeOfPaymentAccount : contains
    ModeOfPayment --> PaymentType : has type
    ModeOfPayment ..> AccountLookup : uses
    ModeOfPayment ..> POSChecker : uses
    ModeOfPayment ..> ValidationError : returns

    note for AccountLookup "Port interface for<br/>external dependencies"
    note for POSChecker "Abstracts Frappe<br/>database queries"
```

### GL Entry Engine Components

```mermaid
classDiagram
    class GLEntry {
        +Name string
        +PostingDate time.Time
        +Account string
        +Debit float64
        +Credit float64
        +Party string
        +PartyType string
        +VoucherType string
        +VoucherNo string
        +Company string
        +IsOpening string
        +IsCancelled int
        ...35 fields total
    }

    class Engine {
        +Accounts AccountLookup
        +Company CompanySettings
        +Budget BudgetValidator
        +GLStore GLEntryStore
        +MakeGLEntries([]GLEntry, PostingOptions) error
        +ProcessGLMap([]GLEntry, bool, bool) []GLEntry
        +validateDisabledAccounts([]GLEntry) error
    }

    class PostingOptions {
        +MergeEntries bool
        +FromRepost bool
        +AdditionalConditions string
    }

    class AccountLookup {
        <<interface>>
        +GetAccount(name string) *Account, error
        +IsDisabled(name string) bool, error
    }

    class GLEntryStore {
        <<interface>>
        +Save(*GLEntry) error
        +SaveBatch([]GLEntry) error
        +GetByVoucher(type, no string) []GLEntry, error
    }

    Engine --> AccountLookup : uses
    Engine --> GLEntryStore : uses
    Engine --> GLEntry : processes
    Engine --> PostingOptions : configured by
```

---

## Data Flow

### GL Entry Posting Flow

```mermaid
sequenceDiagram
    autonumber
    participant Client
    participant Handler as HTTP Handler
    participant Engine as GL Engine
    participant Validator as Validators
    participant Store as GLEntryStore
    participant DB as Database

    Client->>Handler: POST /gl-entries
    Handler->>Engine: MakeGLEntries(glMap, opts)

    activate Engine

    Engine->>Validator: validateDisabledAccounts()
    Validator-->>Engine: ✓ OK

    Engine->>Validator: validateAccountingPeriod()
    Validator-->>Engine: ✓ OK

    Engine->>Engine: ProcessGLMap()
    Note over Engine: Merge similar entries<br/>Toggle negative amounts

    Engine->>Validator: validateDebitCreditBalance()
    Validator-->>Engine: ✓ Balanced

    Engine->>Store: SaveBatch(entries)
    Store->>DB: INSERT INTO gl_entry
    DB-->>Store: ✓ Committed
    Store-->>Engine: ✓ Saved

    deactivate Engine

    Engine-->>Handler: nil (success)
    Handler-->>Client: 201 Created
```

### Shadow Mode Comparison Flow

```mermaid
sequenceDiagram
    autonumber
    participant Client
    participant Gateway as API Gateway
    participant Python as ERPNext (Python)
    participant Go as Go Service
    participant Comparator
    participant Logger

    Client->>Gateway: Request

    par Dual Execution
        Gateway->>Python: Forward Request
        Python-->>Gateway: Python Response
    and
        Gateway->>Go: Shadow Request
        Go-->>Comparator: Go Response
    end

    Gateway-->>Client: Python Response (Primary)

    Comparator->>Comparator: Compare Responses
    alt Responses Match
        Comparator->>Logger: Log: MATCH ✓
    else Responses Differ
        Comparator->>Logger: Log: DIFF ⚠️
        Note over Logger: Field-by-field<br/>difference report
    end
```

---

## Domain Model

### ERPNext Accounts Domain (Core Entities)

```mermaid
erDiagram
    COMPANY ||--o{ ACCOUNT : has
    COMPANY ||--o{ FISCAL_YEAR : defines
    COMPANY ||--o{ GL_ENTRY : records

    ACCOUNT ||--o{ GL_ENTRY : posts_to
    ACCOUNT }o--|| ACCOUNT : parent_of

    GL_ENTRY }o--|| VOUCHER : references
    GL_ENTRY }o--o| PARTY : involves

    JOURNAL_ENTRY ||--|{ GL_ENTRY : creates
    SALES_INVOICE ||--|{ GL_ENTRY : creates
    PAYMENT_ENTRY ||--|{ GL_ENTRY : creates

    COMPANY {
        string name PK
        string default_currency
        string chart_of_accounts
    }

    ACCOUNT {
        string name PK
        string account_type
        string company FK
        string parent_account
        bool is_group
        bool disabled
    }

    GL_ENTRY {
        string name PK
        date posting_date
        string account FK
        decimal debit
        decimal credit
        string voucher_type
        string voucher_no
        string company FK
    }

    VOUCHER {
        string doctype
        string name
        date posting_date
    }

    PARTY {
        string party_type
        string party_name
    }
```

### GL Entry State Machine

```mermaid
stateDiagram-v2
    [*] --> Draft: Create Entry

    Draft --> Validated: Validate()
    Validated --> Draft: Validation Failed

    Validated --> Posted: Submit()
    Posted --> Cancelled: Cancel()

    Cancelled --> [*]: Archived

    state Posted {
        [*] --> InLedger
        InLedger --> Reconciled: Bank Match
        InLedger --> Reversed: Reversal Entry
    }

    note right of Draft: Entry created but<br/>not yet validated
    note right of Validated: All validations pass<br/>Ready to post
    note right of Posted: Committed to GL<br/>Affects balances
    note right of Cancelled: Reversal entry<br/>created
```

---

## Layer Architecture

### Clean Architecture Layers

```mermaid
flowchart TB
    subgraph external["🌐 External World"]
        web["Web UI"]
        api["API Clients"]
        cli["CLI Tools"]
    end

    external --> http

    subgraph http["HTTP/API Layer"]
        direction LR
        handlers["Handlers"]
        routes["Routes"]
    end

    http --> application

    subgraph application["Application Layer"]
        direction LR
        usecases["Use Cases"]
        commands["Commands"]
        queries["Queries"]
    end

    application --> domain

    subgraph domain["💎 Domain Layer (CORE)"]
        direction TB
        entities["Entities"]
        valueobjects["Value Objects"]
        validation["Business Rules"]
        interfaces["Port Interfaces"]
    end

    domain -.->|"defines"| infra

    subgraph infra["Infrastructure Layer"]
        direction LR
        postgres["PostgreSQL Repo"]
        redis["Redis Cache"]
        bridge["Legacy Bridge"]
    end

    infra --> systems

    subgraph systems["External Systems"]
        db[("Database")]
        cache[("Cache")]
        erpnext["ERPNext"]
    end

    style domain fill:#fff3cd,stroke:#856404,stroke-width:3px
    style external fill:#e2e3e5
    style http fill:#cce5ff
    style application fill:#d4edda
    style infra fill:#f8d7da
```

### Dependency Rule

```mermaid
flowchart LR
    subgraph rule["Dependency Rule"]
        direction TB
        outer["Outer Layers"]
        inner["Inner Layers"]
        outer -->|"depend on"| inner
        inner -.-x|"NEVER depend on"| outer
    end

    subgraph allowed["✅ Allowed"]
        h1["HTTP → Application"]
        h2["Application → Domain"]
        h3["Infrastructure → Domain"]
    end

    subgraph forbidden["❌ Forbidden"]
        f1["Domain → HTTP"]
        f2["Domain → Infrastructure"]
        f3["Domain → Database"]
    end

    style allowed fill:#d4edda,stroke:#155724
    style forbidden fill:#f8d7da,stroke:#721c24
```

---

## Integration Patterns

### Anti-Corruption Layer

```mermaid
flowchart LR
    subgraph go["Go Service"]
        domain["Domain Layer"]
        acl["Anti-Corruption<br/>Layer"]
    end

    subgraph legacy["ERPNext Legacy"]
        frappe["Frappe API"]
        schemas["DocType Schemas"]
    end

    domain <--> acl
    acl <-->|"Translate"| frappe

    note1["ACL translates:<br/>• Frappe concepts → Domain concepts<br/>• Snake_case → PascalCase<br/>• frappe.throw → typed errors"]

    acl --- note1

    style acl fill:#fff3cd,stroke:#856404
    style note1 fill:#f8f9fa,stroke:#dee2e6
```

### Event-Driven Sync

```mermaid
flowchart LR
    subgraph go["Go Service"]
        goservice["Go Service"]
        publisher["Event Publisher"]
    end

    subgraph mq["Message Queue"]
        rabbitmq["RabbitMQ"]
    end

    subgraph python["Python Listener"]
        listener["Event Listener"]
        erpnext["ERPNext"]
    end

    goservice --> publisher
    publisher --> rabbitmq
    rabbitmq --> listener
    listener --> erpnext

    subgraph sync["Data Sync"]
        postgres[("PostgreSQL<br/>(Go DB)")]
        mariadb[("MariaDB<br/>(ERPNext DB)")]
    end

    goservice --> postgres
    erpnext --> mariadb
    postgres <-.->|"Bi-directional<br/>Sync"| mariadb

    style mq fill:#ff6b6b,color:#fff
```

### Feature Flag Routing

```mermaid
flowchart TB
    request["Incoming Request"] --> router

    router{"Feature Flag<br/>Enabled?"}

    router -->|"Yes"| gohandler["Go Handler"]
    router -->|"No"| legacyproxy["Legacy Proxy"]

    gohandler --> goresponse["Go Response"]
    legacyproxy --> legacyresponse["Legacy Response"]

    subgraph config["Feature Flag Config"]
        flag["mode_of_payment_go:<br/>enabled: true<br/>rollout: 50%"]
    end

    router -.-> config

    style router fill:#6c757d,color:#fff
    style gohandler fill:#00ADD8,color:#fff
    style legacyproxy fill:#306998,color:#fff
```

---

## Security Architecture

### Authentication Flow

```mermaid
sequenceDiagram
    autonumber
    participant User
    participant Gateway as API Gateway
    participant Auth as Auth Proxy
    participant OAuth as ERPNext OAuth
    participant Go as Go Service

    User->>Gateway: Request + Credentials
    Gateway->>Auth: Validate Session
    Auth->>OAuth: Verify with ERPNext

    alt Valid Session
        OAuth-->>Auth: User Info + Roles
        Auth-->>Gateway: JWT Token
        Gateway->>Go: Request + JWT
        Go->>Go: Validate JWT
        Go-->>Gateway: Response
        Gateway-->>User: Response
    else Invalid Session
        OAuth-->>Auth: 401 Unauthorized
        Auth-->>Gateway: Redirect to Login
        Gateway-->>User: 401 + Login URL
    end
```

### Authorization Model

```mermaid
classDiagram
    class Permission {
        +Role string
        +DocType string
        +Operations []string
        +MatchLevel int
    }

    class Role {
        +Name string
        +Permissions []Permission
    }

    class User {
        +ID string
        +Roles []Role
        +HasPermission(doctype, op) bool
    }

    User "1" --> "*" Role : has
    Role "1" --> "*" Permission : grants

    note for Permission "Mirrors ERPNext's<br/>permission model"
```

---

## Technology Decisions

| Component | Choice | Rationale |
|-----------|--------|-----------|
| Language | Go 1.21+ | Performance, single binary, strong typing |
| HTTP Router | chi or Fiber | Lightweight, middleware support |
| Database | PostgreSQL | Better Go ecosystem than MariaDB |
| Cache | Redis | Familiar from ERPNext stack |
| Testing | stdlib + testify | Simple, no framework lock-in |
| Config | envconfig | 12-factor app compliance |
| Logging | slog (stdlib) | Structured, performant |
| Metrics | Prometheus | Industry standard |

---

## References

- [C4 Model](https://c4model.com/) — Architecture diagram notation
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) — Robert C. Martin
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) — Alistair Cockburn
- [Domain-Driven Design](https://www.domainlanguage.com/ddd/) — Eric Evans
