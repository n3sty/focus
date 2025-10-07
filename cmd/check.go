package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/n3sty/focus/internal/session"
	"github.com/n3sty/focus/internal/tui"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if you're still focused on your goal",
	Long: `Opens an interactive prompt to verify you're still working on your stated goal.

If you've drifted, this helps you log the distraction and decide whether to continue or refocus.`,
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	// Load session
	sess, err := session.Load()
	if err != nil {
		return fmt.Errorf("❌ No active focus session. Run 'focus start' to begin")
	}

	// Launch TUI
	model := tui.NewCheckModel(sess)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	// Save any changes
	if m, ok := finalModel.(tui.CheckModel); ok {
		if m.Updated {
			if err := sess.Save(); err != nil {
				return fmt.Errorf("failed to save session: %w", err)
			}
		}
	}

	return nil
}
