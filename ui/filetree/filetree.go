package filetree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"lamacli/ui/styles"
)

// Item represents a file or directory in the file tree.
type Item struct {
	Path  string
	IsDir bool
}

// FilterValue is used by the list component to filter items.
func (i Item) FilterValue() string { return i.Path }

// itemDelegate is a custom delegate for the list component.
type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	// Render the item's path.
	// Add a trailing slash to directories.
	str := i.Path
	if i.IsDir {
		str = i.Path + "/"
	}

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

// Model is a bubbletea model for displaying a file tree.
type Model struct {
	List list.Model
	path string
}

// New creates a new file tree model.
func New(path string) (*Model, error) {
	m := &Model{path: path}
	if err := m.readDir(path); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Model) readDir(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	items := make([]list.Item, len(files))
	for i, file := range files {
		items[i] = Item{Path: file.Name(), IsDir: file.IsDir()}
	}

	if m.List.Items() == nil {
		// First time setup
		l := list.New(items, itemDelegate{}, 0, 0)
		l.SetShowHelp(false) // We'll manage help text ourselves.
		m.List = l
	} else {
		m.List.SetItems(items)
	}

	m.path, _ = filepath.Abs(path)
	m.List.Title = m.path
	m.List.Styles.Title = styles.TitleStyle.Copy().
		// Show a shorter path in the title
		MaxWidth(50).
		Italic(true)

	return nil
}

// GoTo changes the directory view.
func (m *Model) GoTo(path string) error {
	return m.readDir(path)
}

// View returns the string representation of the UI.
func (m *Model) View() string {
	return m.List.View()
}
