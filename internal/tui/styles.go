package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("212") // Pink
	secondaryColor = lipgloss.Color("241") // Gray
	accentColor    = lipgloss.Color("229") // Yellow
	dangerColor    = lipgloss.Color("196") // Red

	// Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Search input
	searchPromptStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)

	searchInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	// List items
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	dimStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	matchStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	// Delete mode
	deleteStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Strikethrough(true)

	markedStyle = lipgloss.NewStyle().
			Foreground(dangerColor)

	// Footer
	footerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			MarginTop(1)

	// Help
	helpKeyStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)
)
