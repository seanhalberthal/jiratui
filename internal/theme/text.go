package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// WrapStyledText wraps text at the given width, respecting ANSI escape sequences.
// Delegates to lipgloss.Width for accurate width calculation of styled text.
func WrapStyledText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if lipgloss.Width(line) <= width {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		words := strings.Fields(line)
		current := ""
		for _, word := range words {
			if current == "" {
				current = word
			} else if lipgloss.Width(current+" "+word) <= width {
				current += " " + word
			} else {
				result.WriteString(current)
				result.WriteString("\n")
				current = word
			}
		}
		if current != "" {
			result.WriteString(current)
			result.WriteString("\n")
		}
	}

	return strings.TrimRight(result.String(), "\n")
}
