package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

const pidFile = ".focus/daemon.pid"

// IsRunning checks if the daemon is currently running
func IsRunning() bool {
	pid, err := ReadPID()
	if err != nil {
		return false
	}

	// Check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process is alive
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// WritePID writes the current process ID to the PID file
func WritePID() error {
	pid := os.Getpid()
	pidPath := filepath.Join(pidFile)

	// Ensure directory exists
	dir := filepath.Dir(pidPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(pidPath, []byte(fmt.Sprintf("%d", pid)), 0644)
}

// ReadPID reads the process ID from the PID file
func ReadPID() (int, error) {
	pidPath := filepath.Join(pidFile)
	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return pid, nil
}

// Stop stops the daemon by sending SIGTERM
func Stop() error {
	pid, err := ReadPID()
	if err != nil {
		return fmt.Errorf("daemon not running")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	// Remove PID file
	return os.Remove(pidFile)
}

// CleanPID removes the PID file (call this on daemon exit)
func CleanPID() error {
	return os.Remove(pidFile)
}
