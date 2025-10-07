package cmd

import (
	"fmt"
	"time"

	"github.com/jobsiemerink/focus/internal/git"
	"github.com/jobsiemerink/focus/internal/session"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [task]",
	Short: "Start a new focused work session",
	Long: `Start a new focused work session with a clear goal and timebox.

This will:
- Create a git branch for your focused work
- Track your session in .focus/session.json
- Make an empty commit marking the session start

Example:
  focus start "Fix non-PDF OCR support" --time 3h`,
	Args: cobra.ExactArgs(1),
	RunE: runStart,
}

var timeBox string

func init() {
	startCmd.Flags().StringVarP(&timeBox, "time", "t", "3h", "Timebox duration (e.g., 1h, 90m, 2h30m)")
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	task := args[0]

	// Check if already in a focus session
	if session.Exists() {
		return fmt.Errorf("❌ Already in a focus session. Run 'focus end' first or 'focus status' to see current session")
	}

	// Check if in a git repository
	if !git.IsGitRepo() {
		return fmt.Errorf("❌ Not in a git repository. Focus requires git for branch tracking")
	}

	// Create git branch
	fmt.Printf("🎯 Starting focus session: %s\n", task)
	fmt.Printf("⏱️  Timebox: %s\n\n", timeBox)

	branch, err := git.CreateFocusBranch(task)
	if err != nil {
		return fmt.Errorf("failed to create git branch: %w", err)
	}

	fmt.Printf("✓ Created branch: %s\n", branch)

	// Create session
	sess := &session.Session{
		Task:      task,
		StartTime: time.Now(),
		TimeBox:   timeBox,
		Branch:    branch,
		Drifts:    []session.Drift{},
	}

	if err := sess.Save(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	fmt.Println("✓ Session saved")
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🎯 Focus Session Active\n")
	fmt.Printf("   Goal: %s\n", task)
	fmt.Printf("   Time: %s\n", timeBox)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nUse these commands during your session:")
	fmt.Println("  focus check  - Check if you're still on track")
	fmt.Println("  focus status - See session progress")
	fmt.Println("  focus end    - End the session\n")

	return nil
}
