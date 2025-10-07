package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "focus",
	Short: "A CLI tool to help you stay focused and avoid scope creep",
	Long: `Focus is a developer productivity tool that helps you:
- Set clear goals with timeboxed sessions
- Track when you drift into rabbit holes
- Make git commits aligned with your actual goals
- Build the discipline to ship instead of polish

Built in public by Job Siemerink.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here if needed
}
