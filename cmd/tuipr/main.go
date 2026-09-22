package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
)

// ErrNotTerminal is returned when tuipr is run outside a terminal.
var ErrNotTerminal = errors.New("tuipr must be run in a terminal")

// version is set at build time via ldflags: -X main.version=v1.0.0
var version = "dev"

func main() {
	// Parse flags before starting the TUI.
	flags := parseFlags()

	// Handle --version / -v: print version and exit.
	if flags.Version {
		fmt.Printf("tuipr %s\n", version)
		os.Exit(0)
	}

	// Verify we are in a terminal.
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		println("Error:", ErrNotTerminal.Error())
		os.Exit(1)
	}

	// Create the initial model.
	model := NewModel(flags)

	// Start the Bubbletea program in AltScreen mode.
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		println("Error starting tuipr:", err.Error())
		os.Exit(1)
	}
}
