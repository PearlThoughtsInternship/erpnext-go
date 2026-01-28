# Teams Announcement Templates

Copy-paste these into Microsoft Teams #General channel to launch and manage the games.

---

## 🚀 Game Launch Announcement

**Post this to kick off all 3 games:**

```
🎮 CODE INTELLIGENCE GAMES — NOW LIVE!

Teams, your challenge awaits!

We've prepared 3 games that test your Code Intelligence skills:

🔍 **Game 1: Symbol Hunter** (50-100 points)
Find business rules hidden in the ERPNext codebase. Your Code Intelligence tool will need to do this automatically — prove you can do it first!

🐛 **Game 2: The Bug Hunt** (115+ points)
A PR with 7 hidden bugs awaits your review. Easy, Medium, and Hard bugs are scattered throughout. Find them all!

🌳 **Game 3: The Strangler's Trail** (100-175 points)
An epic puzzle journey across code, GitHub, and Teams. Follow the clues, decode the messages, reveal the wisdom of the Strangler Fig.

📍 **Start here:** https://github.com/PearlThoughtsInternship/erpnext-go/tree/main/games

⏰ **Deadline:** [INSERT DATE]

**Teams:**
- Team Alpha: [INSERT NAMES]
- Team Beta: [INSERT NAMES]
- Team Gamma: [INSERT NAMES]

The games run in parallel — choose your strategy! Start with the one that excites you most.

May the best Code Intelligence team win! 🏆
```

---

## 🌳 Stage 5: The Final Fragment (Post after teams reach Stage 4)

**Post this for the final clue of The Strangler's Trail:**

```
🌳 The Strangler's Trail — Final Fragment

For those who have followed the trail...

The fig's patience is legendary. It doesn't rush.

Decode the final piece:

U1RSVUNUVVJFIEJFU0lERSBJVA==

Combine all five fragments in order to reveal the complete wisdom.

Post your answer in #General with this format:

🎯 STRANGLER'S TRAIL COMPLETE
Team: [Your Team Name]
Final Answer: [The complete phrase]
Journey Log: [Brief description of each stage you solved]

Good luck, travelers! 🍀
```

---

## 💡 Daily Hint Posts

### Day 2: Symbol Hunter Hint

```
💡 HINT (Free): Symbol Hunter

The ERPNext accounts module has ~50 Python files.

Start with the obvious ones:
• gl_entry.py
• payment_entry.py
• journal_entry.py
• accounts_controller.py

Use grep or your favorite code search tool. The pattern you're looking for:
frappe.throw(...balance...)
```

### Day 3: Bug Hunt Hint

```
💡 HINT (Free): The Bug Hunt

The easy bugs are really easy. If you haven't found 2 bugs in the first 5 minutes, slow down and read the function signature carefully.

Questions to ask:
• What should this function do vs what does it actually do?
• What happens with edge cases (nil, empty, zero)?
• Is the comparison operator correct?
```

### Day 4: Strangler's Trail Hint

```
💡 HINT (Free): The Strangler's Trail

Stage 1 says "where GL entries are born" — that's the struct definition, not where they're created.

Look in the ledger package. Read the comments carefully. Some comments are... special. 🌳
```

---

## 📊 Progress Updates

### Mid-Week Update Template

```
📊 GAME PROGRESS UPDATE

**Leaderboard (as of [DATE]):**

| Rank | Team | Points | Games Completed |
|------|------|--------|-----------------|
| 🥇 | [Team] | [X] pts | Game 1 ✓, Game 2 (3/7), Game 3 (Stage 2) |
| 🥈 | [Team] | [X] pts | Game 1 (partial), Game 2 ✓, Game 3 (Stage 1) |
| 🥉 | [Team] | [X] pts | Starting strong! |

**Highlights:**
- [Team] found a sneaky bug in Game 2! +30 points
- [Team] cracked Stage 3 of the Strangler's Trail!
- First hint request of the day: [question]

Keep hunting! 🔍🐛🌳
```

### Game Winner Announcement

```
🏆 GAME [X] WINNER!

Congratulations to **[TEAM NAME]** for completing [Game Name] first!

Their achievement:
• [Specific accomplishment]
• [Points earned]
• [Time taken]

Special recognition:
• [Any notable strategy or discovery]

The other teams can still earn full points — just not the first-place bonus. Keep going!
```

---

## 🎉 Final Ceremony

```
🏆 CODE INTELLIGENCE GAMES — FINAL RESULTS

After [X] days of hunting, debugging, and trail-following...

**FINAL STANDINGS:**

🥇 **1st Place: [TEAM NAME]** — [TOTAL] points
   • Symbol Hunter: [X] pts
   • Bug Hunt: [X] pts
   • Strangler's Trail: [X] pts

🥈 **2nd Place: [TEAM NAME]** — [TOTAL] points
   • [breakdown]

🥉 **3rd Place: [TEAM NAME]** — [TOTAL] points
   • [breakdown]

**Special Awards:**
• 🔍 **Best Symbol Hunter**: [Name] — found [X] patterns
• 🐛 **Bug Squasher**: [Name] — found the hardest bug (Race Condition)
• 🌳 **Trail Master**: [Name] — best journey log

**The Wisdom Revealed:**
"WRAPS AROUND THE OLD TREE WHILE GROWING ITS OWN STRUCTURE BESIDE IT"

This is how we modernize legacy systems. Not by cutting down the old tree, but by growing something new alongside it until it can stand on its own.

Thank you all for playing! These skills — code exploration, bug hunting, pattern recognition — are exactly what you'll use when building the Code Intelligence Platform.

🎓 Certificates and prizes will be distributed in Friday's sync.
```

---

## 🆘 Troubleshooting Posts

### If Teams Are Stuck

```
🆘 STUCK? Here's Help!

Noticed some teams haven't made progress in 24 hours. Here's a lifeline:

**Symbol Hunter Stuck?**
Try this command in the ERPNext repo:
grep -r "frappe.throw" --include="*.py" erpnext/accounts/ | grep -i "balance"

**Bug Hunt Stuck?**
The 2 Easy bugs are in:
• Function naming convention
• Missing defensive check

Read the FIRST function carefully!

**Strangler's Trail Stuck?**
Stage 1: The file starts with "m" and ends with ".go"
It's in the ledger/ directory.

Need more help? DM me for a hint (costs 10 points).
```

### Technical Issues

```
⚠️ TECHNICAL ISSUE RESOLVED

[Describe issue]

**Fix:** [What was done]

**If you were affected:**
• [Instructions to continue]
• Any lost progress will be restored

Sorry for the inconvenience! Games continue as normal.
```

---

## 📝 Notes for Game Master

1. **Before posting Stage 5**, verify teams have reached Stage 4 by checking:
   - PR comment activity
   - Issue activity
   - Any questions about "documentation" or "INTERN_GUIDE"

2. **Track submissions** in a spreadsheet with timestamps

3. **Watch for cheating signals:**
   - Teams submitting identical answers at same time
   - Answers appearing faster than humanly possible
   - Direct copy of answer key formatting

4. **Adjust difficulty** if needed:
   - Too hard: Release more hints, extend deadline
   - Too easy: Add bonus challenges

5. **Celebrate effort**, not just winning — even partial completion shows learning
