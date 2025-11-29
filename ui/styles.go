package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorText      = lipgloss.Color("#E0E0E0")
	ColorGhost     = lipgloss.Color("#606060")
	ColorSelected  = lipgloss.Color("#AF87FF")
	ColorHighlight = lipgloss.Color("#00AFFF")
	ColorDim       = lipgloss.Color("#444444")

	// Styles
	InputStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	GhostStyle = lipgloss.NewStyle().
			Foreground(ColorGhost)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorSelected).
				Bold(true)

	ItemStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	TagStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			PaddingLeft(1)
)
