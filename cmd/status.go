package cmd

import (
	"fmt"
	"time"

	"github.com/n3sty/focus/internal/git"
	"github.com/n3sty/focus/internal/session"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current focus session status",
	Long:  `Display information about your current focus session including time elapsed, commits made, and any drifts.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Load session
	sess, err := session.Load()
	if err != nil {
		return fmt.Errorf("❌ No active focus session. Run 'focus start' to begin")
	}

	// Get commit count
	commits, err := git.GetCommitsSince(sess.StartTime)
	if err != nil {
		commits = 0 // Non-fatal, just show 0
	}

	// Calculate elapsed time
	elapsed := time.Since(sess.StartTime)
	elapsedStr := formatDuration(elapsed)

	// Display status
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Focus Session Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Goal:     %s\n", sess.Task)
	fmt.Printf("Started:  %s\n", sess.StartTime.Format("15:04 PM"))
	fmt.Printf("Elapsed:  %s\n", elapsedStr)
	fmt.Printf("Timebox:  %s\n", sess.TimeBox)
	fmt.Printf("Branch:   %s\n", sess.Branch)
	fmt.Printf("Commits:  %d\n", commits)
	fmt.Printf("Drifts:   %d\n", len(sess.Drifts))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Show drifts if any
	if len(sess.Drifts) > 0 {
		fmt.Println("\n🐰 Drift Log:")
		for i, drift := range sess.Drifts {
			timeStr := drift.Timestamp.Format("15:04")
			fmt.Printf("  %d. [%s] %s", i+1, timeStr, drift.Description)
			if drift.Reason != "" {
				fmt.Printf(" (Reason: %s)", drift.Reason)
			}
			fmt.Println()
		}
	}

	fmt.Println("\nCommands:")
	fmt.Println("  focus check - Verify you're still on track")
	fmt.Println("  focus end   - Complete or abandon this session\n")

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
