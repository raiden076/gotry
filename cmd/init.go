package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [bash|zsh|fish|powershell]",
	Short: "Output shell integration script",
	Long:  `Output shell function for your shell. Add to your .bashrc, .zshrc, config.fish, or $PROFILE`,
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
	case "powershell", "pwsh":
		fmt.Print(powershellInit)
	default:
		return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish, powershell)", shell)
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

const powershellInit = `# gotry shell integration
# Add this to your PowerShell profile ($PROFILE):
#   gotry init powershell | Invoke-Expression

function gt {
    $result = gotry @args
    $exitCode = $LASTEXITCODE

    if ($exitCode -eq 0 -and $result -and (Test-Path -Path $result -PathType Container)) {
        Set-Location $result
    } elseif ($result) {
        Write-Output $result
    }

    return $exitCode
}
`
