package theme

import (
	catppuccingo "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color palette for the application.
type Theme struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Subtle    lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
}

var themes = []Theme{
	{
		// Mocha (Original)
		Primary:   lipgloss.Color(catppuccingo.Mocha.Mauve().Hex),
		Secondary: lipgloss.Color(catppuccingo.Mocha.Pink().Hex),
		Subtle:    lipgloss.Color(catppuccingo.Mocha.Subtext0().Hex),
		Success:   lipgloss.Color(catppuccingo.Mocha.Green().Hex),
		Warning:   lipgloss.Color(catppuccingo.Mocha.Yellow().Hex),
		Error:     lipgloss.Color(catppuccingo.Mocha.Red().Hex),
	},
	{
		// Latte (Original)
		Primary:   lipgloss.Color(catppuccingo.Latte.Mauve().Hex),
		Secondary: lipgloss.Color(catppuccingo.Latte.Pink().Hex),
		Subtle:    lipgloss.Color(catppuccingo.Latte.Subtext0().Hex),
		Success:   lipgloss.Color(catppuccingo.Latte.Green().Hex),
		Warning:   lipgloss.Color(catppuccingo.Latte.Yellow().Hex),
		Error:     lipgloss.Color(catppuccingo.Latte.Red().Hex),
	},
	{
		// Frappe (Original)
		Primary:   lipgloss.Color(catppuccingo.Frappe.Mauve().Hex),
		Secondary: lipgloss.Color(catppuccingo.Frappe.Pink().Hex),
		Subtle:    lipgloss.Color(catppuccingo.Frappe.Subtext0().Hex),
		Success:   lipgloss.Color(catppuccingo.Frappe.Green().Hex),
		Warning:   lipgloss.Color(catppuccingo.Frappe.Yellow().Hex),
		Error:     lipgloss.Color(catppuccingo.Frappe.Red().Hex),
	},
	{
		// Macchiato (Original)
		Primary:   lipgloss.Color(catppuccingo.Macchiato.Mauve().Hex),
		Secondary: lipgloss.Color(catppuccingo.Macchiato.Pink().Hex),
		Subtle:    lipgloss.Color(catppuccingo.Macchiato.Subtext0().Hex),
		Success:   lipgloss.Color(catppuccingo.Macchiato.Green().Hex),
		Warning:   lipgloss.Color(catppuccingo.Macchiato.Yellow().Hex),
		Error:     lipgloss.Color(catppuccingo.Macchiato.Red().Hex),
	},
	{
		// Blue Theme
		Primary:   lipgloss.Color("#007BFF"), // Bright Blue
		Secondary: lipgloss.Color("#6C757D"), // Gray
		Subtle:    lipgloss.Color("#ADB5BD"), // Light Gray
		Success:   lipgloss.Color("#28A745"), // Green
		Warning:   lipgloss.Color("#FFC107"), // Yellow
		Error:     lipgloss.Color("#DC3545"), // Red
	},
	{
		// Green Theme
		Primary:   lipgloss.Color("#28A745"), // Green
		Secondary: lipgloss.Color("#6C757D"), // Gray
		Subtle:    lipgloss.Color("#ADB5BD"), // Light Gray
		Success:   lipgloss.Color("#007BFF"), // Blue
		Warning:   lipgloss.Color("#FFC107"), // Yellow
		Error:     lipgloss.Color("#DC3545"), // Red
	},
	{
		// Yellow Theme
		Primary:   lipgloss.Color("#FFC107"), // Yellow
		Secondary: lipgloss.Color("#6C757D"), // Gray
		Subtle:    lipgloss.Color("#ADB5BD"), // Light Gray
		Success:   lipgloss.Color("#28A745"), // Green
		Warning:   lipgloss.Color("#007BFF"), // Blue
		Error:     lipgloss.Color("#DC3545"), // Red
	},
	{
		// Red Theme
		Primary:   lipgloss.Color("#DC3545"), // Red
		Secondary: lipgloss.Color("#6C757D"), // Gray
		Subtle:    lipgloss.Color("#ADB5BD"), // Light Gray
		Success:   lipgloss.Color("#28A745"), // Green
		Warning:   lipgloss.Color("#FFC107"), // Yellow
		Error:     lipgloss.Color("#007BFF"), // Blue
	},
	{
		// White Theme (Light Mode)
		Primary:   lipgloss.Color("#343A40"), // Dark Gray for text
		Secondary: lipgloss.Color("#6C757D"), // Gray
		Subtle:    lipgloss.Color("#ADB5BD"), // Light Gray
		Success:   lipgloss.Color("#28A745"), // Green
		Warning:   lipgloss.Color("#FFC107"), // Yellow
		Error:     lipgloss.Color("#DC3545"), // Red
	},
	{
		// Gold Theme
		Primary:   lipgloss.Color("#FFD700"), // Gold
		Secondary: lipgloss.Color("#6C757D"), // Gray
		Subtle:    lipgloss.Color("#ADB5BD"), // Light Gray
		Success:   lipgloss.Color("#28A745"), // Green
		Warning:   lipgloss.Color("#007BFF"), // Blue
		Error:     lipgloss.Color("#DC3545"), // Red
	},
}

var currentThemeIndex = 0

// CurrentTheme holds the active theme.
var CurrentTheme = themes[currentThemeIndex]

// NextTheme cycles to the next theme.
func NextTheme() {
	currentThemeIndex = (currentThemeIndex + 1) % len(themes)
	CurrentTheme = themes[currentThemeIndex]
}
