package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CreateFocusBranch creates a new git branch for the focus session
func CreateFocusBranch(task string) (string, error) {
	// Generate branch name from task
	branchName := fmt.Sprintf("focus/%s", slugify(task))

	// Create and checkout branch
	cmd := exec.Command("git", "checkout", "-b", branchName)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create branch: %w", err)
	}

	// Make empty commit to mark start
	commitMsg := fmt.Sprintf("🎯 START: %s", task)
	cmd = exec.Command("git", "commit", "--allow-empty", "-m", commitMsg)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to create start commit: %w", err)
	}

	return branchName, nil
}

// GetCommitsSince returns the number of commits since a given time
func GetCommitsSince(since time.Time) (int, error) {
	sinceStr := since.Format("2006-01-02T15:04:05")
	cmd := exec.Command("git", "rev-list", "--count", "--since="+sinceStr, "HEAD")

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to count commits: %w", err)
	}

	var count int
	fmt.Sscanf(string(output), "%d", &count)
	return count, nil
}

// GetCurrentBranch returns the current git branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// MergeToMain merges the current branch to main and deletes the focus branch
func MergeToMain(task string) error {
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	// Checkout main
	cmd := exec.Command("git", "checkout", "main")
	if err := cmd.Run(); err != nil {
		// Try master if main doesn't exist
		cmd = exec.Command("git", "checkout", "master")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout main/master: %w", err)
		}
	}

	// Merge with no-ff to preserve history
	commitMsg := fmt.Sprintf("✅ Completed: %s", task)
	cmd = exec.Command("git", "merge", "--no-ff", currentBranch, "-m", commitMsg)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to merge: %w", err)
	}

	// Delete the focus branch
	cmd = exec.Command("git", "branch", "-d", currentBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

// DeleteBranch deletes the current branch and returns to main
func DeleteBranch() error {
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	// Checkout main
	cmd := exec.Command("git", "checkout", "main")
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("git", "checkout", "master")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout main/master: %w", err)
		}
	}

	// Force delete the branch
	cmd = exec.Command("git", "branch", "-D", currentBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

// IsGitRepo checks if current directory is a git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// slugify converts a task name to a git-safe branch name
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove special characters
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
