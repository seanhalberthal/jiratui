package wikiview

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/seanhalberthal/jiru/internal/confluence"
)

func TestDismissed_SentinelResetsOnRead(t *testing.T) {
	m := New()
	m.dismissed = true
	if !m.Dismissed() {
		t.Error("expected true on first read")
	}
	if m.Dismissed() {
		t.Error("expected false on second read")
	}
}

func TestOpenURL_SentinelResetsOnRead(t *testing.T) {
	m := New()
	m.openURL = "https://example.com"
	url, ok := m.OpenURL()
	if !ok || url != "https://example.com" {
		t.Errorf("first read = (%q, %v), want (url, true)", url, ok)
	}
	_, ok = m.OpenURL()
	if ok {
		t.Error("expected false on second read")
	}
}

func TestSelectedIssue_SentinelResetsOnRead(t *testing.T) {
	m := New()
	m.selIssue = "PROJ-123"
	key, ok := m.SelectedIssue()
	if !ok || key != "PROJ-123" {
		t.Errorf("first read = (%q, %v), want (PROJ-123, true)", key, ok)
	}
	_, ok = m.SelectedIssue()
	if ok {
		t.Error("expected false on second read")
	}
}

func TestSelectedAncestor_SentinelResetsOnRead(t *testing.T) {
	m := New()
	m.selAnc = "456"
	id, ok := m.SelectedAncestor()
	if !ok || id != "456" {
		t.Errorf("first read = (%q, %v), want (456, true)", id, ok)
	}
	_, ok = m.SelectedAncestor()
	if ok {
		t.Error("expected false on second read")
	}
}

func TestSetPage_ViewNonEmpty(t *testing.T) {
	m := New()
	m = m.SetSize(80, 24)
	m.SetPage(&confluence.Page{
		ID:      "123",
		Title:   "Test Page",
		BodyADF: `{"type":"doc","version":1,"content":[{"type":"paragraph","content":[{"type":"text","text":"Hello"}]}]}`,
	})
	view := m.View()
	if view == "" {
		t.Error("View() should not be empty after SetPage")
	}
}

func TestCurrentPage(t *testing.T) {
	m := New()
	if m.CurrentPage() != nil {
		t.Error("CurrentPage should be nil initially")
	}
	page := &confluence.Page{ID: "1", Title: "P"}
	m.SetPage(page)
	if m.CurrentPage() != page {
		t.Error("CurrentPage should return the set page")
	}
}

func TestUpdate_EscDismisses(t *testing.T) {
	m := New()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !m.Dismissed() {
		t.Error("Esc should set dismissed")
	}
}

func TestRenderBreadcrumb_PersonalSpaceKeyOmitted(t *testing.T) {
	m := New()
	m = m.SetSize(80, 24)
	m.SetSpaceKey("~user1")
	m.SetPage(&confluence.Page{ID: "1", Title: "Page", BodyADF: `{"type":"doc","version":1,"content":[]}`})
	bc := m.renderBreadcrumb()
	if bc != "" {
		// Personal space keys (starting with ~) should not appear in breadcrumb
		// and with no ancestors, breadcrumb should be empty
		t.Errorf("breadcrumb = %q, want empty for personal space with no ancestors", bc)
	}
}
