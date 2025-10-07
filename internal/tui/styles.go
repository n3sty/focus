package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Color palette
	ColorPrimary   = lipgloss.Color("#FF6B6B")
	ColorSuccess   = lipgloss.Color("#51CF66")
	ColorWarning   = lipgloss.Color("#FFD93D")
	ColorInfo      = lipgloss.Color("#74C0FC")
	ColorMuted     = lipgloss.Color("#868E96")
	ColorDanger    = lipgloss.Color("#FF6B6B")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(1, 0)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginTop(1)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	// Warning style
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	// Info style
	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// Muted style
	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	// Key binding hint style
	HintStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true).
			MarginTop(1)

	// Emoji styles for visual feedback
	EmojiGoal    = "🎯"
	EmojiTime    = "⏱️"
	EmojiCommit  = "📝"
	EmojiDrift   = "🐰"
	EmojiSuccess = "✅"
	EmojiWarning = "⚠️"
	EmojiThink   = "💭"
	EmojiTrash   = "🗑️"
	EmojiPin     = "📌"
)
