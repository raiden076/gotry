package tui

import (
	"github.com/arkaprav0/gotry/internal/workspace"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeDelete
	ModeConfirm
)

type Model struct {
	// Config
	basePath string

	// State
	directories []workspace.Directory
	filtered    []workspace.Directory
	cursor      int
	mode        Mode
	marked      map[int]bool // indices marked for deletion
	confirmText string

	// Components
	searchInput textinput.Model

	// Output
	selected string
	quitting bool
}

func NewModel(basePath string, initialQuery string) Model {
	ti := textinput.New()
	ti.Placeholder = "Search or create..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	ti.SetValue(initialQuery)

	return Model{
		basePath:    basePath,
		searchInput: ti,
		marked:      make(map[int]bool),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.loadDirectories,
	)
}

func (m Model) loadDirectories() tea.Msg {
	dirs, err := workspace.List(m.basePath)
	if err != nil {
		return errMsg{err}
	}
	return dirsLoadedMsg{dirs}
}

// Messages
type dirsLoadedMsg struct {
	dirs []workspace.Directory
}

type errMsg struct {
	err error
}

// Fuzzy search implementation
type dirSearchable []workspace.Directory

func (d dirSearchable) String(i int) string {
	return d[i].Name
}

func (d dirSearchable) Len() int {
	return len(d)
}

func (m *Model) filterDirectories() {
	query := m.searchInput.Value()

	if query == "" {
		m.filtered = m.directories
		return
	}

	matches := fuzzy.FindFrom(query, dirSearchable(m.directories))

	m.filtered = make([]workspace.Directory, len(matches))
	for i, match := range matches {
		m.filtered[i] = m.directories[match.Index]
	}
}

func (m Model) Selected() string {
	return m.selected
}

func (m Model) Quitting() bool {
	return m.quitting
}
