package theme

import (
	catppuccingo "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color palette for the application.
// We can add more colors here as we need them.
type Theme struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Subtle    lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
}

// DefaultTheme is the default theme for the application.
var DefaultTheme = Theme{
	Primary:   lipgloss.Color(catppuccingo.Mocha.Mauve().Hex),
	Secondary: lipgloss.Color(catppuccingo.Mocha.Pink().Hex),
	Subtle:    lipgloss.Color(catppuccingo.Mocha.Subtext0().Hex),
	Success:   lipgloss.Color(catppuccingo.Mocha.Green().Hex),
	Warning:   lipgloss.Color(catppuccingo.Mocha.Yellow().Hex),
	Error:     lipgloss.Color(catppuccingo.Mocha.Red().Hex),
}
