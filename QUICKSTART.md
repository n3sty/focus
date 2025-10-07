# Quick Start Guide

Get started with Focus in 2 minutes!

## Installation

```bash
# If you haven't already
go install github.com/n3sty/focus@latest

# Make sure ~/go/bin is in your PATH
export PATH="$HOME/go/bin:$PATH"
```

## Your First Focus Session

### 1. Navigate to your project
```bash
cd ~/dev/my-project
```

### 2. Start a focus session
```bash
focus start "Fix the login bug" --time 2h
```

You'll see:
```
🎯 Starting focus session: Fix the login bug
⏱️  Timebox: 2h

✓ Created branch: focus/fix-the-login-bug
✓ Session saved

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 Focus Session Active
   Goal: Fix the login bug
   Time: 2h
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 3. Work on your goal
Make commits as you progress:
```bash
git add .
git commit -m "Identify the auth token issue"
```

### 4. Check yourself
Every 25-30 minutes (or when you feel yourself drifting):
```bash
focus check
```

Interactive prompt appears:
```
🎯 Current Goal
   Fix the login bug

Are you still working on this goal?

[y] Yes - Still on track!
[n] No  - I've drifted...
[d] Defer - Check me later
```

If you drifted, it asks what you're working on instead and why.

### 5. Check your progress
```bash
focus status
```

Shows:
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎯 Focus Session Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Goal:     Fix the login bug
Started:  10:30 AM
Elapsed:  1h 45m
Timebox:  2h
Branch:   focus/fix-the-login-bug
Commits:  5
Drifts:   1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 6. End your session
```bash
focus end
```

Interactive TUI shows your summary and asks what to do:
```
🏁 End Focus Session

┌─────────────────────────────────────┐
│ 🎯  Fix the login bug               │
│                                     │
│ ⏱️  Time: 1h 45m (planned: 2h)     │
│ 📝  Commits: 5                      │
│ 🐰  Drifts: 1                       │
└─────────────────────────────────────┘

What do you want to do?

▸ ✅ Merge to main
   Goal achieved! Merge branch and complete session

  📌 Continue tomorrow
  🗑️ Abandon
```

Select your choice and hit Enter!

## Integration with Timer

Add to your `~/.zshrc`:

```bash
work() {
  timer "${1:-25m}" && terminal-notifier -message 'Run: focus check' \
        -title 'Work Timer - Focus Check!' \
        -sound Crystal
}
```

Now when your Pomodoro ends, you'll get a reminder to run `focus check`!

## Tips

1. **Set realistic timeboxes**: 2-3 hours max per session
2. **Run `focus check` honestly**: The tool only works if you're honest with yourself
3. **Don't fear abandoning**: That's the whole point! Rabbit holes are okay to delete
4. **Commit often**: Makes the git log more meaningful
5. **Read your drift logs**: They reveal patterns in your behavior

## Common Workflows

### Starting your day
```bash
# Decide on your goal
focus start "Implement user profile page" --time 3h

# Start working
```

### Mid-session check
```bash
# Every 25-30 min or when you feel distracted
focus check
```

### End of session
```bash
# Before lunch, end of day, or when timebox expires
focus end
```

### Checking multiple projects
```bash
# Each project can have its own focus session
cd ~/dev/project-a
focus start "Feature A" --time 2h

# Later...
cd ~/dev/project-b
focus start "Bug fix B" --time 1h
```

Each `.focus` directory is independent!

## Need Help?

```bash
focus --help
focus start --help
focus check --help
```

---

**Ready to stay focused?** 🎯

Start your first session now!
