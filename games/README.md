# 🎮 Internship Games

Welcome to the Code Intelligence Challenge Series!

These games test your understanding of legacy codebases, migration patterns, and Code Intelligence skills. Work with your team, explore the codebase, and may the best team win!

---

## Active Games

| Game | Status | Difficulty | Prize |
|------|--------|------------|-------|
| [Game 1: Symbol Hunter](#game-1-symbol-hunter) | 🟢 Active | Medium | First to solve |
| [Game 2: The Bug Hunt](#game-2-the-bug-hunt) | 🟢 Active | Hard | Most bugs found |
| [Game 3: The Strangler's Trail](#game-3-the-stranglers-trail) | 🟢 Active | Epic | Complete the journey |

---

## How Teams Work

You've been assigned to teams. Each team:
- Works together (divide and conquer encouraged!)
- Submits ONE answer per challenge (team consensus)
- Posts answers in the designated GitHub Issue or Teams channel
- Earns points for speed AND accuracy

---

## Game 1: Symbol Hunter

**Objective:** Your Code Intelligence tool needs to find business rules in legacy code. Prove you can do it manually first!

**The Challenge:**
```
Find all functions in the ERPNext Python codebase (erpnext/accounts/) that:
1. Call frappe.throw() - these are validation rules
2. The throw message contains the word "balance" (case-insensitive)

For each function found, report:
- File path
- Function name
- Line number of the frappe.throw() call
- The exact error message
```

**Submission:** Create a comment on GitHub Issue `#SYMBOL_HUNTER_ISSUE` with your findings.

**Scoring:**
- Each correct finding: +10 points
- Each incorrect finding: -5 points
- First team to find ALL: +50 bonus points
- Time bonus: +5 points for each team you beat

**Hint:** The ERPNext repo is at `https://github.com/frappe/erpnext`. Focus on `erpnext/accounts/` directory.

---

## Game 2: The Bug Hunt

**Objective:** Code review is critical for quality. Can you spot the bugs?

**The Challenge:**

A PR has been submitted with **exactly 7 bugs** hidden in the code. The bugs range from:
- 2 Easy (syntax/obvious)
- 3 Medium (logic errors)
- 2 Hard (subtle edge cases)

**Your Mission:** Find all 7 bugs and explain why each is wrong.

**Submission:** Add review comments directly on PR `#BUG_HUNT_PR`. Each bug should be a separate review comment with:
1. The problematic line
2. What's wrong
3. The fix

**Scoring:**
- Easy bug found: +5 points
- Medium bug found: +15 points
- Hard bug found: +30 points
- False positive (marking correct code as bug): -10 points
- First team to find all 7: +50 bonus points

---

## Game 3: The Strangler's Trail

**Objective:** Follow the clues across Teams, GitHub, and code to complete the migration journey.

**The Challenge:**

This is a multi-stage puzzle. Each stage gives you a clue to the next location.

### Stage 1: The Beginning
> *"The strangler fig starts small. Find where GL entries are born in the Go codebase.
> The file that defines their structure holds your first clue — look for a comment
> that doesn't belong."*

**What to do:** Find the clue, decode it, proceed to Stage 2.

### Stage 2-5: Follow the Trail
Each clue leads to the next. Stages span:
- Code comments in `erpnext-go`
- GitHub Issues
- Teams channel posts
- Documentation files

### Final Stage: The Revelation
Assemble all clues to reveal the **final answer** — a phrase that summarizes the Strangler Fig Pattern.

**Submission:** Post your final answer in Teams #General with format:
```
🎯 STRANGLER'S TRAIL COMPLETE
Team: [Your Team Name]
Final Answer: [The phrase]
Journey Log: [Brief description of each stage]
```

**Scoring:**
- Complete the trail: +100 points
- First team to complete: +75 bonus points
- Best journey log (most detailed): +25 bonus points

---

## Leaderboard

| Rank | Team | Game 1 | Game 2 | Game 3 | Total |
|------|------|--------|--------|--------|-------|
| 🥇 | TBD | - | - | - | - |
| 🥈 | TBD | - | - | - | - |
| 🥉 | TBD | - | - | - | - |

*Updated after each game completion*

---

## Rules

1. **No sabotage** — Don't delete others' submissions or mislead other teams
2. **Document your process** — Screenshots, commands used, reasoning
3. **Ask clarifying questions** — Post in #General, answers help everyone
4. **Have fun** — These skills directly apply to your Code Intelligence work!

---

## Getting Help

Stuck? Hints available:
- **Hint 1:** Free (ask in Teams)
- **Hint 2:** Costs 10 points (ask in Teams)
- **Hint 3:** Costs 25 points (direct message mentor)

---

*May the best Code Intelligence team win! 🏆*
