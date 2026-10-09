# Vitualizz Dotfiles

[![CI](https://github.com/vitualizz/dotfiles/actions/workflows/ci.yml/badge.svg)](https://github.com/vitualizz/dotfiles/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.24-blue?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

Interactive installer for my terminal environment on Linux and macOS: Kitty, zsh, Starship, Neovim and the CLI tools I use every day, plus a starting config for each.

> **Platform:** Linux (pacman, apt, apk, dnf) and macOS (Homebrew). Windows and BSD are **not supported**.

## Quick Start

Linux and macOS (macOS requires [Homebrew](https://brew.sh)):

```bash
curl -fsSL https://raw.githubusercontent.com/vitualizz/dotfiles/master/install.sh | bash
```

Do **not** prefix it with `sudo`: your config would be installed for root instead of your user. When a tool needs administrator rights (apt, pacman, Homebrew casks), the installer asks for your password once and keeps it valid until it finishes.

The script downloads a pre-compiled binary from the latest GitHub release, runs it, and deletes it.

## How config works: seed and let go

The installer **installs**; it does not **own** your environment. Each config is copied once to the tool's own location, and from then on it is yours. Re-running the installer never overwrites it.

| Config | Where it goes | If it already exists |
|--------|---------------|----------------------|
| Starship prompt | `~/.config/starship.toml` | kept as is |
| Kitty + Tokyo Night theme | `~/.config/kitty/kitty.conf`, `tokyonight_night.conf` | kept as is (per file) |
| zsh integration | `~/.zshrc` | a commented `vitualizz dotfiles` block is appended once |
| Git | `~/.gitconfig` (delta pager, push/pull/rebase defaults) | kept as is; identity and signing go in `~/.gitconfig.local`, included last |

The `~/.zshrc` block wires up what the installer adds: Homebrew, completion (case-insensitive, menu), history, Starship, zoxide, direnv, fzf (fd + bat/eza previews), atuin, the zsh plugins, Oh My Zsh's git aliases, and `ls`/`cat` through eza/bat. It sources `~/.zshrc.local` last, for secrets and machine-specific settings that never go to git. Each section only acts if that is not already configured above it, so appending it to an existing `~/.zshrc` doesn't override your setup. Edit, move or delete it freely. A backup (`~/.zshrc.bak-<timestamp>`) is made before appending.

Uninstalling tools never touches these files.

## Theme

- **Neovim**: Everblush, with the base16 palette NvChad uses, applied by `mini.base16` (`nvim/lua/vitualizz/plugins/colorscheme.lua`). Transparent background (Kitty's shows through). Comments and line numbers are lightened to 4.5:1 and 3:1 contrast; Everblush's own grays are below 3:1.
- **Kitty**: Tokyo Night, from the official `tokyonight_night.conf`. **Starship** uses the terminal's ANSI colors, so it follows Kitty.

## What's installed

| Category | Tools |
|----------|-------|
| **Terminal** | Kitty |
| **Shell** | zsh, Starship, zsh-autosuggestions, zsh-syntax-highlighting, atuin, zoxide, fzf, direnv |
| **Editor** | Neovim, ripgrep, fd, C compiler (for nvim-treesitter) |
| **Git** | git, lazygit, delta, gh, gnupg |
| **Files** | yazi, eza, bat, sd |
| **Data** | jq, yq |
| **System** | btop, duf, dust, hyperfine |
| **Docs / info** | tealdeer, glow, fastfetch, onefetch |
| **Runtimes** | mise |
| **Fonts** | Hack, JetBrains Mono, Fira Code (Nerd Fonts) |

## Supported Platforms

- **Arch Linux** (pacman)
- **Debian / Ubuntu** (apt)
- **Alpine** (apk)
- **Fedora** (dnf)
- **macOS** on Apple Silicon or Intel (Homebrew)

### macOS notes

- **Prerequisites:** [Homebrew](https://brew.sh) (it also installs the Xcode Command Line Tools). The installer stops early with instructions if it is missing.
- **git** is installed from Homebrew, replacing Apple's older bundled git.
- **Fonts** are installed as Homebrew casks into `~/Library/Fonts`.

## How tools are resolved

Each tool in `cmd/dotfiles/config/tools.yaml` declares install commands per platform. The installer detects the platform and picks the first match:

- **Linux:** `exact distro → all → other distros → fallback`
- **macOS:** `macos → brew → all → fallback`. Linux commands (apt, pacman, sudo…) are never used.

```yaml
- name: ripgrep
  install:
    debian: sudo apt-get install -y ripgrep
    arch: sudo pacman -S --noconfirm ripgrep
    brew: brew install ripgrep
  check: rg --version
```

Homebrew casks and other macOS-only commands go under `macos:`. `internal/config/tools_yaml_test.go` fails if a tool resolves to a Linux-only command on macOS or to a macOS command on Linux, and checks that seeded configs are never overwritten.

## Architecture

Hexagonal (Ports & Adapters): UI and infrastructure only talk through interfaces.

```
cmd/dotfiles/
  config/           embedded into the binary (tools.yaml, kitty/, zsh/, starship/)
internal/
  domain/
    entities/       Tool, Distro, Category
    interfaces/     InstallerPort, SudoPort, ToolRepository
  usecases/         install, uninstall, dependency order
  config/           YAML-backed ToolRepository
  infrastructure/
    executor/       runs shell commands
    installers/     resolves and runs each tool's commands
    logger/         ~/.local/state/vitualizz-dotfiles/install.log
  ui/               Bubble Tea TUI
i18n/               en/es translations
```

## Development

Requires Go 1.24+ (macOS: `brew install go`).

```bash
go run ./cmd/dotfiles/          # interactive TUI
go run ./cmd/dotfiles/ --ci     # headless: no TUI, exit code 1 if a tool fails
go test -race -cover ./...
```

`DOTFILES_CONFIG=/path/to/tools.yaml` runs with an external config instead of the embedded one.

### Isolated environments

```bash
docker compose run app          # Linux container, --ci mode (skips kitty)
docker compose run test         # tests in a clean container
vagrant up ubuntu               # or: vagrant up arch
```

## Releases

Pushing a tag builds binaries for linux/darwin × amd64/arm64 with GoReleaser:

```bash
git tag v1.0.0 && git push origin v1.0.0
```

## License

MIT

---

Con amor @vitualizz
