package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/n3sty/focus/internal/git"
	"github.com/n3sty/focus/internal/session"
)

type endAction int

const (
	actionMerge endAction = iota
	actionContinue
	actionAbandon
)

type EndModel struct {
	session    *session.Session
	selected   int
	commits    int
	elapsed    time.Duration
	choice     endAction
	confirmed  bool
}

func NewEndModel(sess *session.Session) EndModel {
	commits, _ := git.GetCommitsSince(sess.StartTime)
	elapsed := time.Since(sess.StartTime)

	return EndModel{
		session:   sess,
		selected:  0,
		commits:   commits,
		elapsed:   elapsed,
		confirmed: false,
	}
}

func (m EndModel) Init() tea.Cmd {
	return nil
}

func (m EndModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if !m.confirmed {
				return m, tea.Quit
			}

		case tea.KeyUp, tea.KeyShiftTab:
			if m.selected > 0 {
				m.selected--
			}

		case tea.KeyDown, tea.KeyTab:
			if m.selected < 2 {
				m.selected++
			}

		case tea.KeyEnter:
			m.choice = endAction(m.selected)
			m.confirmed = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m EndModel) View() string {
	if m.confirmed {
		return m.renderConfirmation()
	}

	var b strings.Builder

	// Header
	header := TitleStyle.Render("🏁 End Focus Session")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Summary box
	summary := m.renderSummary()
	summaryBox := BoxStyle.Render(summary)
	b.WriteString(summaryBox)
	b.WriteString("\n\n")

	// Drift log if any
	if len(m.session.Drifts) > 0 {
		driftLog := m.renderDriftLog()
		b.WriteString(driftLog)
		b.WriteString("\n\n")
	}

	// Options
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("What do you want to do?"))
	b.WriteString("\n\n")

	options := []struct {
		label string
		desc  string
		emoji string
	}{
		{
			label: "Merge to main",
			desc:  "Goal achieved! Merge branch and complete session",
			emoji: EmojiSuccess,
		},
		{
			label: "Continue tomorrow",
			desc:  "Save progress, keep branch, resume later",
			emoji: EmojiPin,
		},
		{
			label: "Abandon",
			desc:  "This was a rabbit hole, delete branch",
			emoji: EmojiTrash,
		},
	}

	for i, opt := range options {
		cursor := "  "
		style := lipgloss.NewStyle()

		if i == m.selected {
			cursor = "▸ "
			style = style.Foreground(ColorPrimary).Bold(true)
		}

		line := fmt.Sprintf("%s%s %s", cursor, opt.emoji, opt.label)
		b.WriteString(style.Render(line))
		b.WriteString("\n")

		if i == m.selected {
			desc := MutedStyle.Render(fmt.Sprintf("   %s", opt.desc))
			b.WriteString(desc)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(HintStyle.Render("↑/↓ to select • Enter to confirm • Esc to cancel"))

	return BaseStyle.Render(b.String())
}

func (m EndModel) renderSummary() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s  %s\n", EmojiGoal, m.session.Task))
	b.WriteString("\n")

	elapsedStr := formatDuration(m.elapsed)
	b.WriteString(fmt.Sprintf("%s  Time: %s (planned: %s)\n", EmojiTime, elapsedStr, m.session.TimeBox))
	b.WriteString(fmt.Sprintf("%s  Commits: %d\n", EmojiCommit, m.commits))
	b.WriteString(fmt.Sprintf("%s  Drifts: %d\n", EmojiDrift, len(m.session.Drifts)))

	return b.String()
}

func (m EndModel) renderDriftLog() string {
	var b strings.Builder

	b.WriteString(WarningStyle.Render(fmt.Sprintf("%s Drift Log:", EmojiDrift)))
	b.WriteString("\n\n")

	for i, drift := range m.session.Drifts {
		timeStr := drift.Timestamp.Format("15:04")
		b.WriteString(fmt.Sprintf("  %d. [%s] %s", i+1, timeStr, drift.Description))
		if drift.Reason != "" {
			b.WriteString("\n")
			b.WriteString(MutedStyle.Render(fmt.Sprintf("      Reason: %s", drift.Reason)))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m EndModel) renderConfirmation() string {
	var message string

	switch m.choice {
	case actionMerge:
		message = SuccessStyle.Render(fmt.Sprintf("%s Merging to main...", EmojiSuccess))
	case actionContinue:
		message = InfoStyle.Render(fmt.Sprintf("%s Session saved for later", EmojiPin))
	case actionAbandon:
		message = WarningStyle.Render(fmt.Sprintf("%s Abandoning branch...", EmojiTrash))
	}

	return BaseStyle.Render(message)
}

func (m EndModel) HandleAction() error {
	switch m.choice {
	case actionMerge:
		// Merge to main
		if err := git.MergeToMain(m.session.Task); err != nil {
			return fmt.Errorf("failed to merge: %w", err)
		}
		if err := session.Clear(); err != nil {
			return fmt.Errorf("failed to clear session: %w", err)
		}
		fmt.Println(SuccessStyle.Render(fmt.Sprintf("\n%s Session complete! Branch merged to main.", EmojiSuccess)))

	case actionContinue:
		// Just exit, keep everything
		fmt.Println(InfoStyle.Render(fmt.Sprintf("\n%s Session saved. Run 'focus status' to resume.", EmojiPin)))

	case actionAbandon:
		// Delete branch and clear session
		if err := git.DeleteBranch(); err != nil {
			return fmt.Errorf("failed to delete branch: %w", err)
		}
		if err := session.Clear(); err != nil {
			return fmt.Errorf("failed to clear session: %w", err)
		}
		fmt.Println(WarningStyle.Render(fmt.Sprintf("\n%s Branch abandoned. Back to main.", EmojiTrash)))
	}

	return nil
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
