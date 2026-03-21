package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	issueKeyRe   = regexp.MustCompile(`^[A-Z][A-Z0-9]*-[0-9]+$`)
	projectKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,9}$`)

	// Patterns for setup wizard.
	domainRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?){2,}$`)
	emailRe  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// IssueKey validates a Jira issue key (e.g., "PROJ-123").
func IssueKey(key string) error {
	if !issueKeyRe.MatchString(key) {
		return fmt.Errorf("invalid issue key %q: must match [A-Z]+-[0-9]+", key)
	}
	return nil
}

// ProjectKey validates a Jira project key (e.g., "PROJ").
func ProjectKey(key string) error {
	if !projectKeyRe.MatchString(key) {
		return fmt.Errorf("invalid project key %q: must match [A-Z][A-Z0-9_]{0,9}", key)
	}
	return nil
}

// Domain validates a Jira domain (e.g., "mycompany.atlassian.net").
// Accepts bare domains only — no protocol prefix or trailing slash.
func Domain(d string) error {
	if !domainRe.MatchString(d) {
		return fmt.Errorf("invalid domain %q: expected format like mycompany.atlassian.net", d)
	}
	return nil
}

// Email validates an email address (e.g., "user@example.com").
func Email(email string) error {
	if !emailRe.MatchString(email) {
		return fmt.Errorf("invalid email %q: expected format like user@example.com", email)
	}
	return nil
}

// AuthType validates a Jira auth type ("basic" or "bearer").
func AuthType(at string) error {
	switch at {
	case "basic", "bearer":
		return nil
	default:
		return fmt.Errorf("invalid auth type %q: must be 'basic' or 'bearer'", at)
	}
}

// invalidBranchChars matches characters disallowed in git branch names.
var invalidBranchChars = regexp.MustCompile(`[\x00-\x1f\x7f ~^:?*\[\\]`)

// BranchName validates a string as a legal git branch name.
// Enforces rules from git-check-ref-format(1).
func BranchName(name string) error {
	if name == "" {
		return fmt.Errorf("branch name must not be empty")
	}
	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("branch name must not start with '-'")
	}
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("branch name must not start with '.'")
	}
	if strings.HasSuffix(name, ".lock") {
		return fmt.Errorf("branch name must not end with '.lock'")
	}
	if strings.HasSuffix(name, ".") {
		return fmt.Errorf("branch name must not end with '.'")
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("branch name must not contain '..'")
	}
	if strings.Contains(name, "@{") {
		return fmt.Errorf("branch name must not contain '@{'")
	}
	if name == "@" {
		return fmt.Errorf("branch name must not be '@'")
	}
	if invalidBranchChars.MatchString(name) {
		return fmt.Errorf("branch name contains invalid characters")
	}
	return nil
}

// confluenceIDRe matches Confluence numeric IDs.
var confluenceIDRe = regexp.MustCompile(`^[0-9]+$`)

// ConfluenceID validates a Confluence page or space ID (must be numeric).
func ConfluenceID(id string) error {
	if !confluenceIDRe.MatchString(id) {
		return fmt.Errorf("invalid Confluence ID %q: must be numeric", id)
	}
	return nil
}

// BoardID validates a board ID string (positive integer).
func BoardID(id string) error {
	if id == "" {
		return nil // optional field
	}
	n, err := strconv.Atoi(id)
	if err != nil || n <= 0 {
		return fmt.Errorf("invalid board ID %q: must be a positive integer", id)
	}
	return nil
}
