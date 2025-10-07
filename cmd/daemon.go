package cmd

import (
	"fmt"

	"github.com/n3sty/focus/internal/daemon"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the background watcher daemon",
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if the watcher daemon is running",
	RunE:  runDaemonStatus,
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the watcher daemon",
	RunE:  runDaemonStop,
}

func init() {
	daemonCmd.AddCommand(daemonStatusCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	rootCmd.AddCommand(daemonCmd)
}

func runDaemonStatus(cmd *cobra.Command, args []string) error {
	if daemon.IsRunning() {
		pid, _ := daemon.ReadPID()
		fmt.Printf("✓ Watcher daemon is running (PID: %d)\n", pid)
	} else {
		fmt.Println("✗ Watcher daemon is not running")
	}
	return nil
}

func runDaemonStop(cmd *cobra.Command, args []string) error {
	if !daemon.IsRunning() {
		fmt.Println("✗ Watcher daemon is not running")
		return nil
	}

	if err := daemon.Stop(); err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}

	fmt.Println("✓ Watcher daemon stopped")
	return nil
}
