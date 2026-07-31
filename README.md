<div align="center">

# 🐙 `tuipr`

**A keyboard-driven Pull Request Lifecycle Manager for your terminal.**  
*Built with Go, Charm Bubbletea, and Catppuccin Mocha.*

[![Go Version](https://img.shields.io/github/go-mod/go-version/erickcortez/tuipr?style=flat-square&color=cba6f7)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square&color=89b4fa)](LICENSE)
[![Theme](https://img.shields.io/badge/theme-Catppuccin%20Mocha-fab387?style=flat-square)](https://github.com/catppuccin/catppuccin)

</div>

---

## ✨ Features

- **🎨 Catppuccin Mocha Theme** — Elegant, high-contrast terminal aesthetic
- **⌨️ LazyGit-style Navigation** — Numbered panels (1, 2, 3) + Vim keys (`j/k`, `Tab`)
- **📝 Vim Modes** — Normal/Insert mode with `i` and `Esc`
- **🔗 GitHub CLI (`gh`)** — List, create, and merge PRs directly
- **⚠️ Conflict Detection** — Visual indicators with instant `nvim` hotkey
- **🔀 Flexible Merging** — Merge Commit, Squash, or Rebase strategies

---

## 📦 Installation

```bash
git clone https://github.com/erickcortez/tuipr.git
cd tuipr
go build -o tuipr .
mv tuipr /usr/local/bin/
```

### Requirements

- Go 1.21+
- GitHub CLI ([`gh`](https://cli.github.com/)) authenticated
- Git
- Terminal with 256-color support

---

## 🚀 Usage

```bash
tuipr          # Open main dashboard
tuipr -c       # Open Create PR directly
tuipr -m       # Open Merge PR screen
tuipr -m 134   # Merge PR #134 directly
```

---

## ⌨️ Keybindings

### Dashboard

| Key | Action |
| :--- | :--- |
| `1` `2` `3` | Switch between PRs / Details / Conflicts panels |
| `j` / `k` | Navigate up / down |
| `c` | Open Create PR |
| `m` | Open Merge PR |
| `e` | Open nvim for conflicts |
| `r` | Refresh list |
| `q` | Quit |

### Create PR

| Key | Action |
| :--- | :--- |
| `1` `2` `3` | Switch Fields / Title / Description panels |
| `Tab` | Next panel |
| `i` | Enter Insert mode |
| `Esc` | Return to Normal |
| `Enter` | Select target branch |
| `Ctrl+s` | Submit PR |

### Merge PR

| Key | Action |
| :--- | :--- |
| `1` `2` `3` `4` | Switch Merge / Options / Checklist / Commit panels |
| `Tab` / `Shift+Tab` | Next / Previous panel |
| `j` / `k` | Navigate options |
| `Space` | Toggle / Select |
| `i` | Enter Insert mode for commit |
| `Ctrl+s` | Execute merge |

---

## ⚙️ Configuration

`tuipr` reads from `tuipr.toml` in the current directory or `~/.config/tuipr/config.toml`:

```toml
default_target_branch = "master"
transparent_panels = true
```

---

## 📋 Layout

```
┌──────────────────────────────────────────────────────────────┐
│  TUIPR | Branch: feature/auth | 5 PRs                        │
│  ─────────────────────────────────────────────────────────  │
│  ┌────────────┐  ┌──────────────┐  ┌─────────────────┐    │
│  │ (1) PRs    │  │ (2) Details   │  │ (3) Status      │    │
│  │ ────────── │  │ ───────────── │  │ ──────────────── │    │
│  │ -> #142 *  │  │ Fix auth      │  │ CI Checks [x]   │    │
│  │   #141 W   │  │ Author: @you  │  │ Reviews (1/1) [x]│    │
│  │   #140 *   │  │ Branch: ->main│  │ Conflicts [x]    │    │
│  └────────────┘  └──────────────┘  └─────────────────┘    │
│  ─────────────────────────────────────────────────────────  │
│  [1] PRs  [2] Details  [3] Status                            │
│  [j/k] Nav  [c] Create  [m] Merge  [q] Quit                  │
└──────────────────────────────────────────────────────────────┘
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

---

<div align="center">
Made with ❤️ for Terminal Power Users.
</div>
