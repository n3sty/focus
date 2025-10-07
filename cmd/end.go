package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/n3sty/focus/internal/daemon"
	"github.com/n3sty/focus/internal/session"
	"github.com/n3sty/focus/internal/tui"
	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "End the current focus session",
	Long: `End your current focus session with an interactive review.

You'll be able to:
- Review what you accomplished
- Merge to main if goal is complete
- Continue tomorrow if still in progress
- Abandon the branch if it was a rabbit hole`,
	RunE: runEnd,
}

func init() {
	rootCmd.AddCommand(endCmd)
}

func runEnd(cmd *cobra.Command, args []string) error {
	// Load session
	sess, err := session.Load()
	if err != nil {
		return fmt.Errorf("❌ No active focus session. Run 'focus start' to begin")
	}

	// Launch TUI
	model := tui.NewEndModel(sess)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	// Handle the user's choice
	if m, ok := finalModel.(tui.EndModel); ok {
		if err := m.HandleAction(); err != nil {
			return err
		}
	}

	// Stop watcher when session ends (for merge and abandon, not continue)
	if m, ok := finalModel.(tui.EndModel); ok {
		if m.GetChoice() != 1 { // 1 = continue
			if daemon.IsRunning() {
				if err := daemon.Stop(); err != nil {
					fmt.Printf("⚠️  Warning: Could not stop watcher: %v\n", err)
				} else {
					fmt.Println("✓ Background watcher stopped")
				}
			}
		}
	}

	return nil
}
