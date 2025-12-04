# gotry Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go-based ephemeral workspace manager with TUI, fuzzy search, git auto-init, and shell integration.

**Architecture:** CLI built with Cobra handles commands, Bubbletea powers the interactive TUI, internal packages separate concerns (config, workspace, git, tui). Config via Viper with TOML file.

**Tech Stack:** Go, Cobra, Viper, Bubbletea, Bubbles, Lipgloss, sahilm/fuzzy

---

## Phase 1: Project Setup

### Task 1: Initialize Go Module

**Files:**
- Create: `go.mod`
- Create: `main.go`

**Step 1: Initialize go module**

Run:
```bash
cd /home/arkaprav0/clone-try
go mod init github.com/arkaprav0/gotry
```

**Step 2: Create minimal main.go**

```go
package main

import "fmt"

func main() {
	fmt.Println("gotry")
}
```

**Step 3: Verify it builds and runs**

Run: `go build -o gotry . && ./gotry`
Expected: `gotry`

**Step 4: Commit**

```bash
git init
git add go.mod main.go CLAUDE.md
git commit -m "feat: initialize gotry project"
```

---

### Task 2: Add Dependencies

**Files:**
- Modify: `go.mod`

**Step 1: Add all required dependencies**

Run:
```bash
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/sahilm/fuzzy@latest
```

**Step 2: Tidy modules**

Run: `go mod tidy`

**Step 3: Verify go.sum exists**

Run: `ls go.sum`
Expected: `go.sum`

**Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: add dependencies (cobra, viper, bubbletea, lipgloss, fuzzy)"
```

---

## Phase 2: Configuration System

### Task 3: Create Config Package

**Files:**
- Create: `internal/config/config.go`

**Step 1: Create config structure and loader**

```go
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Workspace WorkspaceConfig `mapstructure:"workspace"`
	Git       GitConfig       `mapstructure:"git"`
}

type WorkspaceConfig struct {
	Path string `mapstructure:"path"`
}

type GitConfig struct {
	AutoInit      bool `mapstructure:"auto_init"`
	InitialCommit bool `mapstructure:"initial_commit"`
}

func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Workspace: WorkspaceConfig{
			Path: filepath.Join(homeDir, "tries"),
		},
		Git: GitConfig{
			AutoInit:      true,
			InitialCommit: true,
		},
	}
}

func Load() (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil // Return defaults if home dir unavailable
	}

	configDir := filepath.Join(homeDir, ".config", "gotry")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil // Config file not found, use defaults
		}
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// Expand ~ in workspace path
	if cfg.Workspace.Path[:2] == "~/" {
		cfg.Workspace.Path = filepath.Join(homeDir, cfg.Workspace.Path[2:])
	}

	return cfg, nil
}

func (c *Config) EnsureWorkspaceExists() error {
	return os.MkdirAll(c.Workspace.Path, 0755)
}
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/config/config.go
git commit -m "feat: add config package with defaults and TOML loading"
```

---

## Phase 3: Workspace Operations

### Task 4: Create Workspace Package

**Files:**
- Create: `internal/workspace/workspace.go`

**Step 1: Create workspace operations**

```go
package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Directory struct {
	Name      string
	Path      string
	ModTime   time.Time
	DatePart  string
	NamePart  string
}

func List(basePath string) ([]Directory, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Directory{}, nil
		}
		return nil, err
	}

	var dirs []Directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		datePart, namePart := parseDirectoryName(name)

		dirs = append(dirs, Directory{
			Name:     name,
			Path:     filepath.Join(basePath, name),
			ModTime:  info.ModTime(),
			DatePart: datePart,
			NamePart: namePart,
		})
	}

	// Sort by modification time, most recent first
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].ModTime.After(dirs[j].ModTime)
	})

	return dirs, nil
}

func parseDirectoryName(name string) (datePart, namePart string) {
	// Expected format: YYYY-MM-DD-name
	if len(name) >= 11 && name[4] == '-' && name[7] == '-' && name[10] == '-' {
		return name[:10], name[11:]
	}
	return "", name
}

func Create(basePath, name string) (string, error) {
	today := time.Now().Format("2006-01-02")
	dirName := fmt.Sprintf("%s-%s", today, sanitizeName(name))
	fullPath := filepath.Join(basePath, dirName)

	// Handle collisions
	finalPath := fullPath
	counter := 2
	for {
		if _, err := os.Stat(finalPath); os.IsNotExist(err) {
			break
		}
		finalPath = fmt.Sprintf("%s-%d", fullPath, counter)
		counter++
	}

	if err := os.MkdirAll(finalPath, 0755); err != nil {
		return "", err
	}

	return finalPath, nil
}

