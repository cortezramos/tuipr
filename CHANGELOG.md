# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [0.2.0] - 2024-08-07

### Added
- **Homebrew support** - Install tuipr via Homebrew
  - `brew tap erickcortez/tuipr`
  - `brew install tuipr`
- Cross-platform release builds
  - macOS Intel (amd64)
  - macOS Apple Silicon (arm64)
  - Linux Intel (amd64)
- GitHub Actions workflow for automated releases
- Build script for local testing (`scripts/build-release.sh`)

### Fixed
- Corrected module path in go.mod

---

## [0.1.0] - 2024-07-31

### Added
- Initial release
- Dashboard with PR list, details, and conflict status
- Create PR buffer with Vim-style editing
- Merge PR screen with strategy selection (Merge Commit, Squash, Rebase)
- Catppuccin Mocha theme
- LazyGit-style navigation (numbered panels)
- Neovim integration for conflict resolution
- GitHub CLI (`gh`) integration
- Unit tests with 100% coverage on core functionality
- CI/CD pipeline with GitHub Actions

### Features
- `tuipr` - Open dashboard
- `tuipr -c` - Open Create PR directly
- `tuipr -m` - Open Merge PR screen
- `tuipr -m <num>` - Merge specific PR

---

[0.2.0]: https://github.com/erickcortez/tuipr/releases/tag/v0.2.0
[0.1.0]: https://github.com/erickcortez/tuipr/releases/tag/v0.1.0
