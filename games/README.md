# 🎮 Code Intelligence Games

**Building skills for AI-assisted legacy modernization**

These games teach the exact skills your Code Intelligence Platform needs. Every challenge is something your tool will eventually automate — but first, you need to understand it manually.

---

## 🎯 Why These Games?

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  THE CODE INTELLIGENCE PIPELINE                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ERPNext Python ──► Parse/Index ──► Search ──► Context ──► AI ──► Go Code  │
│  (legacy)          (Game 3)       (Game 1,2) (Game 5,6)       (verify)     │
│                                                                              │
│  Each game teaches ONE step of this pipeline.                               │
│  Master them manually → Build tools to automate them.                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Active Challenges

| Game | Focus Area | Difficulty | Points |
|------|------------|------------|--------|
| [Game 1: Business Rule Hunter](#game-1-business-rule-hunter) | Rule Extraction | 🟡 Medium | 100 |
| [Game 2: Dependency Mapper](#game-2-dependency-mapper) | Semantic Boundaries | 🟡 Medium | 100 |
| [Game 3: AST Explorer](#game-3-ast-explorer) | Code Parsing | 🔴 Hard | 150 |
| [Game 4: Bounded Context Detective](#game-4-bounded-context-detective) | Architecture | 🔴 Hard | 150 |
| [Game 5: Parity Spec Writer](#game-5-parity-spec-writer) | Migration | 🟡 Medium | 100 |
| [Game 6: Context Crafter](#game-6-context-crafter) | RAG/Retrieval | 🔴 Hard | 150 |

**Total possible: 750 points + bonuses**

---

## Game 1: Business Rule Hunter

> **Skill:** Extracting validation logic from legacy code
> **Why it matters:** Your Code Intelligence tool must find and document business rules automatically

### The Challenge

Find **10 business rules** in ERPNext's `accounts` module where:
1. A validation check is performed
2. An error is raised if the check fails (`frappe.throw()`, `frappe.msgprint()`, `raise`)

### Target Codebase

```
https://github.com/frappe/erpnext
Directory: erpnext/accounts/
```

### Submission Format

For each rule found, document:

| Field | Description |
|-------|-------------|
| **Rule Name** | Human-readable name (you create this) |
| **File** | Path relative to erpnext/ |
| **Function** | Function name containing the rule |
| **Line** | Line number of the validation |
| **Condition** | The actual condition being checked |
| **Error Message** | The error shown to users |
| **Business Impact** | Why this rule exists (1 sentence) |

### Example Submission

```markdown
## Rule 1: Debit Credit Balance

| Field | Value |
|-------|-------|
| **Rule Name** | Debit Credit Balance |
| **File** | accounts/general_ledger.py |
| **Function** | check_if_in_list |
| **Line** | 186 |
| **Condition** | `abs(total_debit - total_credit) > 0.001` |
| **Error Message** | "Debit and Credit not equal for {voucher_type}: {voucher_no}" |
| **Business Impact** | Enforces double-entry accounting principle |
```

### Scoring

- Each complete, correct rule: **+10 points**
- Incomplete documentation: **+5 points**
- Incorrect/made-up rule: **-5 points**
- First team to find 10 valid rules: **+25 bonus**

### How You Know You've Won

✅ You found 10 distinct business rules
✅ Each rule has all 7 fields documented
✅ Each condition is verifiable in the actual code
✅ Your team submitted first (for bonus)

---

## Game 2: Dependency Mapper

> **Skill:** Understanding module boundaries and dependencies
> **Why it matters:** Migration requires knowing what depends on what

### The Challenge

Map all dependencies for **ONE DocType**: `GL Entry`

Location: `erpnext/accounts/doctype/gl_entry/`

### What to Map

1. **Internal Dependencies** — Other ERPNext modules this DocType imports
2. **External Dependencies** — Frappe framework, Python stdlib, third-party
3. **Reverse Dependencies** — What OTHER files import GL Entry

### Submission Format

```markdown
## GL Entry Dependency Map

### Internal Dependencies (ERPNext)
| Import | File | Why It's Needed |
|--------|------|-----------------|
| `erpnext.accounts.utils` | gl_entry.py:5 | Balance calculations |
| ... | ... | ... |

### External Dependencies
| Import | Type | Purpose |
|--------|------|---------|
| `frappe` | Framework | ORM, validation |
| `json` | Stdlib | Data serialization |
| ... | ... | ... |

### Reverse Dependencies (What imports GL Entry)
| File | Import Statement |
|------|------------------|
| accounts/general_ledger.py | `from erpnext.accounts.doctype.gl_entry...` |
| ... | ... |

### Dependency Diagram
[ASCII or Mermaid diagram showing relationships]
```

### Scoring

- Internal dependencies (complete list): **+30 points**
- External dependencies (complete list): **+20 points**
- Reverse dependencies (5+ files): **+30 points**
- Diagram quality: **+20 points**
- First complete submission: **+25 bonus**

### How You Know You've Won

✅ You found ALL imports in gl_entry.py and gl_entry.json
✅ You found at least 5 files that import GL Entry
✅ Your diagram shows the relationships clearly

---

## Game 3: AST Explorer

> **Skill:** Programmatic code analysis using Abstract Syntax Trees
> **Why it matters:** Your Code Intelligence tool parses code — this is how

### The Challenge

Write a **Python script** that analyzes ERPNext code and extracts:

1. **All class definitions** in `erpnext/accounts/doctype/gl_entry/gl_entry.py`
2. **All method names** in those classes
3. **All `frappe.throw()` calls** with their error messages

### Requirements

Use Python's `ast` module OR tree-sitter. Your script must:
- Take a file path as input
- Output structured JSON
- Run without errors

### Expected Output Format

```json
{
  "file": "gl_entry.py",
  "classes": [
    {
      "name": "GLEntry",
      "line": 15,
      "methods": [
        {"name": "validate", "line": 25},
        {"name": "on_submit", "line": 45}
      ]
    }
  ],
  "validation_calls": [
    {
      "type": "frappe.throw",
      "line": 52,
      "message": "Account {0} is frozen"
    }
  ]
}
```

### Submission Format

Submit:
1. Your Python script (as GitHub Gist or in Issue comment)
2. The JSON output when run on `gl_entry.py`
3. Brief explanation of your approach

### Scoring

- Script runs without errors: **+30 points**
- Extracts all classes correctly: **+30 points**
- Extracts all methods correctly: **+30 points**
- Extracts all frappe.throw() calls: **+40 points**
- Uses tree-sitter (bonus for learning new tool): **+20 bonus**
- First working submission: **+25 bonus**

### How You Know You've Won

✅ Your script runs: `python your_script.py gl_entry.py`
✅ Output is valid JSON
✅ All classes, methods, and throws are captured
✅ Your extraction matches manual inspection

### Starter Code

```python
import ast
import json
import sys

def analyze_file(filepath):
    with open(filepath, 'r') as f:
        tree = ast.parse(f.read())

    result = {
        "file": filepath,
        "classes": [],
        "validation_calls": []
    }

    # TODO: Walk the AST and extract information
    # Hint: Use ast.walk(tree) or ast.NodeVisitor

    for node in ast.walk(tree):
        if isinstance(node, ast.ClassDef):
            # TODO: Extract class info
            pass
        if isinstance(node, ast.Call):
            # TODO: Check if it's frappe.throw()
            pass

    return result

if __name__ == "__main__":
    filepath = sys.argv[1]
    result = analyze_file(filepath)
    print(json.dumps(result, indent=2))
```

---

## Game 4: Bounded Context Detective

> **Skill:** Identifying domain boundaries for migration
> **Why it matters:** Strangler Fig requires knowing where to draw the lines

### The Challenge

Analyze ERPNext's `accounts` module and identify **3 bounded contexts** that could be migrated independently.

### What is a Bounded Context?

From Domain-Driven Design: A bounded context is a subsystem with:
- Clear responsibility (one job)
- Defined interfaces (how it talks to other contexts)
- Minimal external dependencies (can be extracted)

### Your Task

For each bounded context you identify:

```markdown
## Bounded Context: [Name]

### Responsibility
[What business capability does this handle?]

### Core Entities (DocTypes)
| DocType | Purpose |
|---------|---------|
| GL Entry | Records individual ledger transactions |
| ... | ... |

### Interfaces (How it communicates)
| Interface | Direction | Description |
|-----------|-----------|-------------|
| `make_gl_entries()` | Inbound | Other modules call this to post entries |
| `get_balance_on()` | Outbound | Queries account balances |
| ... | ... | ... |

### Dependencies
| Depends On | Why | Can Be Mocked? |
|------------|-----|----------------|
| Company | Currency, fiscal year | Yes |
| Account | Account hierarchy | Yes |
| ... | ... | ... |

### Migration Feasibility
| Factor | Assessment |
|--------|------------|
| **Complexity** | Low / Medium / High |
| **Dependencies** | Count: X internal, Y external |
| **Test Coverage** | Can behavior be verified? |
| **Recommendation** | Migrate first / second / later |
```

### Scoring

- Each well-defined bounded context: **+40 points**
- Clear interfaces identified: **+10 points each**
- Dependencies accurately mapped: **+10 points**
- Migration recommendation justified: **+10 points**
- First submission with 3 valid contexts: **+25 bonus**

### How You Know You've Won

✅ You identified 3 distinct bounded contexts
✅ Each has clear responsibility (not overlapping)
✅ Interfaces are actual functions/methods in the code
✅ Your migration recommendation makes sense

---

## Game 5: Parity Spec Writer

> **Skill:** Translating Python behavior into implementation spec
> **Why it matters:** AI assistants need clear specs to generate correct Go code

### The Challenge

Write a **Parity Specification** for this ERPNext Python function:

```python
# From erpnext/accounts/utils.py
def get_balance_on(account, date=None, party_type=None, party=None,
                   company=None, in_account_currency=True,
                   cost_center=None, ignore_account_permission=False):
    """
    Get the balance of an account on a specific date.

    Returns the sum of (debit - credit) for all GL entries
    up to and including the given date.
    """
```

### Your Parity Spec Must Include

```markdown
## Parity Specification: get_balance_on

### Function Signature (Go)
```go
func GetBalanceOn(params BalanceParams) (decimal.Decimal, error)
```

### Input Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| account | string | Yes | - | Account name |
| date | time.Time | No | today | Balance as of date |
| ... | ... | ... | ... | ... |

### Business Logic
1. [Step-by-step description of what the function does]
2. [Include edge cases]
3. [Include error conditions]

### Test Cases for Parity Verification
| Scenario | Input | Expected Output |
|----------|-------|-----------------|
| Basic balance | account="Sales", date=2024-01-15 | 50000.00 |
| No entries | account="Empty", date=2024-01-01 | 0.00 |
| With party filter | account="Debtors", party="ACME" | 11800.00 |
| Future date | account="Sales", date=2099-01-01 | [same as today] |

### SQL Equivalent
```sql
SELECT SUM(debit) - SUM(credit) as balance
FROM `tabGL Entry`
WHERE account = %(account)s
  AND posting_date <= %(date)s
  AND is_cancelled = 0
  [additional filters...]
```

### Edge Cases
| Case | Behavior |
|------|----------|
| Account doesn't exist | Return 0 or error? |
| Date is None | Use today's date |
| ... | ... |
```

### Scoring

- Complete function signature: **+15 points**
- All parameters documented: **+20 points**
- Business logic accurate: **+25 points**
- Test cases (5+ scenarios): **+25 points**
- Edge cases identified: **+15 points**
- First complete submission: **+25 bonus**

### How You Know You've Won

✅ A developer could implement the Go function from your spec alone
✅ Test cases would catch if Go doesn't match Python
✅ Edge cases are documented (not just happy path)

---

## Game 6: Context Crafter

> **Skill:** Selecting relevant code for AI context windows
> **Why it matters:** RAG quality determines AI output quality

### The Challenge

A developer asks:
> **"How does ERPNext create GL entries when a Sales Invoice is submitted?"**

Your job: Select the **minimum code** that answers this question completely.

### Constraints

- Maximum **2000 tokens** (~8000 characters)
- Must include actual code snippets (not just descriptions)
- Must answer the question completely

### Submission Format

```markdown
## Context for: "How does ERPNext create GL entries when a Sales Invoice is submitted?"

### Selected Files (in order of relevance)

#### 1. [filename] (lines X-Y)
```python
[code snippet]
```
**Why included:** [1 sentence]
**Tokens:** ~XXX

#### 2. [filename] (lines X-Y)
```python
[code snippet]
```
**Why included:** [1 sentence]
**Tokens:** ~XXX

[Continue for each snippet...]

### Total Tokens: XXXX / 2000

### What This Context Enables
- [What questions can an AI answer with this context?]
- [What questions can it NOT answer?]
```

### Scoring

| Criteria | Points |
|----------|--------|
| Answers the question completely | +40 |
| Under 2000 tokens | +20 |
| Most concise complete answer | +30 (only 1 team) |
| Includes entry point (on_submit) | +15 |
| Includes GL creation (make_gl_entries) | +15 |
| Includes actual GL Entry structure | +15 |
| Explains what's NOT included and why | +15 |
| First submission | +25 bonus |

### How You Know You've Won

✅ An AI given ONLY your context could explain the GL creation flow
✅ You used fewer tokens than other teams (efficiency matters)
✅ You didn't miss critical code paths

### Hints

Start here:
- `erpnext/accounts/doctype/sales_invoice/sales_invoice.py` → `on_submit()`
- `erpnext/controllers/accounts_controller.py` → `make_gl_entries()`
- `erpnext/accounts/general_ledger.py` → `make_gl_entries()`

---

## 📊 Leaderboard

| Rank | Team | G1 | G2 | G3 | G4 | G5 | G6 | Bonus | Total |
|------|------|----|----|----|----|----|----|-------|-------|
| 🥇 | TBD | - | - | - | - | - | - | - | - |
| 🥈 | TBD | - | - | - | - | - | - | - | - |
| 🥉 | TBD | - | - | - | - | - | - | - | - |

---

## 📝 Submission Instructions

1. **Create your submission** as a Markdown document
2. **Post in the designated GitHub Issue** for each game
3. **One submission per team** — coordinate with your teammates
4. **Include team name** at the top of every submission

---

## 🏆 Prizes

| Achievement | Prize |
|-------------|-------|
| **Overall Winner** | LinkedIn recommendation + Certificate |
| **Most Creative Solution** | Featured in documentation |
| **Best AST Script** | Code included in actual tool |
| **Fastest Submission** | Recognition in team meeting |

---

## 💡 Getting Help

| Hint Level | Cost | How to Request |
|------------|------|----------------|
| **Free Hint** | 0 pts | Ask in Teams #General |
| **Guided Hint** | -10 pts | DM mentor |
| **Full Walkthrough** | -25 pts | Scheduled call |

---

## 🔗 Resources

- **ERPNext Source:** https://github.com/frappe/erpnext
- **Frappe Framework:** https://github.com/frappe/frappe
- **Python AST Docs:** https://docs.python.org/3/library/ast.html
- **Tree-sitter Python:** https://github.com/tree-sitter/tree-sitter-python

---

*These games teach you to think like a Code Intelligence tool. Master them, then build the tool that does this automatically.*
