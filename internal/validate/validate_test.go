package validate

import (
	"testing"
)

func TestIssueKey(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{"PROJ-123", true},
		{"AB-1", true},
		{"LONGPROJ-99999", true},
		{"A-1", true},
		{"proj-123", false},                   // lowercase
		{"PROJ123", false},                    // missing dash
		{"PROJ-", false},                      // missing number
		{"-123", false},                       // missing project
		{"", false},                           // empty
		{"PROJ-0", true},                      // zero is valid
		{"PROJ-123 OR project = EVIL", false}, // injection attempt
		{"PROJ-123)", false},                  // JQL metacharacter
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			err := IssueKey(tt.key)
			if tt.valid && err != nil {
				t.Errorf("IssueKey(%q) returned error: %v", tt.key, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("IssueKey(%q) expected error, got nil", tt.key)
			}
		})
	}
}

func TestProjectKey(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{"PROJ", true},
		{"A", true},
		{"AB_CD", true},
		{"ABCDEFGHIJ", true},   // 10 chars (max)
		{"ABCDEFGHIJK", false}, // 11 chars (too long)
		{"proj", false},        // lowercase
		{"1PROJ", false},       // starts with digit
		{"", false},
		{"PROJ 123", false},               // space
		{"PROJ OR project = EVIL", false}, // injection
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			err := ProjectKey(tt.key)
			if tt.valid && err != nil {
				t.Errorf("ProjectKey(%q) returned error: %v", tt.key, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("ProjectKey(%q) expected error, got nil", tt.key)
			}
		})
	}
}

func TestDomain(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"myco.atlassian.net", true},
		{"sub.domain.example.com", true},
		{"a.b.c", true},
		{"example.com", false},                // only two segments
		{"", false},                           // empty
		{"https://myco.atlassian.net", false}, // protocol prefix
		{"-bad.domain.com", false},            // leading hyphen
		{"domain.com/path", false},            // path component
		{"my co.atlassian.net", false},        // space
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := Domain(tt.input)
			if (err == nil) != tt.valid {
				t.Errorf("Domain(%q) valid=%v, want %v (err=%v)", tt.input, err == nil, tt.valid, err)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"user@example.com", true},
		{"user+tag@sub.domain.co.uk", true},
		{"first.last@company.com", true},
		{"", false},
		{"noatsign", false},
		{"@nodomain.com", false},
		{"user@", false},
		{"user@.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := Email(tt.input)
			if (err == nil) != tt.valid {
				t.Errorf("Email(%q) valid=%v, want %v (err=%v)", tt.input, err == nil, tt.valid, err)
			}
		})
	}
}

func TestAuthType(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"basic", true},
		{"bearer", true},
		{"", false},
		{"oauth", false},
		{"Basic", false},
		{"BEARER", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := AuthType(tt.input)
			if (err == nil) != tt.valid {
				t.Errorf("AuthType(%q) valid=%v, want %v (err=%v)", tt.input, err == nil, tt.valid, err)
			}
		})
	}
}

func TestBoardID(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"", false},
		{"42", false},
		{"1", false},
		{"0", true},
		{"-1", true},
		{"abc", true},
		{"3.14", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := BoardID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoardID(%q) err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestBranchName(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"feat/my-branch", false},
		{"PROJ-123-fix-login", false},
		{"simple", false},
		{"", true},
		{"-flag", true},
		{".hidden", true},
		{"name.lock", true},
		{"has..double", true},
		{"has@{ref", true},
		{"@", true},
		{"has space", true},
		{"has~tilde", true},
		{"has^caret", true},
		{"has:colon", true},
		{"has?question", true},
		{"has*star", true},
		{"has[bracket", true},
		{"has\\backslash", true},
		{"trailing.", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := BranchName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BranchName(%q) err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			}
		})
	}
}
