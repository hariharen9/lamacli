package styles

import (
	"github.com/hariharen9/lamacli/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

// AppStyle is the base style for the application.
var AppStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Margin(0, 0)

// TitleStyle is the style for titles.
var TitleStyle = lipgloss.NewStyle().
	Foreground(theme.DefaultTheme.Primary).
	Bold(true)

// SubtleStyle is a style for subtle text.
var SubtleStyle = lipgloss.NewStyle().
	Foreground(theme.DefaultTheme.Subtle)

// ItemStyle is the style for a list item.
var ItemStyle = lipgloss.NewStyle().PaddingLeft(4)

// SelectedItemStyle is the style for a selected list item.
var SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(theme.DefaultTheme.Primary)

// ErrorStyle is the style for error messages.
var ErrorStyle = lipgloss.NewStyle().
	Foreground(theme.DefaultTheme.Error).
	Bold(true)

// PromptStyle is the style for the text input prompt.
var PromptStyle = lipgloss.NewStyle().Foreground(theme.DefaultTheme.Primary)

// TextInputStyle is the style for the text input.
var TextInputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

// UserPromptStyle is the style for the user's prompt in chat.
var UserPromptStyle = lipgloss.NewStyle().Foreground(theme.DefaultTheme.Secondary).Bold(true)

// LLMResponseStyle is the style for the LLM's response in chat.
var LLMResponseStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

// ChatBoxStyle is the style for the chat box.
var ChatBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(theme.DefaultTheme.Primary).
	Padding(1, 2).
	Margin(1, 0)

// WelcomeStyle is the style for welcome messages.
var WelcomeStyle = lipgloss.NewStyle().
	Foreground(theme.DefaultTheme.Primary).
	Bold(true).
	Align(lipgloss.Center).
	Padding(2, 4)

// StatusStyle is the style for status messages.
var StatusStyle = lipgloss.NewStyle().
	Foreground(theme.DefaultTheme.Success).
	Bold(true)
