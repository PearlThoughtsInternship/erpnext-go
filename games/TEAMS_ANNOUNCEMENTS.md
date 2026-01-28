# Teams Announcement Templates

Copy-paste these into Microsoft Teams #General channel.

---

## 🚀 Game Launch Announcement

```
🎮 CODE INTELLIGENCE GAMES — NOW LIVE!

6 challenges that teach the exact skills for building AI-assisted legacy modernization tools.

Each game = one step in the Code Intelligence pipeline:
ERPNext Python → Parse → Search → Context → AI → Go Code

📍 **Full Rules & Examples:**
https://github.com/PearlThoughtsInternship/erpnext-go/tree/main/games

📖 **Solved Example (Study This First):**
https://github.com/PearlThoughtsInternship/erpnext-go/blob/main/games/SOLVED_EXAMPLE_GAME1.md

---

**THE 6 CHALLENGES:**

🔍 **Game 1: Business Rule Hunter** (100 pts)
Find 10 validation rules in ERPNext Python. This IS what your Code Intelligence tool must do.

🗺️ **Game 2: Dependency Mapper** (100 pts)
Map all imports for GL Entry. Understanding module boundaries = understanding where to cut.

🌳 **Game 3: AST Explorer** (150 pts)
Write a Python script using `ast` module to extract classes and methods. This is how indexers work.

🏛️ **Game 4: Bounded Context Detective** (150 pts)
Identify 3 subsystems in ERPNext that could be migrated independently. DDD in practice.

📋 **Game 5: Parity Spec Writer** (100 pts)
Write a spec for `get_balance_on()` so detailed that a dev could implement Go from it alone.

🎯 **Game 6: Context Crafter** (150 pts)
Select minimum code to answer "How does ERPNext create GL entries?" Under 2000 tokens. This is RAG.

---

**TEAMS:**
- Team Alpha: [names]
- Team Beta: [names]
- Team Gamma: [names]

⏰ **Deadline:** [DATE]

**Submit answers as comments on the GitHub Issues (links in README).**

Go! 🏁
```

---

## 📊 How to Create Submission Issues

Create these 6 GitHub Issues for submissions:

### Issue Template for Each Game

**Title:** `🎮 Game [X]: [Name] — Submissions`

**Body:**
```markdown
## Game [X]: [Name]

**Points:** [X] points + bonuses

### Challenge Summary
[1-2 sentence description]

### Submission Format
Post your team's submission as a comment below.

**Required format:**
- Team name at top
- All required fields filled
- Verifiable against actual ERPNext code

### Resources
- [Link to game rules in README]
- [Link to solved example if applicable]

### Scoring
- [Scoring breakdown]
- First complete submission: +25 bonus

---

**Submissions below ⬇️**
```

---

## 💡 Daily Hints (Post One Per Day)

### Day 1: Getting Started

```
💡 Day 1 Tip: Getting Started

Haven't started yet? Here's the path:

1. Clone ERPNext: git clone https://github.com/frappe/erpnext
2. Read the solved example for Game 1
3. Pick ONE game to start with (recommend Game 1 — clearest)
4. Use grep/search to find patterns

The solved example shows EXACTLY how to submit. Study it.

Link: https://github.com/PearlThoughtsInternship/erpnext-go/blob/main/games/SOLVED_EXAMPLE_GAME1.md
```

### Day 2: Finding Business Rules

```
💡 Day 2 Tip: Finding Business Rules (Game 1)

The magic command:

grep -rn "frappe.throw" erpnext/accounts --include="*.py" | head -50

This shows you 50 potential rules. Now:
1. Pick one
2. Read the surrounding code
3. Document the condition
4. Explain why it matters

Remember: 10 rules = 100 points. Quality over speed.
```

### Day 3: AST Explorer Help

```
💡 Day 3 Tip: AST Explorer (Game 3)

Python's ast module walkthrough:

import ast

code = open("gl_entry.py").read()
tree = ast.parse(code)

for node in ast.walk(tree):
    if isinstance(node, ast.ClassDef):
        print(f"Class: {node.name} at line {node.lineno}")
    if isinstance(node, ast.FunctionDef):
        print(f"Function: {node.name} at line {node.lineno}")

Start here. Then add frappe.throw() detection.

Docs: https://docs.python.org/3/library/ast.html
```

### Day 4: Context Crafter Strategy

