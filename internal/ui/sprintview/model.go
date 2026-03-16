package sprintview

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/seanhalberthal/jiru/internal/jira"
	"github.com/seanhalberthal/jiru/internal/theme"
	"github.com/seanhalberthal/jiru/internal/ui/issuedelegate"
)

// Model is the sprint issue list view.
type Model struct {
	list      list.Model
	issues    []jira.Issue
	listStale bool // True when issues have been updated but not synced to the list widget (e.g. during filtering).
	width     int
	height    int
	selected  *jira.Issue // set when user presses enter.
	openKeys  key.Binding
}

// New creates a new sprint view model.
func New() Model {
	delegate := issuedelegate.Delegate{}
	l := list.New(nil, delegate, 0, 0)
	l.Title = "Issues"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false) // We handle help ourselves.
	l.Styles.Title = theme.StyleTitle

	return Model{
		list: l,
		openKeys: key.NewBinding(
			key.WithKeys("enter"),
		),
	}
}

// SetSize updates the dimensions.
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
	return m
}

// SetIssues populates the list with issues.
func (m Model) SetIssues(issues []jira.Issue) Model {
	m.issues = issues
	m.list.SetItems(issuedelegate.ToItems(issues))
	m.list.Title = fmt.Sprintf("Issues (%d)", len(issues))
	return m
}

// AppendIssues adds more issues to the existing list (for progressive pagination).
// Deduplicates by issue key to handle overlapping API pages.
// When a filter is active, items are buffered to avoid disrupting the user's
// filtering interaction — they are flushed when the filter is cleared.
func (m Model) AppendIssues(issues []jira.Issue) Model {
	seen := make(map[string]bool, len(m.issues))
	for _, iss := range m.issues {
		seen[iss.Key] = true
	}
	for _, iss := range issues {
		if !seen[iss.Key] {
			m.issues = append(m.issues, iss)
			seen[iss.Key] = true
		}
	}

	if m.list.FilterState() != list.Unfiltered {
		// Don't call SetItems while the user is filtering — just mark as stale.
		m.listStale = true
		m.list.Title = fmt.Sprintf("Issues (%d) loading...", len(m.issues))
		return m
	}

	items := issuedelegate.ToItems(m.issues)
	m.list.SetItems(items)
	m.list.Title = fmt.Sprintf("Issues (%d)", len(m.issues))
	return m
}

// SetLoading updates the title to show pagination progress.
func (m Model) SetLoading(loading bool) Model {
	if loading {
		m.list.Title = fmt.Sprintf("Issues (%d) loading...", len(m.issues))
	} else {
		m.list.Title = fmt.Sprintf("Issues (%d)", len(m.issues))
	}
	return m
}

// SelectedIssue returns the issue the user selected (if any) and resets the selection.
func (m *Model) SelectedIssue() (jira.Issue, bool) {
	if m.selected == nil {
		return jira.Issue{}, false
	}
	iss := *m.selected
	m.selected = nil
	return iss, true
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't handle keys when filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		if key.Matches(msg, m.openKeys) {
			if item, ok := m.list.SelectedItem().(issueItem); ok {
				m.selected = &item.Issue
				return m, nil
			}
		}

		// d/u for half-page scrolling (forwarded as pgdown/pgup to the list).
		switch msg.String() {
		case "d":
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(tea.KeyMsg{Type: tea.KeyPgDown})
			return m, cmd
		case "u":
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(tea.KeyMsg{Type: tea.KeyPgUp})
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Flush buffered items once the user clears the filter.
	if m.listStale && m.list.FilterState() == list.Unfiltered {
		m.listStale = false
		items := issuedelegate.ToItems(m.issues)
		m.list.SetItems(items)
		m.list.Title = fmt.Sprintf("Issues (%d)", len(m.issues))
	}

	return m, cmd
}

// Filtering returns true when the list filter input is active.
func (m Model) Filtering() bool {
	return m.list.FilterState() == list.Filtering
}

// View renders the sprint view.
func (m Model) View() string {
	return m.list.View()
}
