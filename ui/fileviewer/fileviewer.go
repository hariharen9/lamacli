package fileviewer

import (
	"lamacli/fileops"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"lamacli/ui/styles"
)

// Model is a bubbletea model for viewing a file.
// It uses the viewport component from the bubbles library.
// https://github.com/charmbracelet/bubbles/tree/master/viewport
type Model struct {
	Viewport viewport.Model
}

// New creates a new file viewer model.
func New() Model {
	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.TitleStyle.GetForeground()).
		PaddingRight(2)

	return Model{Viewport: vp}
}

// SetContent sets the content of the file viewer.
func (m *Model) SetContent(path string) error {
	content, err := fileops.ReadFile(path)
	if err != nil {
		return err
	}

	m.Viewport.SetContent(string(content))
	return nil
}
