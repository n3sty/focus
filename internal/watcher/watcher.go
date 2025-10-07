package watcher

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n3sty/focus/internal/daemon"
	"github.com/n3sty/focus/internal/notify"
	"github.com/n3sty/focus/internal/session"
)

// Config holds watcher configuration
type Config struct {
	CheckInterval    time.Duration // How often to check session state
	ReminderInterval time.Duration // How often to send "focus check" reminders
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		CheckInterval:    30 * time.Second, // Check every 30s
		ReminderInterval: 25 * time.Minute, // Reminder every 25 min (Pomodoro)
	}
}

// Watch starts watching the focus session
func Watch(cfg Config) error {
	// Write PID file
	if err := daemon.WritePID(); err != nil {
		return fmt.Errorf("failed to write PID: %w", err)
	}
	defer daemon.CleanPID()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Tracking state
	lastReminder := time.Now()
	timeboxExpiredNotified := false

	ticker := time.NewTicker(cfg.CheckInterval)
	defer ticker.Stop()

	fmt.Println("🔍 Focus watcher started (running in background)")

	for {
		select {
		case <-ticker.C:
			// Load session
			sess, err := session.Load()
			if err != nil {
				// Session doesn't exist, stop watching
				return nil
			}

			now := time.Now()
			elapsed := now.Sub(sess.StartTime)

			// Parse timebox duration
			timeboxDuration, err := parseTimebox(sess.TimeBox)
			if err != nil {
				continue
			}

			// Check if timebox expired
			if elapsed >= timeboxDuration && !timeboxExpiredNotified {
				notify.SendUrgent(
					"⏱️ Focus Timebox Expired!",
					fmt.Sprintf("Your %s timebox for '%s' has ended. Run 'focus check' or 'focus end'", sess.TimeBox, sess.Task),
				)
				timeboxExpiredNotified = true
			}

			// Send periodic reminders
			if now.Sub(lastReminder) >= cfg.ReminderInterval && !timeboxExpiredNotified {
				notify.Send(
					"🎯 Focus Check",
					fmt.Sprintf("Still working on: %s? Run 'focus check'", sess.Task),
				)
				lastReminder = now
			}

		case <-sigChan:
			fmt.Println("🛑 Focus watcher stopped")
			return nil
		}
	}
}

// parseTimebox converts "3h", "90m", "2h30m" to time.Duration
func parseTimebox(timebox string) (time.Duration, error) {
	return time.ParseDuration(timebox)
}
