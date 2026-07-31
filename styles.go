package main

import "github.com/charmbracelet/lipgloss"

// Colores Catppuccin Mocha
var (
	ColorBase     = lipgloss.Color("#1e1e2e")
	ColorMantle   = lipgloss.Color("#181825")
	ColorSurface0 = lipgloss.Color("#313244")
	ColorSurface1 = lipgloss.Color("#45475a")
	ColorSurface2 = lipgloss.Color("#585b70")
	ColorText     = lipgloss.Color("#cdd6f4")
	ColorSubtext0 = lipgloss.Color("#a6adc8")
	ColorOverlay0 = lipgloss.Color("#6c7086")
	ColorMauve    = lipgloss.Color("#cba6f7")
	ColorBlue     = lipgloss.Color("#89b4fa")
	ColorSapphire = lipgloss.Color("#74c7ec")
	ColorGreen    = lipgloss.Color("#a6e3a1")
	ColorYellow   = lipgloss.Color("#f9e2af")
	ColorPeach    = lipgloss.Color("#fab387")
	ColorRed      = lipgloss.Color("#f38ba8")
	ColorPink     = lipgloss.Color("#f5c2e7")
)

// getPanelStyle retorna el estilo del panel según si está activo y la configuración
func getPanelStyle(active, transparent bool) lipgloss.Style {
	if active {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(ColorMauve).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			Padding(0, 1).
			Margin(0, 0)
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ColorSurface0).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		Padding(0, 1).
		Margin(0, 0)

	if transparent {
		// Fondo transparente en paneles inactivos
		style = style.Background(lipgloss.Color("#00000000"))
	} else {
		style = style.Background(ColorBase)
	}

	return style
}

// Estilos base
var (
	// Paneles activos/inactivos (se usan con getPanelStyle)
	ActivePanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(ColorMauve).
				Padding(0, 1)

	InactivePanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(ColorSurface0).
				Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	BranchStyle = lipgloss.NewStyle().
			Foreground(ColorSapphire)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorMauve).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext0)

	PRNumberStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	KeybindStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay0)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	NormalModeStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	InsertModeStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(ColorText).
				Background(ColorSurface0)

	InsertInputStyle = lipgloss.NewStyle().
				Foreground(ColorText).
				Background(ColorSurface0)

	SeparatorStyle = lipgloss.NewStyle().
			Foreground(ColorSurface0)

	PanelActiveStyle = lipgloss.NewStyle().
			Foreground(ColorMauve)

	PanelInactiveStyle = lipgloss.NewStyle().
			Foreground(ColorSurface0)

	TitleTextStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	GreenStyle  = lipgloss.NewStyle().Foreground(ColorGreen)
	RedStyle    = lipgloss.NewStyle().Foreground(ColorRed)
	YellowStyle = lipgloss.NewStyle().Foreground(ColorYellow)
)
