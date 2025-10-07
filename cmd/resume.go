package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/n3sty/focus/internal/session"
	"github.com/n3sty/focus/internal/tui"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused focus session",
	Long:  `Lists all paused sessions and allows you to select one to resume.`,
	RunE:  runResume,
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func runResume(cmd *cobra.Command, args []string) error {
	// Get all paused sessions
	sessions, err := session.ListPaused()
	if err != nil {
		return fmt.Errorf("failed to list paused sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No paused sessions found")
		return nil
	}

	// If only one paused session, resume it automatically
	if len(sessions) == 1 {
		sess := sessions[0]
		if err := sess.Activate(); err != nil {
			return fmt.Errorf("failed to resume session: %w", err)
		}

		fmt.Println("\nSession resumed:")
		fmt.Printf("  Goal: %s\n", sess.Task)
		fmt.Printf("  Branch: %s\n", sess.Branch)
		return nil
	}

	// Multiple sessions - show TUI selector
	model := tui.NewResumeModel(sessions)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	// Get the user's choice
	if m, ok := finalModel.(tui.ResumeModel); ok {
		chosen := m.GetChoice()
		if chosen == nil {
			fmt.Println("Resume cancelled")
			return nil
		}

		if err := chosen.Activate(); err != nil {
			return fmt.Errorf("failed to resume session: %w", err)
		}

		fmt.Println("\nSession resumed:")
		fmt.Printf("  Goal: %s\n", chosen.Task)
		fmt.Printf("  Branch: %s\n", chosen.Branch)
	}

	return nil
}
