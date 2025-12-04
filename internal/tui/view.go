package tui

import (
	"fmt"
	"strings"

	"github.com/arkaprav0/gotry/internal/workspace"
)

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("gotry"))
	b.WriteString("\n")

	// Search input
	b.WriteString(searchPromptStyle.Render("Search: "))
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	// Directory list
	if len(m.filtered) == 0 {
		if m.searchInput.Value() != "" {
			b.WriteString(dimStyle.Render("  No matches. Press enter to create: "))
			b.WriteString(normalStyle.Render(m.searchInput.Value()))
		} else {
			b.WriteString(dimStyle.Render("  No experiments yet. Type a name to create one."))
		}
		b.WriteString("\n")
	} else {
		for i, dir := range m.filtered {
			b.WriteString(m.renderDirectoryItem(i, dir))
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	return b.String()
}

func (m Model) renderDirectoryItem(index int, dir workspace.Directory) string {
	var b strings.Builder

	// Selection indicator
	if index == m.cursor {
		b.WriteString(selectedStyle.Render(" → "))
	} else {
		b.WriteString("   ")
	}

	// Icon
	if m.mode == ModeDelete && m.marked[index] {
		b.WriteString(markedStyle.Render("🗑️  "))
	} else {
		b.WriteString("📁 ")
	}

	// Directory name
	name := dir.Name
	if m.mode == ModeDelete && m.marked[index] {
		name = deleteStyle.Render(name)
	} else if index == m.cursor {
		if dir.DatePart != "" {
			name = dimStyle.Render(dir.DatePart+"-") + selectedStyle.Render(dir.NamePart)
		} else {
			name = selectedStyle.Render(name)
		}
	} else {
		if dir.DatePart != "" {
			name = dimStyle.Render(dir.DatePart+"-") + normalStyle.Render(dir.NamePart)
		} else {
			name = normalStyle.Render(name)
		}
	}
	b.WriteString(name)

	// Relative time
	relTime := workspace.RelativeTime(dir.ModTime)
	padding := 40 - len(dir.Name)
	if padding < 2 {
		padding = 2
	}
	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(dimStyle.Render(relTime))

	return b.String()
}

func (m Model) renderFooter() string {
	switch m.mode {
	case ModeConfirm:
		return fmt.Sprintf(
			"%s Type %s to confirm deletion (%d items): %s",
			markedStyle.Render("⚠️"),
			markedStyle.Render("YES"),
			len(m.marked),
			m.confirmText,
		)

	case ModeDelete:
		count := len(m.marked)
		if count > 0 {
			return fmt.Sprintf(
				"%s · %s · %s",
				helpKeyStyle.Render("space")+" "+helpDescStyle.Render("toggle mark"),
				helpKeyStyle.Render("enter")+" "+helpDescStyle.Render(fmt.Sprintf("delete %d", count)),
				helpKeyStyle.Render("esc")+" "+helpDescStyle.Render("cancel"),
			)
		}
		return fmt.Sprintf(
			"%s · %s",
			helpKeyStyle.Render("space")+" "+helpDescStyle.Render("mark for deletion"),
			helpKeyStyle.Render("esc")+" "+helpDescStyle.Render("cancel"),
		)

	default:
		return fmt.Sprintf(
			"%s · %s · %s",
			helpKeyStyle.Render("enter")+" "+helpDescStyle.Render("select"),
			helpKeyStyle.Render("ctrl+d")+" "+helpDescStyle.Render("delete"),
			helpKeyStyle.Render("esc")+" "+helpDescStyle.Render("quit"),
		)
	}
}
