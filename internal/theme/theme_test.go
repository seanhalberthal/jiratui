package theme

import "testing"

func TestStatusCategory(t *testing.T) {
	tests := []struct {
		status string
		want   int
	}{
		{"Done", 2},
		{"Closed", 2},
		{"Resolved", 2},
		{"In Progress", 1},
		{"In Review", 1},
		{"To Do", 0},
		{"Open", 0},
		{"Backlog", 0},
		{"Unknown Status", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := StatusCategory(tt.status)
			if got != tt.want {
				t.Errorf("StatusCategory(%q) = %d, want %d", tt.status, got, tt.want)
			}
		})
	}
}

func TestStatusStyle_ReturnsNonNil(t *testing.T) {
	statuses := []string{"Done", "In Progress", "To Do", "Unknown"}
	for _, s := range statuses {
		t.Run(s, func(t *testing.T) {
			style := StatusStyle(s)
			// Verify the style can render without panic.
			_ = style.Render("test")
		})
	}
}

func TestStatusStyle_AllKnownStatuses(t *testing.T) {
	statuses := []string{
		"To Do", "Open", "Backlog", "In Progress", "In Review",
		"Done", "Closed", "Resolved",
		"Unknown Status", "", // should return a default style
	}
	for _, s := range statuses {
		t.Run(s, func(t *testing.T) {
			style := StatusStyle(s)
			_ = style.Render(s) // must not panic
		})
	}
}

func TestRenderLogo_NarrowTerminal(t *testing.T) {
	result := RenderLogo(10)
	if result != "" {
		t.Errorf("expected empty logo for narrow terminal, got non-empty (%d bytes)", len(result))
	}
}

func TestRenderLogo_ExactWidth(t *testing.T) {
	result := RenderLogo(LogoWidth)
	if result == "" {
		t.Error("expected non-empty logo at exact LogoWidth")
	}
}

func TestRenderLogo_BelowWidth(t *testing.T) {
	result := RenderLogo(LogoWidth - 1)
	if result != "" {
		t.Error("expected empty logo just below LogoWidth")
	}
}

func TestRenderLogo_WideTerminal(t *testing.T) {
	result := RenderLogo(120)
	if result == "" {
		t.Error("expected non-empty logo for wide terminal")
	}
}

func TestStatusCategory_AllBranches(t *testing.T) {
	// Verify coverage of all three category branches.
	if StatusCategory("Done") != 2 {
		t.Error("Done should be category 2")
	}
	if StatusCategory("In Progress") != 1 {
		t.Error("In Progress should be category 1")
	}
	if StatusCategory("Custom") != 0 {
		t.Error("Custom should be category 0")
	}
}

func TestStatusStyle_Categories(t *testing.T) {
	// Done statuses should use StyleStatusDone.
	for _, s := range []string{"Done", "Closed", "Resolved"} {
		got := StatusStyle(s)
		if got.GetForeground() != StyleStatusDone.GetForeground() {
			t.Errorf("StatusStyle(%q) should match StyleStatusDone", s)
		}
	}

	// In Progress statuses should use StyleStatusInProgress.
	for _, s := range []string{"In Progress", "In Review"} {
		got := StatusStyle(s)
		if got.GetForeground() != StyleStatusInProgress.GetForeground() {
			t.Errorf("StatusStyle(%q) should match StyleStatusInProgress", s)
		}
	}
}