func sanitizeName(name string) string {
	// Replace spaces with hyphens, lowercase
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	return name
}

func Delete(paths []string) error {
	for _, path := range paths {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func RelativeTime(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "now"
	case duration < time.Hour:
		mins := int(duration.Minutes())
		return fmt.Sprintf("%dm", mins)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	default:
		weeks := int(duration.Hours() / 24 / 7)
		return fmt.Sprintf("%dw", weeks)
	}
}
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/workspace/workspace.go
git commit -m "feat: add workspace package for directory operations"
```

---

## Phase 4: Git Operations

### Task 5: Create Git Package

**Files:**
- Create: `internal/git/git.go`

**Step 1: Create git operations**

```go
package git

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const commitMessage = `✨ Let's try something new

🤖 Created with gotry (https://github.com/arkaprav0/gotry)`

func Init(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitialCommit(path string) error {
	// Create .gitkeep to have something to commit
	gitkeep := filepath.Join(path, ".gitkeep")
	if err := os.WriteFile(gitkeep, []byte{}, 0644); err != nil {
		return err
	}

	// git add .
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = path
	if err := addCmd.Run(); err != nil {
		return err
	}

	// git commit
	commitCmd := exec.Command("git", "commit", "-m", commitMessage)
	commitCmd.Dir = path
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	return commitCmd.Run()
}

func Clone(repoURL, destPath string) error {
	cmd := exec.Command("git", "clone", repoURL, destPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type RepoInfo struct {
	Host string
	User string
	Repo string
}

func ParseGitURL(rawURL string) (*RepoInfo, error) {
	// Handle SSH format: git@github.com:user/repo.git
	sshRegex := regexp.MustCompile(`^git@([^:]+):([^/]+)/(.+?)(?:\.git)?$`)
	if matches := sshRegex.FindStringSubmatch(rawURL); matches != nil {
		return &RepoInfo{
			Host: matches[1],
			User: matches[2],
			Repo: matches[3],
		}, nil
	}

	// Handle HTTPS format
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid repository URL: %s", rawURL)
	}

	repo := pathParts[1]
	repo = strings.TrimSuffix(repo, ".git")

	return &RepoInfo{
		Host: parsed.Host,
		User: pathParts[0],
		Repo: repo,
	}, nil
}

func (r *RepoInfo) DirectoryName() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s-%s-%s", today, r.User, r.Repo)
}

func IsGitURL(s string) bool {
	return strings.HasPrefix(s, "git@") ||
		strings.HasPrefix(s, "https://github.com") ||
		strings.HasPrefix(s, "https://gitlab.com") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") && strings.Contains(s, ".git")
}
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/git/git.go
git commit -m "feat: add git package for init, commit, clone operations"
```

---

## Phase 5: TUI Implementation

### Task 6: Create TUI Styles

**Files:**
- Create: `internal/tui/styles.go`

**Step 1: Define lipgloss styles**

```go
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
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/tui/styles.go
git commit -m "feat: add TUI styles with lipgloss"
```

---

### Task 7: Create TUI Model

**Files:**
- Create: `internal/tui/model.go`

**Step 1: Define TUI model and state**

```go
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
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/tui/model.go
git commit -m "feat: add TUI model with fuzzy search"
```

---

### Task 8: Create TUI Update Handler

**Files:**
- Create: `internal/tui/update.go`

**Step 1: Implement update logic**

```go
package tui

import (
	"strings"

	"github.com/arkaprav0/gotry/internal/workspace"
	"github.com/charmbracelet/bubbles/textinput"
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
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/tui/update.go
git commit -m "feat: add TUI update handler with all modes"
```

---

### Task 9: Create TUI View

**Files:**
- Create: `internal/tui/view.go`

**Step 1: Implement view rendering**

```go
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
```

**Step 2: Verify it compiles**

Run: `go build ./...`
Expected: No errors

**Step 3: Commit**

```bash
git add internal/tui/view.go
git commit -m "feat: add TUI view rendering"
```

---

## Phase 6: CLI Commands

### Task 10: Create Root Command

**Files:**
- Create: `cmd/root.go`
- Modify: `main.go`

**Step 1: Create root command with TUI launcher**

```go
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/arkaprav0/gotry/internal/config"
	"github.com/arkaprav0/gotry/internal/git"
	"github.com/arkaprav0/gotry/internal/tui"
	"github.com/arkaprav0/gotry/internal/workspace"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	flagNoGit    bool
	flagNoCommit bool
	flagPath     string
)

var rootCmd = &cobra.Command{
	Use:   "gotry [name|url|query]",
	Short: "Ephemeral workspace manager",
	Long:  `gotry (gt) - A universal alternative to try. Manage experimental project directories with ease.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&flagNoGit, "no-git", false, "Skip git initialization")
	rootCmd.Flags().BoolVar(&flagNoCommit, "no-commit", false, "Skip initial commit")
	rootCmd.Flags().StringVar(&flagPath, "path", "", "Override workspace path")
}

func runRoot(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if flagPath != "" {
		cfg.Workspace.Path = flagPath
	}

	if err := cfg.EnsureWorkspaceExists(); err != nil {
		return err
	}

	// Check if argument is a git URL
	if len(args) == 1 && git.IsGitURL(args[0]) {
		return handleClone(cfg, args[0])
	}

	// Launch TUI
	initialQuery := ""
	if len(args) == 1 {
		initialQuery = args[0]
	}

	model := tui.NewModel(cfg.Workspace.Path, initialQuery)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	m := finalModel.(tui.Model)

	if m.Quitting() && m.Selected() == "" {
		return nil // User cancelled
	}

	selected := m.Selected()

	// Handle create new directory
	if strings.HasPrefix(selected, "CREATE:") {
		name := strings.TrimPrefix(selected, "CREATE:")
		return handleCreate(cfg, name)
	}

	// Output selected path for shell integration
	if selected != "" {
		fmt.Println(selected)
	}

	return nil
}

func handleCreate(cfg *config.Config, name string) error {
	path, err := workspace.Create(cfg.Workspace.Path, name)
	if err != nil {
		return err
	}

	// Git init
	if cfg.Git.AutoInit && !flagNoGit {
		if err := git.Init(path); err != nil {
			return err
		}

		// Initial commit
		if cfg.Git.InitialCommit && !flagNoCommit {
			if err := git.InitialCommit(path); err != nil {
				return err
			}
		}
	}

	fmt.Println(path)
	return nil
}

func handleClone(cfg *config.Config, url string) error {
	info, err := git.ParseGitURL(url)
	if err != nil {
		return err
	}

	dirName := info.DirectoryName()
	destPath := cfg.Workspace.Path + "/" + dirName

	// Handle collision
	counter := 2
	basePath := destPath
	for {
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			break
		}
		destPath = fmt.Sprintf("%s-%d", basePath, counter)
		counter++
	}

	if err := git.Clone(url, destPath); err != nil {
		return err
	}

	fmt.Println(destPath)
	return nil
}
```

**Step 2: Update main.go**

```go
package main

import "github.com/arkaprav0/gotry/cmd"

func main() {
	cmd.Execute()
}
```

**Step 3: Verify it compiles**

Run: `go build -o gotry .`
Expected: Binary created

**Step 4: Commit**

```bash
git add cmd/root.go main.go
git commit -m "feat: add root command with TUI launcher"
```

---

### Task 11: Create Init Command (Shell Integration)

**Files:**
- Create: `cmd/init.go`

**Step 1: Create shell integration command**

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [bash|zsh|fish]",
	Short: "Output shell integration script",
	Long:  `Output shell function for your shell. Add to your .bashrc, .zshrc, or config.fish`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	shell := args[0]

	switch shell {
	case "bash", "zsh":
		fmt.Print(bashZshInit)
	case "fish":
		fmt.Print(fishInit)
	default:
		return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish)", shell)
	}

	return nil
}

const bashZshInit = `# gotry shell integration
# Add this to your .bashrc or .zshrc:
#   eval "$(gotry init bash)"  # or zsh

gt() {
    local result
    result=$(gotry "$@")
    local exit_code=$?

    if [ $exit_code -eq 0 ] && [ -n "$result" ] && [ -d "$result" ]; then
        cd "$result"
    elif [ -n "$result" ]; then
        echo "$result"
    fi

    return $exit_code
}
`

const fishInit = `# gotry shell integration
# Add this to your config.fish:
#   gotry init fish | source

function gt
    set -l result (gotry $argv)
    set -l exit_code $status

    if test $exit_code -eq 0; and test -n "$result"; and test -d "$result"
        cd "$result"
    else if test -n "$result"
        echo "$result"
    end

    return $exit_code
end
`
```

**Step 2: Verify it compiles**

Run: `go build -o gotry . && ./gotry init zsh`
Expected: Shell function output

**Step 3: Commit**

```bash
git add cmd/init.go
git commit -m "feat: add init command for shell integration"
```

---

### Task 12: Create Config Command

**Files:**
- Create: `cmd/config.go`

**Step 1: Create config display command**

```go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arkaprav0/gotry/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "gotry", "config.toml")

	fmt.Println("Configuration")
	fmt.Println("─────────────")
	fmt.Printf("Config file:    %s\n", configPath)
	fmt.Printf("Workspace:      %s\n", cfg.Workspace.Path)
	fmt.Printf("Auto git init:  %t\n", cfg.Git.AutoInit)
	fmt.Printf("Initial commit: %t\n", cfg.Git.InitialCommit)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("\n(Using defaults - no config file found)")
		fmt.Println("\nCreate config file with:")
		fmt.Printf("  mkdir -p %s\n", filepath.Dir(configPath))
		fmt.Printf("  cat > %s << 'EOF'\n", configPath)
		fmt.Println(`[workspace]
path = "~/tries"

[git]
auto_init = true
initial_commit = true
EOF`)
	}

	return nil
}
```

**Step 2: Verify it compiles**

Run: `go build -o gotry . && ./gotry config`
Expected: Config output

**Step 3: Commit**

```bash
git add cmd/config.go
git commit -m "feat: add config command"
```

---

### Task 13: Create Version Command

**Files:**
- Create: `cmd/version.go`

**Step 1: Create version command**

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gotry version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

**Step 2: Verify it compiles**

Run: `go build -o gotry . && ./gotry version`
Expected: `gotry version 0.1.0`

**Step 3: Commit**

```bash
git add cmd/version.go
git commit -m "feat: add version command"
```

---

## Phase 7: Final Polish

### Task 14: Test Full Flow

**Step 1: Build final binary**

Run: `go build -o gotry .`

**Step 2: Test create directory**

Run: `./gotry test-experiment`
Expected: TUI opens, can create directory

**Step 3: Test TUI navigation**

Run: `./gotry`
Expected: TUI shows created directory, can navigate

**Step 4: Test shell integration**

Run: `eval "$(./gotry init zsh)" && gt`
Expected: Shell function works, cd's into selected directory

**Step 5: Test clone**

Run: `./gotry https://github.com/charmbracelet/bubbletea`
Expected: Clones repo into dated directory

