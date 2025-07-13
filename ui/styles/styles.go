package styles

import (
	"github.com/hariharen9/lamacli/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

func AppStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Margin(0, 0)
}

func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Primary).
		Bold(true)
}

func SubtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Subtle)
}

func ItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(4)
}

func SelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(2).Foreground(theme.CurrentTheme.Primary)
}

func ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Error).
		Bold(true)
}

func PromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.CurrentTheme.Primary)
}

func TextInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
}

func UserPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(theme.CurrentTheme.Secondary).Bold(true)
}

func LLMResponseStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
}

func ChatBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary).
		Padding(1, 2).
		Margin(1, 0)
}

func WelcomeStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Padding(2, 4)
}

func StatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Success).
		Bold(true)
}