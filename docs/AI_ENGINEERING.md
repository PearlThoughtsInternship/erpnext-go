# AI Engineering in Legacy Modernization

> Using AI assistants to accelerate and de-risk ERPNext to Go migration

---

## Table of Contents

- [Overview](#overview)
- [AI Engineering Phases](#ai-engineering-phases)
- [AI Capabilities Matrix](#ai-capabilities-matrix)
- [Human-AI Collaboration Model](#human-ai-collaboration-model)
- [Practical Workflows](#practical-workflows)
- [Lessons from ERPNext Migration](#lessons-from-erpnext-migration)
- [Prompting Strategies](#prompting-strategies)
- [Quality Gates](#quality-gates)

---

## Overview

### What is AI Engineering for Modernization?

AI Engineering applies AI assistants to accelerate legacy code modernization through:

| Capability | Description |
|------------|-------------|
| **Code Analysis** | AI extracts patterns, business rules, and dependencies from legacy code |
| **Design Assistance** | AI helps map legacy concepts to modern architecture |
| **Code Generation** | AI translates Python idioms to Go idioms |
| **Validation** | AI generates tests and identifies parity gaps |

### Traditional vs AI-Assisted Approach

| Traditional Approach | AI-Assisted Approach |
|---------------------|----------------------|
| Manual code reading (hours) | Automated pattern extraction (minutes) |
| Implicit business rules | Documented rule extraction |
| Sparse test coverage | Generated comprehensive test cases |
| Months of analysis | Days of analysis |
| Senior developer bottleneck | Democratized understanding |

### ROI of AI-Assisted Modernization

Based on the ERPNext migration experience:

| Metric | Traditional | AI-Assisted | Improvement |
|--------|-------------|-------------|-------------|
| Code comprehension | 4-8 hours/module | 30-60 min/module | 8x faster |
| Test case generation | 2-4 hours/module | 30 min/module | 4-8x faster |
| Documentation | Often skipped | Auto-generated | 100% coverage |
| Parity verification | Manual spot-checks | Systematic comparison | Higher confidence |

---

## AI Engineering Phases

### Phase Diagram

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Discovery  │ → │   Design    │ → │Implementation│ → │ Validation  │ → │   Cutover   │
│             │    │             │    │             │    │             │    │             │
│ • Analyze   │    │ • Interfaces│    │ • Translate │    │ • Parity    │    │ • Shadow    │
│ • Extract   │    │ • Structs   │    │ • Test      │    │ • Coverage  │    │ • Feature   │
│ • Map deps  │    │ • Errors    │    │ • Document  │    │ • Edge cases│    │   flags     │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### Phase 1: Discovery

**Goal:** Understand the legacy codebase without modifying it.

| Activity | AI Role | Human Role |
|----------|---------|------------|
| **Code Analysis** | Parse Python files, identify classes/methods | Validate AI's understanding |
| **Business Rule Extraction** | Find validation logic, identify invariants | Confirm rules match business intent |
| **Dependency Mapping** | Trace `frappe.get_value()` calls | Prioritize which dependencies to abstract |
| **Test Analysis** | Identify existing tests, coverage gaps | Decide test strategy |

**Deliverables:**
- Field mapping table (Python → Go)
- Business rules document
- External dependency list
- Test gap analysis

**Example Prompt:**
```
Read the Python file at:
erpnext/accounts/doctype/mode_of_payment/mode_of_payment.py

Extract all business rules in this format:

| Rule ID | Method | Description | Error Condition |
|---------|--------|-------------|-----------------|
| R1 | validate_repeating_companies | No duplicate companies | frappe.throw() |
```

### Phase 2: Design

**Goal:** Map legacy architecture to modern patterns.

| Activity | AI Role | Human Role |
|----------|---------|------------|
| **Interface Definition** | Suggest port interfaces based on external calls | Review for completeness |
| **Type Design** | Propose Go structs from JSON schemas | Adjust naming conventions |
| **Error Strategy** | Map `frappe.throw()` to typed errors | Decide error granularity |
| **Test Strategy** | Propose test case matrix | Approve test boundaries |

**Deliverables:**
- Go interface definitions
- Struct designs with comments
- Error type catalog
- Test case matrix

**Example Output:**
```go
// Port interface derived from frappe.get_cached_value() calls
type AccountLookup interface {
    GetAccountCompany(accountName string) (string, error)
}
```

### Phase 3: Implementation

**Goal:** Generate Go code with parity to Python.

| Activity | AI Role | Human Role |
|----------|---------|------------|
| **Code Translation** | Convert Python methods to Go | Review idiom translation |
| **Comment Generation** | Add Python equivalent comments | Verify accuracy |
| **Test Generation** | Create table-driven tests | Add edge cases |
| **Documentation** | Generate parity report | Review completeness |

**Deliverables:**
- Go source files with Python-equivalent comments
- Test files with table-driven tests
- Parity report
- Coverage report

**Code Pattern:**
```go
// MergeSimilarEntries combines GL entries with the same merge key.
//
// Maps to: merge_similar_entries() in general_ledger.py (lines 273-326)
//
// Python equivalent:
//   def merge_similar_entries(gl_map, precision=None):
//       merged_gl_map = []
//       for entry in gl_map:
//           same_head = check_if_in_list(entry, merged_gl_map)
//           if same_head:
//               same_head.debit += entry.debit
func MergeSimilarEntries(glMap []GLEntry) []GLEntry {
    // Go implementation
}
```

### Phase 4: Validation

**Goal:** Prove Go implementation matches Python behavior.

| Activity | AI Role | Human Role |
|----------|---------|------------|
| **Parity Testing** | Generate tests with same inputs | Set up test infrastructure |
| **Edge Case Discovery** | Analyze Python for unusual branches | Validate edge cases matter |
| **Error Message Matching** | Compare error messages | Accept intentional differences |
| **Performance Baseline** | Generate benchmark tests | Set performance targets |

**Deliverables:**
- Passing test suite (target: 90%+ coverage)
- Parity verification report
- Performance comparison

### Phase 5: Cutover

**Goal:** Safely route traffic to Go implementation.

| Activity | AI Role | Human Role |
|----------|---------|------------|
| **Shadow Mode Setup** | Generate comparison logic | Configure infrastructure |
| **Difference Analysis** | Analyze shadow logs | Decide on fixes |
| **Rollback Plan** | Document feature flag config | Approve rollback criteria |
| **Monitoring** | Suggest metrics to track | Set up dashboards |

**Deliverables:**
- Shadow mode running
- Feature flag configuration
- Monitoring dashboard
- Runbook

---

## AI Capabilities Matrix

### What AI Does Well

| Capability | Strength | Example from ERPNext |
|------------|----------|---------------------|
| **Pattern Recognition** | High | Identified 3 validation methods in mode_of_payment.py |
| **Business Logic Extraction** | High | Extracted duplicate company check algorithm |
| **Test Case Generation** | High | Generated 19 test cases from 3 validation rules |
| **Code Translation** | High | Converted `frappe.throw()` to Go `ValidationError` |
| **Documentation** | High | Created parity reports with side-by-side comparisons |
| **Interface Design** | Medium-High | Proposed `AccountLookup` and `POSChecker` interfaces |
| **Error Message Matching** | High | Preserved ERPNext error text in Go errors |

### What AI Needs Help With

| Challenge | Example | Mitigation |
|-----------|---------|------------|
| **Business Context** | Why does POS check exist? (regulatory? UX?) | Human provides context |
| **Edge Cases Beyond Code** | What happens in production edge cases? | Human provides production experience |
| **Performance Requirements** | How fast is "fast enough"? | Human sets targets |
| **Integration Decisions** | PostgreSQL vs MariaDB? | Human makes strategic choice |
| **Naming Conventions** | `ModeOfPayment` vs `PaymentMode`? | Human decides domain language |

### Capability by Task Type

| Task Type | AI Capability | Human Requirement |
|-----------|--------------|-------------------|
| Read Python code | Autonomous | None |
| Identify validation methods | Autonomous | Confirmation |
| Extract business rules | Autonomous | Validation |
| Design Go interfaces | Propose | Approve |
| Generate Go code | Draft | Review |
| Write tests | Draft | Extend |
| Identify missing rules | Limited | Required |
| Production deployment | None | Full ownership |

---

## Human-AI Collaboration Model

### The Human Role: Strategic Decisions

| Decision Type | Example | Why Human? |
|---------------|---------|------------|
| **Module Priority** | "Start with Mode of Payment, not Sales Invoice" | Business impact assessment |
| **Architecture** | "Use hexagonal architecture" | Long-term maintainability |
| **Technology** | "Go over Rust" | Team skills, ecosystem |
| **Parity vs Improvement** | "Match Python behavior exactly first" | Risk management |
| **Cutover Timing** | "After Q1 close" | Business calendar |

### The AI Role: Tactical Execution

| Task Type | Example | AI Advantage |
|-----------|---------|--------------|
| **Code Reading** | "Find all `frappe.throw()` calls" | Speed, completeness |
| **Pattern Extraction** | "Identify validation patterns" | Pattern recognition |
| **Code Generation** | "Write Go struct for JSON schema" | Consistent output |
| **Test Generation** | "Create table-driven tests" | Thoroughness |
| **Documentation** | "Generate parity report" | Structured output |

### Approval Gates

| Gate | Trigger | Approver | Criteria |
|------|---------|----------|----------|
| **Design Review** | Interface definitions complete | Tech Lead | Matches architecture principles |
| **Code Review** | Implementation complete | Developer | Follows patterns, has tests |
| **Parity Review** | Parity report complete | QA | All business rules covered |
| **Deployment Review** | Shadow mode successful | Ops | No errors in shadow logs |

---

## Practical Workflows

### Workflow 1: DocType Migration Session (2-4 hours)

| Phase | Duration | Activities |
|-------|----------|------------|
| **1. Analysis** | 30-60 min | AI reads Python, extracts schema, lists rules |
| **2. Design** | 30 min | AI proposes interfaces, human refines |
| **3. Implementation** | 60-90 min | AI generates code, human reviews |
| **4. Testing** | 30-60 min | AI generates tests, human adds edge cases |
| **5. Documentation** | 15-30 min | AI generates parity report |

### Workflow 2: Business Rule Extraction

**Prompt Template:**
```
Read the Python file at:
[path/to/file.py]

Extract all business rules in this format:

| Rule ID | Method | Description | Input | Output | Error Condition |
|---------|--------|-------------|-------|--------|-----------------|
| R1 | ... | ... | ... | ... | ... |

Include:
1. Validation rules in validate() and its sub-methods
2. Computation rules
3. State transition rules
4. External dependency rules
```

### Workflow 3: Test Generation

**Prompt Template:**
```
For the business rule:
"[Rule description]"

Generate table-driven tests covering:
1. Happy path (valid inputs)
2. Error path (invalid inputs)
3. Edge cases (empty, single, boundary)
4. Error message verification

Format:
tests := []struct {
    name     string
    input    [Input Type]
    wantErr  error
}{...}
```

### Progress Tracking Template

| Iteration | DocType | Analysis | Design | Implement | Test | Deploy |
|-----------|---------|:--------:|:------:|:---------:|:----:|:------:|
| 1 | Mode of Payment | ✅ | ✅ | ✅ | ✅ | ⏳ |
| 2 | Tax Calculator | ✅ | ✅ | ✅ | ✅ | ⏳ |
| 3 | GL Entry Engine | ✅ | ✅ | ✅ | ⏳ | - |
| 4 | Account Master | - | - | - | - | - |

---

## Lessons from ERPNext Migration

### Lesson 1: Schema First, Logic Second

**What We Learned:**
- Start with `doctype.json` to understand data model
- JSON schema is the source of truth for fields
- Python code may have computed fields not in schema

**Example:**
```
mode_of_payment.json → ModeOfPayment struct
mode_of_payment_account.json → ModeOfPaymentAccount struct
```

### Lesson 2: Validation Methods Are Gold

**What We Learned:**
- ERPNext's `validate()` method contains most business rules
- Sub-methods like `validate_accounts()` are self-documenting
- Error messages reveal business intent

**Example:**
```python
# This method name tells us the rule
def validate_repeating_companies(self):
    # This error message tells us why
    frappe.throw(_("Same Company is entered more than once"))
```

### Lesson 3: External Dependencies Need Interfaces

**What We Learned:**
- `frappe.get_value()` calls indicate external dependencies
- `frappe.db.sql()` queries reveal data access patterns
- Each external call becomes an interface method

**Example:**
```python
# This call → AccountLookup interface
frappe.get_cached_value("Account", entry.default_account, "company")
```

**Becomes:**
```go
type AccountLookup interface {
    GetAccountCompany(accountName string) (string, error)
}
```

### Lesson 4: Legacy Test Files May Be Empty

**What We Learned:**
- ERPNext test files often contain skeleton classes
- Don't assume Python has good test coverage
- AI-generated tests often exceed legacy coverage

**Example:**
```python
# erpnext/accounts/doctype/mode_of_payment/test_mode_of_payment.py
class TestModeofPayment(IntegrationTestCase):
    pass  # No actual tests!
```

**Result:** Go has 19 comprehensive tests vs Python's 0.

### Lesson 5: Complex Algorithms Work Well

**What We Learned:**
- `taxes_and_totals.py` demonstrates complex business logic
- 5 different charge types with cascading dependencies
- AI handled 350+ lines of calculation logic effectively

**Charge Types Migrated:**
```
- Actual (fixed amount, proportional distribution)
- On Net Total (percentage of line item)
- On Previous Row Amount (cascading tax)
- On Previous Row Total (compound tax)
- On Item Quantity (per-unit charge)
```

### Lesson 6: Parity Reports Build Confidence

**What We Learned:**
- Side-by-side Python/Go comparison is invaluable
- Table format shows field-by-field parity
- Test count comparison shows coverage improvement

**Example from PARITY_REPORT.md:**
```
| Metric | Python | Go | Parity |
|--------|--------|-----|--------|
| Test Coverage | 0 tests | 19 tests | Exceeds |
```

---

## Prompting Strategies

### Strategy 1: Context Loading

**Pattern:** Load relevant files before asking questions

```
Read these files:
1. erpnext/accounts/doctype/mode_of_payment/mode_of_payment.py
2. erpnext/accounts/doctype/mode_of_payment/mode_of_payment.json

Then answer: What validation rules exist in this DocType?
```

### Strategy 2: Iterative Refinement

**Pattern:** Start broad, then narrow

1. "List all methods in the file"
2. "Explain what `validate_accounts()` does"
3. "What error does it throw and when?"
4. "Generate Go code for this method"

### Strategy 3: Reference Existing Patterns

**Pattern:** Point to examples in the codebase

```
Using the pattern from modeofpayment/validation.go:

Implement a similar validation for [new rule].

Follow:
- Same error struct pattern
- Same comment style
- Same test structure
```

### Strategy 4: Parity Verification

**Pattern:** Ask AI to verify its own work

```
Compare:

Python:
[paste Python code]

Go:
[paste Go code]

Are they functionally equivalent? List any differences.
```

---

## Quality Gates

### Gate 1: Analysis Complete

| Checkpoint | Criteria | Evidence |
|------------|----------|----------|
| Schema documented | All fields mapped to Go types | Field mapping table |
| Rules extracted | All validation methods identified | Rules table |
| Dependencies listed | All external calls identified | Interface list |

### Gate 2: Design Complete

| Checkpoint | Criteria | Evidence |
|------------|----------|----------|
| Interfaces defined | All external dependencies abstracted | Interface code |
| Structs designed | All entities have Go structs | Model code |
| Errors defined | All error conditions have typed errors | Error definitions |

### Gate 3: Implementation Complete

| Checkpoint | Criteria | Evidence |
|------------|----------|----------|
| Code compiles | `go build ./...` passes | Build log |
| Tests pass | `go test ./...` passes | Test output |
| Coverage adequate | >80% coverage | Coverage report |

### Gate 4: Parity Verified

| Checkpoint | Criteria | Evidence |
|------------|----------|----------|
| Fields match | All Python fields in Go | Parity table |
| Rules match | All validations produce same results | Test cases |
| Errors match | Error messages are equivalent | Error comparison |

---

## References

- [Working Effectively with Legacy Code](https://www.oreilly.com/library/view/working-effectively-with/0131177052/) — Michael Feathers
- [Strangler Fig Pattern](https://martinfowler.com/bliki/StranglerFigApplication.html) — Martin Fowler
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) — Alistair Cockburn

---

## Appendix: ERPNext Migration Statistics

### Current Progress

| Package | Test Count | Coverage | Status |
|---------|------------|----------|--------|
| modeofpayment | 19 | 85.3% | ✅ Complete |
| taxcalc | 24 | 90.2% | ✅ Complete |
| ledger | 25 | 49.1% | 🔄 In Progress |

### Lines of Code Migrated

| Module | Python Lines | Go Lines | Ratio |
|--------|--------------|----------|-------|
| Mode of Payment | ~150 | ~175 | 1.17x |
| Tax Calculator | ~350 | ~440 | 1.26x |
| GL Entry Engine | ~550 | ~650 | 1.18x |

*Note: Go tends to be slightly more verbose due to explicit type handling and error management, but the additional safety is worth the trade-off.*
