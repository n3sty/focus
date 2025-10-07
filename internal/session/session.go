package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Session represents an active focus session
type Session struct {
	Task      string    `json:"task"`
	StartTime time.Time `json:"start_time"`
	TimeBox   string    `json:"timebox"`
	Branch    string    `json:"branch"`
	Drifts    []Drift   `json:"drifts"`
}

// Drift represents a moment when the user went off-track
type Drift struct {
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	Reason      string    `json:"reason,omitempty"`
}

const focusDir = ".focus"
const sessionFile = "session.json"

// Load reads the current session from .focus/session.json
func Load() (*Session, error) {
	path := filepath.Join(focusDir, sessionFile)

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

// Save writes the session to .focus/session.json
func (s *Session) Save() error {
	if err := os.MkdirAll(focusDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(focusDir, sessionFile)
	return os.WriteFile(path, data, 0644)
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

// Exists checks if a focus session exists
func Exists() bool {
	path := filepath.Join(focusDir, sessionFile)
	_, err := os.Stat(path)
	return err == nil
}

// Clear removes the focus session
func Clear() error {
	return os.RemoveAll(focusDir)
}
