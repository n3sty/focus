package tui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/n3sty/focus/internal/session"
)

// sessionItem wraps a session for the list
type sessionItem struct {
	session *session.Session
}

func (i sessionItem) FilterValue() string { return i.session.Task }

// Custom item delegate for rendering
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 3 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(sessionItem)
	if !ok {
		return
	}

	sess := i.session
	elapsed := time.Since(sess.StartTime)
	elapsedStr := formatDuration(elapsed)

	// Render the session info
	str := fmt.Sprintf("%s", sess.Task)
	desc := fmt.Sprintf("Started: %s | Elapsed: %s | Branch: %s",
		sess.StartTime.Format("15:04 PM"),
		elapsedStr,
		sess.Branch,
	)

	isSelected := index == m.Index()

	if isSelected {
		str = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Render("> " + str)
		desc = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Render("  " + desc)
	} else {
		str = "  " + str
		desc = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Render("  " + desc)
	}

	fmt.Fprintf(w, "%s\n%s", str, desc)
}

type ResumeModel struct {
	list     list.Model
	choice   *session.Session
	quitting bool
}

func NewResumeModel(sessions []*session.Session) ResumeModel {
	items := make([]list.Item, len(sessions))
	for i, s := range sessions {
		items[i] = sessionItem{session: s}
	}

	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = "Paused Sessions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle

	return ResumeModel{
		list: l,
	}
}

func (m ResumeModel) Init() tea.Cmd {
	return nil
}

func (m ResumeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			if i, ok := m.list.SelectedItem().(sessionItem); ok {
				m.choice = i.session
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ResumeModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.list.View())
	b.WriteString("\n\n")
	b.WriteString(HintStyle.Render("↑/↓ to select • Enter to resume • Esc to cancel"))

	return b.String()
}

func (m ResumeModel) GetChoice() *session.Session {
	return m.choice
}
