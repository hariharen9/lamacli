package chathistory

import (
	"fmt"

	"github.com/hariharen9/lamacli/chathistory"
	"github.com/hariharen9/lamacli/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SessionItem represents a chat session in the list
type SessionItem struct {
	Session *chathistory.ChatSession
}

// FilterValue returns the value to filter by
func (s SessionItem) FilterValue() string {
	return s.Session.Title
}

// Title returns the title of the session
func (s SessionItem) Title() string {
	return s.Session.Title
}

// Description returns the description of the session
func (s SessionItem) Description() string {
	return s.Session.GetSessionSummary()
}

// SessionSelectedMsg is sent when a session is selected
type SessionSelectedMsg struct {
	Session *chathistory.ChatSession
}

// SessionDeletedMsg is sent when a session is deleted
type SessionDeletedMsg struct {
	SessionID string
}

// Model represents the chat history browser model
type Model struct {
	list           list.Model
	historyManager *chathistory.ChatHistoryManager
	width          int
	height         int
	err            error
}

// New creates a new chat history browser
func New() (*Model, error) {
	historyManager, err := chathistory.NewChatHistoryManager()
	if err != nil {
		return nil, err
	}

	// Create list
	l := list.New([]list.Item{}, NewItemDelegate(), 0, 0)
	l.Title = "📚 Chat History"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.TitleStyle()
	l.Styles.PaginationStyle = styles.SubtleStyle()
	l.Styles.HelpStyle = styles.SubtleStyle()

	m := &Model{
		list:           l,
		historyManager: historyManager,
	}

	// Load sessions
	if err := m.loadSessions(); err != nil {
		m.err = err
	}

	return m, nil
}

// NewItemDelegate creates a new item delegate for the list
func NewItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = true
	d.SetHeight(3)
	d.Styles.SelectedTitle = styles.SelectedItemStyle()
	d.Styles.SelectedDesc = styles.SelectedItemStyle().Copy().Foreground(styles.SubtleStyle().GetForeground())
	d.Styles.NormalTitle = styles.ItemStyle()
	d.Styles.NormalDesc = styles.ItemStyle().Copy().Foreground(styles.SubtleStyle().GetForeground())
	return d
}

// loadSessions loads all chat sessions from disk
func (m *Model) loadSessions() error {
	sessions, err := m.historyManager.ListSessions()
	if err != nil {
		return err
	}

	items := make([]list.Item, len(sessions))
	for i, session := range sessions {
		items[i] = SessionItem{Session: session}
	}

	m.list.SetItems(items)
	return nil
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width - styles.AppStyle().GetHorizontalFrameSize())
		m.list.SetHeight(msg.Height - styles.AppStyle().GetVerticalFrameSize() - 4)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selectedItem, ok := m.list.SelectedItem().(SessionItem); ok {
				return m, func() tea.Msg {
					return SessionSelectedMsg{Session: selectedItem.Session}
				}
			}
		case "delete", "d":
			if selectedItem, ok := m.list.SelectedItem().(SessionItem); ok {
				sessionID := selectedItem.Session.ID
				if err := m.historyManager.DeleteSession(sessionID); err == nil {
					// Reload sessions after deletion
					m.loadSessions()
					return m, func() tea.Msg {
						return SessionDeletedMsg{SessionID: sessionID}
					}
				}
			}
		case "r":
			// Refresh the list
			if err := m.loadSessions(); err != nil {
				m.err = err
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the chat history browser
func (m *Model) View() string {
	if m.err != nil {
		return styles.ErrorStyle().Render(fmt.Sprintf("Error: %v", m.err))
	}

	if len(m.list.Items()) == 0 {
		emptyMessage := styles.SubtleStyle().Render("No chat history found.\nStart a new conversation to create your first session!")
		return lipgloss.JoinVertical(
			lipgloss.Center,
			m.list.View(),
			"\n",
			emptyMessage,
		)
	}

	return m.list.View()
}

// GetSelectedSession returns the currently selected session
func (m *Model) GetSelectedSession() *chathistory.ChatSession {
	if selectedItem, ok := m.list.SelectedItem().(SessionItem); ok {
		return selectedItem.Session
	}
	return nil
}

// SaveCurrentSession saves the current chat session
func (m *Model) SaveCurrentSession(history []string, model string) error {
	session := &chathistory.ChatSession{
		Model:   model,
		History: history,
	}

	if err := m.historyManager.SaveSession(session); err != nil {
		return err
	}

	// Reload sessions to show the new one
	return m.loadSessions()
}

// SaveExistingSession updates an existing session
func (m *Model) SaveExistingSession(session *chathistory.ChatSession) error {
	if err := m.historyManager.SaveSession(session); err != nil {
		return err
	}

	// Reload sessions to show updates
	return m.loadSessions()
}
