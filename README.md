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

# PowerShell ($PROFILE)
gotry init powershell | Invoke-Expression
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
