# Contributing to tuipr

🎉 Thanks for your interest in contributing to tuipr!

## Development Setup

```bash
git clone https://github.com/erickcortez/tuipr.git
cd tuipr
go build -o tuipr .
```

## Project Structure

```
tuipr/
├── main.go      # Entry point
├── model.go     # State and messages
├── update.go    # Key handlers and actions
├── view.go      # TUI rendering
├── styles.go    # Catppuccin Mocha theme
└── flags.go     # CLI flags
```

## Making Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes and test
4. Build: `go build -o tuipr .`
5. Commit your changes
6. Push and open a Pull Request

## Code Style

- Run `go fmt` before committing
- Run `go vet` to check for issues
- Add documentation for new functions

## Reporting Issues

Please report bugs via GitHub Issues with:
- Your OS and terminal
- Steps to reproduce
- Expected vs actual behavior

---

Made with ❤️ by Erick Cortez and Dereck Curtis
