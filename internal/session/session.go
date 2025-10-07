package session

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/n3sty/focus/internal/git"
)

// Session represents a focus session
type Session struct {
	ID        string    `json:"id"`
	Task      string    `json:"task"`
	StartTime time.Time `json:"start_time"`
	TimeBox   string    `json:"timebox"`
	Branch    string    `json:"branch"`
	Drifts    []Drift   `json:"drifts"`
	Status    string    `json:"status"` // "active" or "paused"
}

// Drift represents a moment when the user went off-track
type Drift struct {
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	Reason      string    `json:"reason,omitempty"`
}

const focusDir = ".focus"
const sessionsDir = ".focus/sessions"
const activeFile = ".focus/active"

// Load reads the currently active session
func Load() (*Session, error) {
	id, err := getActiveID()
	if err != nil {
		return nil, err
	}

	return LoadByID(id)
}

// LoadByID loads a specific session by ID
func LoadByID(id string) (*Session, error) {
	path := filepath.Join(sessionsDir, id+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Save writes the session to disk
func (s *Session) Save() error {
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(sessionsDir, s.ID+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	// If this is the active session, update the active pointer
	if s.Status == "active" {
		return setActiveID(s.ID)
	}

	return nil
}

// AddDrift logs a distraction
func (s *Session) AddDrift(description, reason string) {
	drift := Drift{
		Timestamp:   time.Now(),
		Description: description,
		Reason:      reason,
	}
	s.Drifts = append(s.Drifts, drift)
}

// Pause marks the session as paused
func (s *Session) Pause() error {
	s.Status = "paused"
	return s.Save()
}

// Activate marks the session as active
func (s *Session) Activate() error {
	// Pause any currently active session
	if err := PauseActive(); err != nil && !os.IsNotExist(err) {
		return err
	}

	s.Status = "active"

	// Make sure we're on the right git branch
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if currentBranch != s.Branch {
		cmd := exec.Command("git", "switch", s.Branch)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to switch branch: %w", err)
		}
	}

	return s.Save()
}

// Delete removes the session from disk
func (s *Session) Delete() error {
	path := filepath.Join(sessionsDir, s.ID+".json")
	return os.Remove(path)
}

// Exists checks if an active session exists
func Exists() bool {
	_, err := getActiveID()
	return err == nil
}

// ListPaused returns all paused sessions
func ListPaused() ([]*Session, error) {
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Session{}, nil
		}
		return nil, err
	}

	var paused []*Session
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := entry.Name()[:len(entry.Name())-5]
		sess, err := LoadByID(id)
		if err != nil {
			continue
		}

		if sess.Status == "paused" {
			paused = append(paused, sess)
		}
	}

	return paused, nil
}

// PauseActive pauses the currently active session
func PauseActive() error {
	sess, err := Load()
	if err != nil {
		return err
	}

	return sess.Pause()
}

// Clear removes all session data
func Clear() error {
	return os.RemoveAll(focusDir)
}

// Helper functions

func getActiveID() (string, error) {
	data, err := os.ReadFile(activeFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func setActiveID(id string) error {
	if err := os.MkdirAll(focusDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(activeFile, []byte(id), 0644)
}

// GenerateID generates a unique session ID
func GenerateID(task string) string {
	// Use timestamp + first few chars of task
	timestamp := time.Now().Unix()
	taskSlug := ""
	for i, r := range task {
		if i >= 8 {
			break
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			taskSlug += string(r)
		} else if r == ' ' {
			taskSlug += "-"
		}
	}
	return fmt.Sprintf("%d-%s", timestamp, taskSlug)
}