```
💡 Day 4 Tip: Context Crafter (Game 6)

Strategy for finding minimum context:

1. Start at entry point: sales_invoice.py → on_submit()
2. Follow the call chain: on_submit() → make_gl_entries()
3. Stop when you hit the actual GL Entry creation
4. Count tokens (rough: 4 chars = 1 token)

The question is "How does ERPNext create GL entries?"
You need: trigger → creation → storage. Nothing more.

Less is more. The team with smallest complete answer wins the bonus.
```

### Day 5: Bounded Context Hints

```
💡 Day 5 Tip: Bounded Contexts (Game 4)

Three candidates in accounts module:

1. **GL Posting** — general_ledger.py + gl_entry/
   Interface: make_gl_entries()

2. **Account Master** — doctype/account/
   Interface: get_account(), get_balance_on()

3. **Payment Ledger** — payment_ledger_entry/
   Interface: make_payment_ledger_entry()

Your job: Prove these are separate contexts with clear boundaries.
Look for: What does each import? What imports each?
```

---

## 📊 Progress Update Template

```
📊 GAME PROGRESS — Day [X]

**Leaderboard:**

| Team | G1 | G2 | G3 | G4 | G5 | G6 | Total |
|------|----|----|----|----|----|----|-------|
| Alpha | 80 | - | 30 | - | - | - | 110 |
| Beta | 100 | 50 | - | - | - | - | 150 |
| Gamma | 60 | - | - | 40 | - | - | 100 |

**Highlights:**
- Team Beta completed Game 1 with 10 valid rules! 🎉
- Team Alpha's AST script runs but misses some methods
- Great discussion in Issue #43 about boundary detection

**Tips:**
- Game 1 solved example is your friend
- Quality > Speed (except for bonuses)

Keep going! 🚀
```

---

## 🏆 Game Winner Announcement

```
🏆 GAME [X] WINNER: Team [NAME]!

[Team] completed [Game] with [score] points!

**What they did well:**
- [Specific achievement]
- [Quality of submission]
- [Creative approach if any]

**Their submission is now the benchmark.**
Other teams: You can still earn full points, just not the first-place bonus.

Direct link to their submission: [link]
```

---

## 🎉 Final Results Ceremony

```
🏆 CODE INTELLIGENCE GAMES — FINAL RESULTS

After [X] days of hunting rules, mapping dependencies, and crafting contexts...

**FINAL STANDINGS:**

🥇 **1st Place: Team [NAME]** — [XXX] points
   - G1: [X] | G2: [X] | G3: [X] | G4: [X] | G5: [X] | G6: [X]
   - Bonuses: [X]

🥈 **2nd Place: Team [NAME]** — [XXX] points

🥉 **3rd Place: Team [NAME]** — [XXX] points

---

**SPECIAL AWARDS:**

🔍 **Best Rule Hunter:** [Name]
   Found the most obscure but valid business rule

🌳 **Best AST Script:** [Name]
   Cleanest, most reusable code

🎯 **Most Efficient Context:** [Name]
   Answered the question in fewest tokens

📚 **Best Documentation:** [Name]
   Clearest, most helpful submission

---

**WHAT YOU LEARNED:**

| Game | Skill | Your Tool Will |
|------|-------|----------------|
| Business Rule Hunter | Finding validation logic | Extract rules automatically |
| Dependency Mapper | Understanding boundaries | Map module relationships |
| AST Explorer | Parsing code structure | Index symbols and types |
| Bounded Context Detective | DDD analysis | Identify migration boundaries |
| Parity Spec Writer | Behavior documentation | Generate migration specs |
| Context Crafter | RAG selection | Provide focused AI context |

**These are the exact skills needed to build your Code Intelligence Platform.**

You did them manually. Now automate them.

🎓 Certificates will be distributed in Friday's sync.
```

---

## 🆘 If Teams Are Stuck

```
🆘 Stuck? Here's Help.

**Haven't submitted anything yet?**
1. Start with Game 1 — it's the most straightforward
2. Read the solved example COMPLETELY
3. Clone ERPNext and run: grep -rn "frappe.throw" erpnext/accounts

**Stuck on AST (Game 3)?**
The starter code is in the README. It runs. Start there.
Just add: if isinstance(node, ast.ClassDef): print(node.name)

**Stuck on Bounded Context (Game 4)?**
Start by listing what files are in accounts/doctype/
Group them by "responsibility" — what business job does each do?

**Need more help?**
- Free hint: Ask in #General
- Guided hint (-10 pts): DM mentor
- We WANT you to succeed. Ask for help!
```
