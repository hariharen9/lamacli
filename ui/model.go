package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lamacli/fileops"
	"lamacli/llm"
	"lamacli/ui/chat"
	"lamacli/ui/filetree"
	"lamacli/ui/fileviewer"
	"lamacli/ui/modelselect"
	"lamacli/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// viewMode is used to determine which view to render.
type viewMode int

const (
	chatView viewMode = iota // Chat is now the default view
	fileTreeView
	fileViewerView
	modelSelectView
)

// Msg for when an error occurs
type errMsg struct{ error }

// Model represents the state of our UI.
type Model struct {
	filetree      *filetree.Model
	fileviewer    fileviewer.Model
	modelselect   *modelselect.Model
	chat          chat.Model
	llmClient     *llm.OllamaClient
	viewMode      viewMode
	width         int
	height        int
	selectedModel string
	Err           error // Stores errors to display to the user
}

// InitialModel returns an initialized Model.
func InitialModel() Model {
	ft, err := filetree.New(".")
	if err != nil {
		panic(err)
	}

	llmClient, err := llm.NewOllamaClient()
	var initialErr error
	if err != nil {
		initialErr = fmt.Errorf("Ollama client initialization failed: %w", err)
	}

	ms, err := modelselect.New(llmClient)
	if err != nil && initialErr == nil {
		initialErr = fmt.Errorf("Model selection initialization failed: %w", err)
	}

	var defaultModel string
	if llmClient != nil {
		models, err := llmClient.ListModels()
		if err == nil && len(models) > 0 {
			defaultModel = models[0] // Set first available model as default
		} else if initialErr == nil {
			initialErr = fmt.Errorf("No Ollama models found. Please pull a model (e.g., 'ollama pull llama2').")
		}
	}

	return Model{
		filetree:      ft,
		fileviewer:    fileviewer.New(),
		modelselect:   ms,
		chat:          chat.New(llmClient, defaultModel),
		llmClient:     llmClient,
		viewMode:      chatView, // Start with chat view
		selectedModel: defaultModel,
		Err:           initialErr,
	}
}

// Init is a command that can be run when the program starts.
func (m Model) Init() tea.Cmd {
	// Initialize the chat view since it's the default
	return m.chat.Init()
}

// Update handles messages and updates the model accordingly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.filetree.List.SetSize(msg.Width, msg.Height-styles.AppStyle.GetVerticalFrameSize()-lipgloss.Height(m.helpView()))
		m.fileviewer.Viewport.Width = msg.Width - styles.AppStyle.GetHorizontalFrameSize()
		m.fileviewer.Viewport.Height = msg.Height - styles.AppStyle.GetVerticalFrameSize() - lipgloss.Height(m.helpView())

	

	case errMsg:
		m.Err = msg
		return m, nil

	case tea.KeyMsg:
		// Handle global shortcuts first
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			// Only quit if not in chat view or if chat input is empty
			if m.viewMode != chatView || m.chat.TextInput.Value() == "" {
				return m, tea.Quit
			}

		case "enter":
			if m.viewMode == fileTreeView {
				selectedItem := m.filetree.List.SelectedItem().(filetree.Item)
				currentPath := filepath.Join(m.filetree.List.Title, selectedItem.Path)
				if selectedItem.IsDir {
					m.filetree.GoTo(currentPath)
				} else {
					m.fileviewer.SetContent(currentPath)
					m.viewMode = fileViewerView
				}
			} else if m.viewMode == modelSelectView {
				// Handle model selection
				m.selectedModel = m.modelselect.GetSelectedModel()
				m.viewMode = chatView // Go back to chat view after selection
				// Re-initialize chat with the new model
				m.chat = chat.New(m.llmClient, m.selectedModel)
				return m, nil
			}

		case "backspace":
			if m.viewMode == fileTreeView {
				currentPath := m.filetree.List.Title
				parentPath := filepath.Dir(currentPath)
				m.filetree.GoTo(parentPath)
			} else if m.viewMode == fileViewerView {
				m.viewMode = fileTreeView
			}
			
		case "escape":
			// Escape should only navigate between views, not exit the program
			if m.viewMode == fileTreeView {
				m.viewMode = chatView // Go back to chat view
				return m, nil
			} else if m.viewMode == fileViewerView {
				m.viewMode = chatView // Go back to chat view
				return m, nil
			} else if m.viewMode == modelSelectView {
				m.viewMode = chatView // Go back to chat view
				return m, nil
			}
			// In chat view, escape does nothing (we don't want to exit)

		case "delete":
			if m.viewMode == fileTreeView {
				selectedItem := m.filetree.List.SelectedItem().(filetree.Item)
				filePath := filepath.Join(m.filetree.List.Title, selectedItem.Path)
				go func() {
					fileops.DeleteFile(filePath)
					// Refresh the file tree after deletion
					m.filetree.GoTo(m.filetree.List.Title)
				}()
			}

		case "f", "F": // 'F' for file tree view
			// Only allow if not typing in chat
			if m.viewMode != chatView || m.chat.TextInput.Value() == "" {
				m.viewMode = fileTreeView
				return m, nil
			}

		case "m", "M": // 'M' for model selection
			// Only allow if not typing in chat
			if m.viewMode != chatView || m.chat.TextInput.Value() == "" {
				m.viewMode = modelSelectView
				// If there was an error initializing modelselect, we should display it.
				if m.modelselect == nil {
					return m, func() tea.Msg { return errMsg{fmt.Errorf("Ollama client not initialized. Please ensure Ollama is running.")} }
				}
				// Set the current model as the default selection
				if m.selectedModel != "" {
					m.modelselect.SetSelectedModel(m.selectedModel)
				}
				return m, m.modelselect.Init()
			}
			
		case "r", "R": // 'R' for reset/clear chat
			// Only allow if not typing in chat
			if m.viewMode == chatView && m.chat.TextInput.Value() == "" {
				// Clear chat history and reinitialize
				m.chat = chat.New(m.llmClient, m.selectedModel)
				return m, m.chat.Init()
			}
		}
	}

	var cmd tea.Cmd
	switch m.viewMode {
	case fileTreeView:
		m.filetree.List, cmd = m.filetree.List.Update(msg)
	case fileViewerView:
		m.fileviewer.Viewport, cmd = m.fileviewer.Viewport.Update(msg)
	case modelSelectView:
		// Only update modelselect if it's not nil (i.e., no init error)
		if m.modelselect != nil {
			updatedModel, updateCmd := m.modelselect.Update(msg)
			if ms, ok := updatedModel.(*modelselect.Model); ok {
				m.modelselect = ms
			if m.modelselect.FormCompleted() {
				m.selectedModel = m.modelselect.GetSelectedModel()
				f, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				log.SetOutput(f)
				log.Printf("DEBUG: Model selection completed, selected model: %s", m.selectedModel)
				f.Close()
				m.chat = chat.New(m.llmClient, m.selectedModel) // Re-initialize chat with the new model
				m.viewMode = fileTreeView
				return m, nil
			}
			}
			cmd = updateCmd
		}
	case chatView:
		updatedChat, newCmd := m.chat.Update(msg)
		m.chat = updatedChat.(chat.Model)
		cmd = newCmd
	}

	return m, cmd
}

