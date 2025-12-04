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
