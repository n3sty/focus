package cmd

import (
	"fmt"

	"github.com/n3sty/focus/internal/daemon"
	"github.com/n3sty/focus/internal/watcher"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:    "watch",
	Short:  "Start the background watcher (automatically started by 'focus start')",
	Hidden: true, // Hidden from help - users shouldn't need to run this directly
	RunE:   runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	// Check if already running
	if daemon.IsRunning() {
		return fmt.Errorf("watcher already running")
	}

	// Start watching
	cfg := watcher.DefaultConfig()
	return watcher.Watch(cfg)
}
