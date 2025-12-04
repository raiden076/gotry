package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/raiden076/gotry/internal/config"
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
