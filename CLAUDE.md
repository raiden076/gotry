# gotry (gt)

A universal alternative to [tobi/try](https://github.com/tobi/try) - an ephemeral workspace manager written in Go.

## What is this?

A CLI tool for managing experimental project directories. Create dated directories, fuzzy-search through them, auto-initialize git, and jump into any experiment instantly.

## Quick Start

```bash
gt                          # Open TUI selector
gt my-experiment            # Create ~/tries/2025-12-04-my-experiment/
gt https://github.com/u/r   # Clone into dated directory
eval "$(gt init zsh)"       # Add to .zshrc for shell integration
```

## Design Decisions

### Names
- **Binary:** `gotry`
- **Alias:** `gt` (shell function)
- **Config dir:** `~/.config/gotry/`
- **Default workspace:** `~/tries/`

### Tech Stack
- **Language:** Go (single binary distribution)
- **TUI:** Bubbletea + Bubbles + Lipgloss (Charm ecosystem)
- **Fuzzy:** sahilm/fuzzy + custom recency scoring
- **CLI:** Cobra + Viper

### Config File
Location: `~/.config/gotry/config.toml`

```toml
[workspace]
path = "~/tries"

[git]
auto_init = true
initial_commit = true
```

### Directory Structure
```
~/tries/
├── 2025-12-04-redis-experiment/
├── 2025-12-04-gotry/
├── 2025-12-03-api-testing/
└── 2025-11-28-tobi-try/
```

Collision handling: Append `-2`, `-3`, etc. if name exists on same day.

### Git Integration
- **New directory:** Auto `git init` + initial commit (configurable)
- **Clone:** Parse URL, create dated dir, clone into it
- **Initial commit message:**
  ```
  ✨ Let's try something new

  🤖 Created with gotry (https://github.com/YOUR_USER/gotry)
  ```

### Shell Integration
`gt init bash|zsh|fish|powershell` outputs shell-specific function.
User adds `eval "$(gt init zsh)"` to rc file (or `gotry init powershell | Invoke-Expression` in PowerShell `$PROFILE`).

---

## V1 Scope (Current)

### Commands
```
gt                     # Open TUI selector
gt <query>             # Open TUI with pre-filled search
gt <name>              # Create new directory if no match
gt <url>               # Clone repo into dated directory
gt init bash|zsh|fish|powershell  # Output shell integration
gt config              # Show current config
gt version             # Show version
```

### CLI Flags
```
gt --no-git <name>      # Skip git init
gt --no-commit <name>   # Init but no initial commit
gt --path <dir>         # Override workspace path
```

### TUI Keybindings
| Key | Action |
|-----|--------|
| `↑/↓` or `ctrl+p/n` | Navigate list |
| `Enter` | Select directory / Create new |
| `Backspace` | Delete search character |
| `ctrl+d` | Toggle delete mode |
| `ctrl+c` / `Esc` | Cancel / Exit |
| Typing | Fuzzy filter |

### TUI Display
```
┌─ gotry ─────────────────────────────────┐
│ Search: red█                            │
│                                         │
│ → 📁 2025-12-04-redis-experiment   2h   │
│   📁 2025-11-28-redis-testing      6d   │
│                                         │
│              [1-2 of 2]                 │
└─ enter: select · ctrl+d: delete · esc ──┘
```

### Features
- [x] Interactive TUI with fuzzy search
- [x] Create date-prefixed directories
- [x] Auto git init + initial commit (configurable)
- [x] Clone repos into dated directories
- [x] Recency-aware sorting
- [x] Batch delete with YES confirmation
- [x] Shell integration (bash/zsh/fish/powershell)
- [x] Config file support

---

## Future Features (Post-V1)

### Parity with try
- [ ] Git worktree support (`gt worktree <branch>`)
- [ ] Link/copy current directory (`gt . <name>`)
- [ ] `ctrl+a` / `ctrl+e` - Jump to start/end of search
- [ ] `ctrl+b` / `ctrl+f` - Move cursor left/right in search
- [ ] `ctrl+k` - Delete to end of line
- [ ] `ctrl+w` - Delete word backward

### Distribution & Installation
- [ ] Post-install shell setup instructions (detect shell, show specific commands)
- [ ] Homebrew formula (`brew install gotry`)
- [ ] Winget package (`winget install gotry`)
- [ ] Post-install hooks for package managers (brew, winget)

### Enhancements
- [ ] Configurable git init commit message
- [ ] Template scaffolding (`gt --template nextjs my-app`)
- [ ] Framework presets (Next.js, Hono, Remix, etc.)
- [ ] Cloud platform configs (Cloudflare, Netlify, Vercel)
- [ ] Custom templates from git repos
- [ ] Archive old experiments (compress + move)
- [ ] Tags/labels for directories
- [ ] Notes attached to experiments
- [ ] `gt list` - Non-interactive directory listing
- [ ] `gt clean` - Remove old/empty experiments
- [ ] `gt stats` - Show experiment statistics

### Template System (Future Project)
Vision: Universal template scaffolding tool.
```bash
gt --template nextjs-cloudflare my-app
gt --template hono-netlify my-api
gt --template rust-cli my-tool
```

---

## Project Structure

```
gotry/
├── cmd/
│   ├── root.go          # Main command, TUI launcher
│   ├── init.go          # Shell integration
│   ├── config.go        # Config management
│   └── version.go       # Version info
├── internal/
│   ├── config/          # Config loading/parsing
│   │   └── config.go
│   ├── tui/             # Bubbletea TUI
│   │   ├── model.go     # Main TUI model
│   │   ├── view.go      # Rendering
│   │   ├── update.go    # Event handling
│   │   └── styles.go    # Lipgloss styles
│   ├── workspace/       # Directory operations
│   │   └── workspace.go
│   └── git/             # Git operations
│       └── git.go
├── main.go
├── go.mod
├── go.sum
├── CLAUDE.md            # This file
├── README.md            # User-facing docs
└── LICENSE
```

---

## Development Notes

### Building
```bash
go build -o gotry .
```

### Testing Locally
```bash
./gotry
```

### Git Commits
**IMPORTANT:** Always use signed commits with `-S` flag:
```bash
git commit -S -m "feat: your message"
```

### Design Principles
1. **Simple yet production-ready** - No unwanted complexity
2. **Extensible** - Easy to add template system later
3. **Single binary** - No runtime dependencies
4. **Fast startup** - TUI should feel instant
