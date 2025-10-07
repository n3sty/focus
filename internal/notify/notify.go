package notify

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Send sends a desktop notification
func Send(title, message string) error {
	switch runtime.GOOS {
	case "darwin":
		return sendMacOS(title, message)
	case "linux":
		return sendLinux(title, message)
	default:
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}

// SendUrgent sends an urgent notification with sound
func SendUrgent(title, message string) error {
	switch runtime.GOOS {
	case "darwin":
		return sendMacOSWithSound(title, message, "Crystal")
	case "linux":
		return sendLinux(title, message)
	default:
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}

// sendMacOS sends notification on macOS using osascript
func sendMacOS(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// sendMacOSWithSound sends notification with sound on macOS
func sendMacOSWithSound(title, message, sound string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "%s"`, message, title, sound)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// sendLinux sends notification on Linux using notify-send
func sendLinux(title, message string) error {
	cmd := exec.Command("notify-send", title, message)
	return cmd.Run()
}