// helpView returns the help text for the current view with enhanced styling.
func (m Model) helpView() string {
	var helpItems []string
	var title string
	
	switch m.viewMode {
	case chatView:
		title = "💬 Chat"
		helpItems = []string{
			"↑/↓: scroll history",
			"enter: send message",
			"F: file explorer",
			"M: switch model",
			"R: reset chat",
			"C: copy code blocks",
			"ctrl+c: exit",
		}
	case fileTreeView:
		title = "📁 File Explorer"
		helpItems = []string{
			"↑/↓: navigate",
			"enter: open file/folder",
			"backspace: parent folder",
			"del: delete file",
			"esc: back to chat",
			"ctrl+c: exit",
		}
	case fileViewerView:
		title = "📄 File Viewer"
		helpItems = []string{
			"↑/↓: scroll",
			"backspace: back to explorer",
			"esc: back to chat",
			"ctrl+c: exit",
		}
	case modelSelectView:
		title = "🤖 Model Selection"
		helpItems = []string{
			"↑/↓: navigate models",
			"enter: select model",
			"esc: back to chat",
			"ctrl+c: exit",
		}
	}
	
	// Create a more appealing help bar
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.TitleStyle.GetForeground()).
		Bold(true).
		PaddingRight(2)
	
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.SubtleStyle.GetForeground()).
		PaddingLeft(1).
		PaddingRight(1)
	
	separatorStyle := lipgloss.NewStyle().
		Foreground(styles.SubtleStyle.GetForeground()).
		SetString(" • ")
	
	var helpBar []string
	helpBar = append(helpBar, titleStyle.Render(title))
	
	for i, item := range helpItems {
		if i > 0 {
			helpBar = append(helpBar, separatorStyle.Render())
		}
		helpBar = append(helpBar, helpStyle.Render(item))
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Left, helpBar...)
}

// View returns the string representation of the UI.
func (m Model) View() string {
	var s string
	switch m.viewMode {
	case fileTreeView:
		s = m.filetree.View()
	case fileViewerView:
		s = m.fileviewer.Viewport.View()
	case modelSelectView:
		s = m.modelselect.View()
	case chatView:
		s = m.chat.View()
	}

	// Display error message if any
	if m.Err != nil {
		errorView := lipgloss.JoinVertical(
			lipgloss.Top,
			styles.ErrorStyle.Render(fmt.Sprintf("❌ Error: %v", m.Err)),
			styles.SubtleStyle.Render("Please ensure Ollama is running and you have models available."),
			styles.SubtleStyle.Render("Try running: ollama pull llama3.2:3b"),
		)
		return lipgloss.JoinVertical(
			lipgloss.Top,
			styles.AppStyle.Render(errorView),
			m.helpView(),
		)
	}

	// Main content with better spacing
	mainContent := styles.AppStyle.Render(s)
	helpBar := m.helpView()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		mainContent,
		helpBar,
	)
}

