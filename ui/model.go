package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hariharen9/lamacli/fileops"
	"github.com/hariharen9/lamacli/llm"
	"github.com/hariharen9/lamacli/ui/chat"
	"github.com/hariharen9/lamacli/ui/filetree"
	"github.com/hariharen9/lamacli/ui/fileviewer"
	"github.com/hariharen9/lamacli/ui/modelselect"
	"github.com/hariharen9/lamacli/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// viewMode is used to determine which view to render.
type viewMode int

// fileSelectedMsg is a message to indicate a file has been selected for context.
type fileSelectedMsg struct {
	path    string
	content []byte
}

const (
	chatView viewMode = iota // Chat is now the default view
	fileTreeView
	fileViewerView
	modelSelectView
	helpView
)

// Msg for when an error occurs
type errMsg struct{ error }

// Model represents the state of our UI.
type Model struct {
	filetree         *filetree.Model
	fileviewer       fileviewer.Model
	modelselect      *modelselect.Model
	chat             chat.Model
	llmClient        *llm.OllamaClient
	viewMode         viewMode
	width            int
	height           int
	selectedModel    string
	fileContextMode  bool  // True when selecting a file for chat context
	exitConfirmation bool  // True when waiting for exit confirmation
	Err              error // Stores errors to display to the user
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
			initialErr = fmt.Errorf("No Ollama models found. Please pull a model (e.g., 'ollama pull llama2')")
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

	case chat.FileContextRequestMsg:
		m.viewMode = fileTreeView
		m.fileContextMode = true
		return m, nil

	case fileSelectedMsg:
		// Add file content to chat input, replacing the @
		currentInput := m.chat.TextInput.Value()
		newInput := strings.TrimSuffix(currentInput, "@") + fmt.Sprintf("\n--- Start of File: %s ---\n%s\n--- End of File ---\n", filepath.Base(msg.path), string(msg.content))
		m.chat.TextInput.SetValue(newInput)
		m.chat.ContextFileName = filepath.Base(msg.path)
		m.viewMode = chatView
		m.fileContextMode = false
		return m, nil

	case errMsg:
		m.Err = msg
		return m, nil

	case tea.KeyMsg:
		// Centralized escape handling
		if msg.Type == tea.KeyEscape || msg.String() == "escape" {
			if m.exitConfirmation {
				m.exitConfirmation = false
				return m, nil
			}
			if m.fileContextMode {
				m.viewMode = chatView
				m.fileContextMode = false
			} else if m.viewMode != chatView {
				m.viewMode = chatView
			}
			return m, nil
		}

		// Handle other keys based on the current view
		switch m.viewMode {
		case chatView:
			// Let the chat view handle its own updates
		case fileTreeView:
			switch msg.String() {
			case "enter":
				selectedItem := m.filetree.List.SelectedItem().(filetree.Item)
				currentPath := filepath.Join(m.filetree.List.Title, selectedItem.Path)
				if selectedItem.IsDir {
					m.filetree.GoTo(currentPath)
				} else { // It's a file
					if m.fileContextMode {
						content, err := fileops.ReadFile(currentPath)
						if err != nil {
							return m, func() tea.Msg { return errMsg{err} }
						}
						return m, func() tea.Msg {
							return fileSelectedMsg{path: currentPath, content: content}
						}
					} else {
						m.fileviewer.SetContent(currentPath)
						m.viewMode = fileViewerView
					}
				}
			case "backspace":
				currentPath := m.filetree.List.Title
				parentPath := filepath.Dir(currentPath)
				m.filetree.GoTo(parentPath)
			case "delete":
				selectedItem := m.filetree.List.SelectedItem().(filetree.Item)
				filePath := filepath.Join(m.filetree.List.Title, selectedItem.Path)
				go func() {
					fileops.DeleteFile(filePath)
					// Refresh the file tree after deletion
					m.filetree.GoTo(m.filetree.List.Title)
				}()
			}
		}

		// Global shortcuts that are not escape
		switch msg.String() {
		case "ctrl+c":
			if m.exitConfirmation {
				return m, tea.Quit
			} else {
				m.exitConfirmation = true
				return m, nil
			}
		case "F":
			if m.viewMode != chatView || m.chat.TextInput.Value() == "" {
				m.viewMode = fileTreeView
				return m, nil
			}
		case "M":
			if m.viewMode != chatView || m.chat.TextInput.Value() == "" {
				m.viewMode = modelSelectView
				if m.modelselect == nil {
					return m, func() tea.Msg {
						return errMsg{fmt.Errorf("Ollama client not initialized. Please ensure Ollama is running.")}
					}
				}
				if m.selectedModel != "" {
					m.modelselect.SetSelectedModel(m.selectedModel)
				}
				return m, tea.Batch(m.modelselect.Init(), func() tea.Msg {
					return tea.WindowSizeMsg{Width: m.width, Height: m.height}
				})
			}
		case "R":
			if m.viewMode == chatView && m.chat.TextInput.Value() == "" {
				m.chat.Reset()
				return m, nil
			}
		case "H":
			if m.viewMode != chatView || m.chat.TextInput.Value() == "" {
				m.viewMode = helpView
				return m, nil
			}
		default:
			// Any other key press cancels the exit confirmation
			if m.exitConfirmation {
				m.exitConfirmation = false
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
					m.chat = chat.New(m.llmClient, m.selectedModel)
					m.viewMode = chatView
					return m, m.chat.Init()
				}
			}
			cmd = updateCmd
		}
	case chatView:
		updatedChat, newCmd := m.chat.Update(msg)
		m.chat = updatedChat.(chat.Model)
		cmd = newCmd
	case helpView:
		cmd = nil
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
			"H: help",
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
	case helpView:
		title = "🆘 Help"
		helpItems = []string{
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

// renderHelpContent returns the detailed help content
func (m Model) renderHelpContent() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.TitleStyle.GetForeground()).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 2)

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.TitleStyle.GetForeground()).
		Bold(true).
		PaddingTop(1)

	itemStyle := lipgloss.NewStyle().
		Foreground(styles.SubtleStyle.GetForeground()).
		PaddingLeft(2)

	keyStyle := lipgloss.NewStyle().
		Foreground(styles.StatusStyle.GetForeground()).
		Bold(true)

	var content strings.Builder

	// Main title
	content.WriteString(titleStyle.Render("🦙 LamaCLI - Complete User Guide"))
	content.WriteString("\n\n")

	// Chat Commands
	content.WriteString(headerStyle.Render("💬 Chat Commands"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• Type your message and press " + keyStyle.Render("Enter") + " to send"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• Use " + keyStyle.Render("↑/↓") + " to scroll through chat history"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• Chat supports full markdown rendering with syntax highlighting"))
	content.WriteString("\n\n")

	// Navigation Commands
	content.WriteString(headerStyle.Render("🧭 Navigation Commands"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("F") + " - Open file explorer to browse project files"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("M") + " - Switch between different AI models"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("R") + " - Reset chat history (clears conversation)"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("H") + " - Show this help screen"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("Esc") + " - Return to chat from any view"))
	content.WriteString("\n\n")

	// Code Block Features
	content.WriteString(headerStyle.Render("📎 Code Block Features"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("C") + " - Open code block copy mode when code is available"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("↑/↓") + " or " + keyStyle.Render("j/k") + " - Navigate between code blocks"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("Enter") + " - Copy selected code block to clipboard"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• Code blocks are automatically extracted from AI responses"))
	content.WriteString("\n\n")

	// File Explorer
	content.WriteString(headerStyle.Render("📁 File Explorer"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("↑/↓") + " - Navigate through files and folders"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("Enter") + " - Open file for viewing or enter directory"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("Backspace") + " - Go to parent directory"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("Del") + " - Delete selected file (with confirmation)"))
	content.WriteString("\n\n")

	// Tips and Tricks
	content.WriteString(headerStyle.Render("💡 Tips and Tricks"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• All command keys (F, M, R, C, H) are uppercase only"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• Lowercase letters are used for typing messages normally"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• The AI supports markdown, so you can ask for formatted responses"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• You can ask the AI to generate code in any programming language"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• Use the file explorer to understand project structure"))
	content.WriteString("\n\n")

	// Exit
	content.WriteString(headerStyle.Render("🚪 Exit"))
	content.WriteString("\n")
	content.WriteString(itemStyle.Render("• " + keyStyle.Render("Ctrl+C") + " - Exit the application (only way to quit)"))
	content.WriteString("\n\n")

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.SubtleStyle.GetForeground()).
		Align(lipgloss.Center).
		Padding(1, 2)
	content.WriteString(footerStyle.Render("Press Esc to return to chat"))

	return content.String()
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
	case helpView:
		s = m.renderHelpContent()
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

	if m.exitConfirmation {
		exitMessage := lipgloss.NewStyle().
			Foreground(styles.ErrorStyle.GetForeground()).
			Bold(true).
			Render("Are you sure you want to exit? Press Ctrl+C again to confirm.")
		helpBar = lipgloss.JoinVertical(lipgloss.Left, helpBar, exitMessage)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		mainContent,
		helpBar,
	)
}
