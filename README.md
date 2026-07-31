<div align="center">

# 🐙 `tuipr`

**A keyboard-driven Pull Request Lifecycle Manager for your terminal.**  
*Built with Go, Charm Bubbletea, and Catppuccin Mocha.*

[![Go Version](https://img.shields.io/github/gomod/go-version/erickcortez/tuipr?style=flat-square&color=cba6f7)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square&color=89b4fa)](LICENSE)
[![Theme](https://img.shields.io/badge/theme-Catppuccin%20Mocha-fab387?style=flat-square)](https://github.com/catppuccin/catppuccin)

[Key Features](#-features) •
[Workflow Integration](#-workflow-integration) •
[Installation](#-installation) •
[Keybindings](#-keybindings) •
[Configuration](#-configuration)

</div>

---

## ✨ Features

- **⚡ Neovim-Inspired Navigation:** Panes, buffers, and Vim keybindings (`h/j/k/l`, `Tab`, Normal/Insert modes).
- **🎨 Catppuccin Mocha Theme:** Clean, high-contrast, distraction-free terminal aesthetic.
- **🚀 Zero-Browser PR Lifecycle:** Create, review status, inspect conflicts, and merge PRs directly from your terminal.
- **⚠️ Smart Conflict Detection:** Visual banners for merge conflicts with instant `nvim` hotkey integration.
- **🔀 Flexible Merging:** Support for Merge Commits, Squash, and Rebase with automatic branch cleanup.
- **🛠️ Power-User Integrations:** Works alongside `lazygit`, `tuicr`, and `gh-dash`.

---

## 🔄 Workflow Integration

`tuipr` bridges the gap between local Git operations and GitHub PR management:

```text
 1. LazyGit        ──> Commit changes and git push branch.
 2. tuipr (Create) ──> Press 'c' to open the Create PR buffer, set target branch, edit body, and submit.
 3. tuicr / Review ──> Review code diffs and approvals.
 4. tuipr (Merge)  ──> Inspect conflict status. Press 'e' to resolve in Nvim OR press 'm' to merge.
```

---

## ⌨️ Keybindings

### Global & Dashboard (`Normal Mode`)

| Key | Action |
| :--- | :--- |
| `j` / `k` | Move cursor Up / Down in lists |
| `h` / `l` or `Tab` | Switch active panel |
| `c` | Open **Create PR** buffer |
| `m` | Open **Merge PR** screen (blocked if conflicts exist) |
| `e` | Launch `Neovim` to resolve merge conflicts in local working directory |
| `r` | Refresh PR list & GitHub status |
| `q` / `Ctrl+c` | Quit `tuipr` |

### Create PR Buffer (`Normal` & `Insert` Modes)

| Key | Mode | Action |
| :--- | :--- | :--- |
| `i` | Normal | Enter **Insert Mode** on focused field |
| `<Esc>` | Insert | Return to **Normal Mode** |
| `j` / `k` / `Tab` | Normal | Switch between Target Branch, Title, and Description Buffer |
| `Ctrl+s` | Normal/Insert | Submit Pull Request & Push |
| `Ctrl+c` | Normal/Insert | Cancel and exit buffer |

---

## ⚙️ Configuration

`tuipr` reads settings from `~/.config/tuipr/config.toml`:

```toml
[defaults]
target_branch = "main"
editor = "nvim"

[merge]
default_strategy = "merge" # Options: merge, squash, rebase
delete_remote_branch = true
delete_local_branch = true
default_commit_message = "Merged via tuipr 🚀"

[theme]
palette = "catppuccin-mocha"
```

---

## 📦 Installation

### Prerequisites

- [Go](https://golang.org/doc/install) 1.21+
- [GitHub CLI (`gh`)](https://cli.github.com/) authenticated (`gh auth login`)
- [Git](https://git-scm.com/)

### Build from Source

```bash
git clone [https://github.com/erickcortez/tuipr.git](https://github.com/erickcortez/tuipr.git)
cd tuipr
go build -o tuipr main.go
mv tuipr /usr/local/bin/
```

### `gh-dash` Integration

To invoke `tuipr` directly from `gh-dash`, add the following keybinding to your `~/.config/gh-dash/config.yml`:

```yaml
keybindings:
  prs:
    - key: "m"
      command: "tuipr"
```

---

<div align="center">
Made with ❤️ for Terminal Power Users.
</div>
