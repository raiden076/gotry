package tui

import (
	"strings"

	"github.com/raiden076/gotry/internal/workspace"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case dirsLoadedMsg:
		m.directories = msg.dirs
		m.filterDirectories()
		return m, nil

	case errMsg:
		// Handle error - for now just quit
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	// Update text input
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.filterDirectories()

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}

	return m, cmd
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeConfirm:
		return m.handleConfirmMode(msg)
	case ModeDelete:
		return m.handleDeleteMode(msg)
	default:
		return m.handleNormalMode(msg)
	}
}

func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.quitting = true
		return m, tea.Quit

	case "enter":
		return m.handleEnter()

	case "up", "ctrl+p":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case "down", "ctrl+n":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
		return m, nil

	case "ctrl+d":
		if len(m.filtered) > 0 {
			m.mode = ModeDelete
		}
		return m, nil
	}

	// Pass to text input
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.filterDirectories()

	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}

	return m, cmd
}

func (m Model) handleDeleteMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.mode = ModeNormal
		m.marked = make(map[int]bool)
		return m, nil

	case "enter":
		if len(m.marked) > 0 {
			m.mode = ModeConfirm
			m.confirmText = ""
		}
		return m, nil

	case "up", "ctrl+p":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case "down", "ctrl+n":
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
		return m, nil

	case "ctrl+d", " ":
		// Toggle mark on current item
		if m.cursor < len(m.filtered) {
			if m.marked[m.cursor] {
				delete(m.marked, m.cursor)
			} else {
				m.marked[m.cursor] = true
			}
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleConfirmMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.mode = ModeDelete
		m.confirmText = ""
		return m, nil

	case "backspace":
		if len(m.confirmText) > 0 {
			m.confirmText = m.confirmText[:len(m.confirmText)-1]
		}
		return m, nil

	case "enter":
		if m.confirmText == "YES" {
			return m.executeDelete()
		}
		return m, nil

	default:
		if len(msg.String()) == 1 {
			m.confirmText += strings.ToUpper(msg.String())
		}
		return m, nil
	}
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	query := m.searchInput.Value()

	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		// Select existing directory
		m.selected = m.filtered[m.cursor].Path
		return m, tea.Quit
	}

	if query != "" {
		// Create new directory
		m.selected = "CREATE:" + query
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) executeDelete() (tea.Model, tea.Cmd) {
	var paths []string
	for idx := range m.marked {
		if idx < len(m.filtered) {
			paths = append(paths, m.filtered[idx].Path)
		}
	}

	if err := workspace.Delete(paths); err != nil {
		// Handle error
		m.mode = ModeNormal
		m.marked = make(map[int]bool)
		return m, nil
	}

	// Reload directories
	m.mode = ModeNormal
	m.marked = make(map[int]bool)
	m.confirmText = ""

	return m, m.loadDirectories
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
