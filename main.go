package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
)

func main() {
	// Parsear flags antes de iniciar el TUI
	flags := parseFlags()

	// Verificar que estamos en un terminal
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		println("Error: tuipr debe ejecutarse en un terminal")
		os.Exit(1)
	}

	// Crear el modelo inicial
	model := NewModel(flags)

	// Iniciar el programa Bubbletea en modo AltScreen
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		println("Error al iniciar tuipr:", err.Error())
		os.Exit(1)
	}
}
