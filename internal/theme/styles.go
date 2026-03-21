package theme

import "github.com/charmbracelet/lipgloss"

// Shared rendering styles used by both ADF and wiki markup renderers.
var (
	StyleBold          = lipgloss.NewStyle().Bold(true)
	StyleItalic        = lipgloss.NewStyle().Italic(true)
	StyleUnderline     = lipgloss.NewStyle().Underline(true)
	StyleStrikethrough = lipgloss.NewStyle().Strikethrough(true)
	StyleCode          = lipgloss.NewStyle().Foreground(ColourWarning)
	StyleLink          = lipgloss.NewStyle().Foreground(ColourPrimary).Underline(true)
	StyleLinkURL       = lipgloss.NewStyle().Foreground(ColourSubtle)
	StyleHeading       = lipgloss.NewStyle().Bold(true).Foreground(ColourPrimary)
	StyleBlockquote    = lipgloss.NewStyle().Foreground(ColourSubtle).Italic(true)
	StyleCodeBlock     = lipgloss.NewStyle().Foreground(ColourWarning)
	StyleHRule         = lipgloss.NewStyle().Foreground(ColourSubtle)
	StyleBullet        = lipgloss.NewStyle().Foreground(ColourPrimary).Bold(true)
	StyleImage         = lipgloss.NewStyle().Foreground(ColourSubtle).Italic(true)
)