---

### Task 15: Create README

**Files:**
- Create: `README.md`

**Step 1: Create user-facing documentation**

```markdown
# gotry (gt)

A universal alternative to [try](https://github.com/tobi/try) - an ephemeral workspace manager written in Go.

## Installation

### From source

```bash
go install github.com/arkaprav0/gotry@latest
```

### Shell integration

Add to your shell rc file:

```bash
# bash (~/.bashrc)
eval "$(gotry init bash)"

# zsh (~/.zshrc)
eval "$(gotry init zsh)"

# fish (~/.config/fish/config.fish)
gotry init fish | source
```

## Usage

```bash
gt                          # Open interactive selector
gt my-experiment            # Create ~/tries/2025-12-04-my-experiment/
gt redis                    # Fuzzy search, select or create
gt https://github.com/u/r   # Clone repo into dated directory
```

## Features

- **Interactive TUI** with fuzzy search
- **Date-prefixed directories** for chronological organization
- **Auto git init** with configurable initial commit
- **Clone repos** directly into your tries directory
- **Recency sorting** - recent experiments appear first
- **Batch delete** with safety confirmation

## Configuration

Create `~/.config/gotry/config.toml`:

```toml
[workspace]
path = "~/tries"

[git]
auto_init = true
initial_commit = true
```

## Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate |
| `Enter` | Select / Create |
| `Ctrl+D` | Delete mode |
| `Esc` | Cancel / Quit |

## License

MIT
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README"
```

---

### Task 16: Final Commit and Tag

**Step 1: Ensure everything builds**

Run: `go build -o gotry . && go mod tidy`

**Step 2: Final commit**

```bash
git add -A
git commit -m "chore: v0.1.0 release prep"
git tag v0.1.0
```

---

## Summary

**16 tasks total across 7 phases:**

1. **Project Setup** (Tasks 1-2): Initialize Go module, add dependencies
2. **Configuration** (Task 3): Config package with TOML support
3. **Workspace** (Task 4): Directory listing, creation, deletion
4. **Git** (Task 5): Init, commit, clone, URL parsing
5. **TUI** (Tasks 6-9): Styles, model, update, view
6. **CLI** (Tasks 10-13): Root, init, config, version commands
7. **Polish** (Tasks 14-16): Testing, README, release

**Estimated commits:** 16 focused commits
