package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/n3sty/focus/internal/session"
)

type checkState int

const (
	stateQuestion checkState = iota
	stateDriftDescription
	stateDriftReason
	stateComplete
)

type CheckModel struct {
	session   *session.Session
	state     checkState
	textarea  textarea.Model
	viewport  viewport.Model
	Updated   bool
	stillOnTrack bool
	driftDesc string
	driftReason string
	width     int
	height    int
}

func NewCheckModel(sess *session.Session) CheckModel {
	ta := textarea.New()
	ta.Placeholder = "Type here..."
	ta.Focus()
	ta.CharLimit = 200
	ta.SetWidth(60)
	ta.SetHeight(3)

	vp := viewport.New(80, 20)

	return CheckModel{
		session:  sess,
		state:    stateQuestion,
		textarea: ta,
		viewport: vp,
		Updated:  false,
	}
}

func (m CheckModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m CheckModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			return m.handleEnter()

		default:
			// Handle state-specific key bindings
			switch m.state {
			case stateQuestion:
				return m.handleQuestionKeys(msg)
			case stateDriftDescription, stateDriftReason:
				m.textarea, cmd = m.textarea.Update(msg)
				return m, cmd
			case stateComplete:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m CheckModel) handleQuestionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.stillOnTrack = true
		m.state = stateComplete
		return m, tea.Quit
	case "n", "N":
		m.stillOnTrack = false
		m.state = stateDriftDescription
		m.textarea.Reset()
		m.textarea.Placeholder = "What are you working on instead?"
		return m, nil
	case "d", "D":
		// Defer - skip for now without logging
		m.state = stateComplete
		return m, tea.Quit
	}
	return m, nil
}

func (m CheckModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateDriftDescription:
		m.driftDesc = strings.TrimSpace(m.textarea.Value())
		if m.driftDesc == "" {
			return m, nil
		}
		m.state = stateDriftReason
		m.textarea.Reset()
		m.textarea.Placeholder = "Why is this necessary? (optional, press Enter to skip)"
		return m, nil

	case stateDriftReason:
		m.driftReason = strings.TrimSpace(m.textarea.Value())
		m.session.AddDrift(m.driftDesc, m.driftReason)
		m.Updated = true
		m.state = stateComplete
		return m, tea.Quit
	}
	return m, nil
}

func (m CheckModel) View() string {
	var b strings.Builder

	// Header
	header := TitleStyle.Render("🎯 Focus Check")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Current goal box
	goalBox := BoxStyle.Render(fmt.Sprintf("%s Current Goal\n\n%s", EmojiGoal, m.session.Task))
	b.WriteString(goalBox)
	b.WriteString("\n\n")

	// State-specific content
	switch m.state {
	case stateQuestion:
		b.WriteString(m.renderQuestion())
	case stateDriftDescription:
		b.WriteString(m.renderDriftDescription())
	case stateDriftReason:
		b.WriteString(m.renderDriftReason())
	case stateComplete:
		b.WriteString(m.renderComplete())
	}

	return BaseStyle.Render(b.String())
}

func (m CheckModel) renderQuestion() string {
	var b strings.Builder

	question := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorInfo).
		Render("Are you still working on this goal?")

	b.WriteString(question)
	b.WriteString("\n\n")

	options := []string{
		SuccessStyle.Render("[y] Yes") + " - Still on track!",
		WarningStyle.Render("[n] No") + "  - I've drifted...",
		MutedStyle.Render("[d] Defer") + " - Check me later",
	}

	b.WriteString(strings.Join(options, "\n"))
	b.WriteString("\n\n")
	b.WriteString(HintStyle.Render("Press Esc to cancel"))

	return b.String()
}

func (m CheckModel) renderDriftDescription() string {
	var b strings.Builder

	prompt := lipgloss.NewStyle().
		Foreground(ColorWarning).
		Render(fmt.Sprintf("%s What are you actually working on?", EmojiDrift))

	b.WriteString(prompt)
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())
	b.WriteString("\n\n")
	b.WriteString(HintStyle.Render("Press Enter to continue, Esc to cancel"))

	return b.String()
}

func (m CheckModel) renderDriftReason() string {
	var b strings.Builder

	prompt := lipgloss.NewStyle().
		Foreground(ColorInfo).
		Render(fmt.Sprintf("%s Why is this necessary?", EmojiThink))

	b.WriteString(prompt)
	b.WriteString("\n\n")
	b.WriteString(MutedStyle.Render("This helps you reflect on whether it's a valid detour or scope creep."))
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())
	b.WriteString("\n\n")
	b.WriteString(HintStyle.Render("Press Enter to save, Esc to cancel"))

	return b.String()
}

func (m CheckModel) renderComplete() string {
	if m.stillOnTrack {
		return SuccessStyle.Render(fmt.Sprintf("%s Great! Keep going.", EmojiSuccess))
	}

	var b strings.Builder
	b.WriteString(WarningStyle.Render(fmt.Sprintf("%s Drift logged.", EmojiDrift)))
	b.WriteString("\n\n")
	b.WriteString(InfoStyle.Render("Consider:"))
	b.WriteString("\n")
	b.WriteString("  • git stash (save current work)\n")
	b.WriteString("  • git checkout main (return to main branch)\n")
	b.WriteString("  • Or create a new focus session for this work\n")
	return b.String()
}
