# Focus 🎯

> A CLI tool to help developers stay focused and avoid scope creep

**Built in public by [Job Siemerink](https://github.com/n3sty)**

## The Problem

As developers, we've all been there:

- You start working on fixing a critical bug (non-PDF OCR support)
- 2 hours later, you're redesigning the upload button UI
- The original bug? Still not fixed.

This is **scope creep** + **bikeshedding** + **lack of timeboxing**, and it kills productivity.

## The Solution

Focus is a CLI tool that helps you:

- **Set clear goals** with timeboxed work sessions
- **Track drift** when you go down rabbit holes
- **Stay accountable** with periodic focus checks
- **Use git intentionally** - branch per goal, clean history

## Features

### 🎯 Focused Sessions
Start work sessions with a clear goal and time limit:
```bash
focus start "Fix non-PDF OCR support" --time 3h
```

### ✅ Focus Checks
Periodically verify you're still on track with an interactive TUI:
```bash
focus check
```
- Answer if you're still working on your goal
- Log "drifts" when you've wandered off
- Reflect on whether detours are necessary

### 📊 Session Status
See your progress at a glance:
```bash
focus status
```
Shows:
- Time elapsed vs. timebox
- Number of commits made
- Drift log (all the rabbit holes)

### 🏁 Session Review
End sessions with intention:
```bash
focus end
```
Interactive TUI lets you:
- ✅ **Merge to main** if goal achieved
- 📌 **Continue tomorrow** if still in progress
- 🗑️ **Abandon branch** if it was a rabbit hole

## Installation

### Quick Install
```bash
go install github.com/n3sty/focus@latest
```

### Manual Build
```bash
git clone https://github.com/n3sty/focus.git
cd focus
go build -o focus .
go install
```

Make sure `~/go/bin` is in your PATH:
```bash
export PATH="$HOME/go/bin:$PATH"
```

## Usage

### Basic Workflow

1. **Start a focused session:**
   ```bash
   focus start "Implement user authentication" --time 2h
   ```
   This creates a git branch and starts tracking your session.

2. **Work on your goal** and commit regularly:
   ```bash
   git add .
   git commit -m "Add login form component"
   ```

3. **Check yourself periodically:**
   ```bash
   focus check
   ```
   Set a timer reminder or run manually when you feel yourself drifting.

4. **Review your progress:**
   ```bash
   focus status
   ```

5. **End the session:**
   ```bash
   focus end
   ```
   Choose whether to merge, continue tomorrow, or abandon.

### Integration with Existing Timer

If you use a Pomodoro timer (like the `timer` command), integrate focus checks:

```bash
# In your ~/.zshrc
work() {
  timer "${1:-25m}" && terminal-notifier -message 'Run: focus check' \
        -title 'Work Timer - Focus Check!' \
        -sound Crystal
}
```

## Why This Approach?

### The Psychology
- **Timeboxing** creates urgency and prevents perfectionism
- **Explicit drift logging** builds self-awareness
- **Git branches per goal** make it easy to abandon rabbit holes
- **Periodic checks** act as pattern interrupts

### The Git Strategy
- Each focus session = one feature branch
- Easy to see what you actually worked on (`git log`)
- Can revert rabbit holes without losing main work
- Clean merge history when goals complete

## Tech Stack

Built with:
- **Go** - Fast, single-binary distribution
- **Cobra** - Powerful CLI framework
- **Bubble Tea** - Beautiful terminal UIs (Elm Architecture for TUI)
- **Lipgloss** - Styling and layout

## Roadmap

- [x] Core focus session management
- [x] Interactive TUI for checks and reviews
- [x] Git integration with automatic branching
- [ ] Timer integration with notifications
- [ ] Claude API integration for drift validation
- [ ] Session analytics and insights
- [ ] Template goals for common tasks
- [ ] Team shared focus sessions

## Contributing

This project is built in public! Contributions, issues, and feedback are welcome.

## License

MIT License - see LICENSE file for details

## Author

Built by [Job Siemerink](https://github.com/n3sty) - Student, developer, and entrepreneur learning to ship instead of polish.

---

**Built in public 🚀** | Follow the journey at [@n3sty](https://twitter.com/n3sty)
