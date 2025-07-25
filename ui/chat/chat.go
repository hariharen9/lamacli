package chat

import (
	"fmt"
	os "os"
	"regexp"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
	glamour "github.com/charmbracelet/glamour"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hariharen9/lamacli/chathistory"
	"github.com/hariharen9/lamacli/llm"
	"github.com/hariharen9/lamacli/ui/styles"
)

// FileContextRequestMsg is a message to request file selection for chat context.
type FileContextRequestMsg struct{}

// llmResponseChunkMsg is a message that contains a chunk of the LLM's response.
type llmResponseChunkMsg string

// streamCompleteMsg is a message that indicates the LLM stream has completed.
type streamCompleteMsg struct{}

// errMsg is a message that contains an error.
type errMsg struct{ err error }

func (e errMsg) Error() string {
	return e.err.Error()
}

type Model struct {
	viewport        viewport.Model
	TextInput       textinput.Model
	llmClient       *llm.OllamaClient
	SelectedModel   string
	History         []string
	streaming       bool
	ready           bool
	responseChan    chan string
	err             error
	width           int
	height          int
	renderer        *glamour.TermRenderer
	codeBlocks      []string                 // Store extracted code blocks
	selectedCode    int                      // Currently selected code block for copying
	showCodeHelp    bool                     // Show code copy help
	ContextFileName string                   // Name of the file added to context
	currentSession  *chathistory.ChatSession // Current chat session for auto-saving

	// New field for chat templates
	chatTemplates    map[string]string
	selectedTemplate string
}

// New creates a new chat model.
func New(llmClient *llm.OllamaClient, selectedModel string) Model {
	ti := textinput.New()
	ti.Placeholder = "Type your message here..."
	ti.Focus()
	ti.CharLimit = 4096
	ti.PromptStyle = styles.PromptStyle()
	ti.TextStyle = styles.TextInputStyle()
	// Clear any initial value to prevent gibberish on macOS
	ti.SetValue("")

	// Apply platform-specific fixes
	if runtime.GOOS == "darwin" {
		// macOS-specific initialization to prevent input issues
		ti.SetValue("")
		// Force a clean state by resetting the cursor position
		ti.SetCursor(0)
	}

	vp := viewport.New(80, 20)

	// Initialize markdown renderer with custom styling
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(80),
	)

	// Add welcome message to history
	welcomeMessage := "Welcome to LamaCLI! 🦙✨\n\nI'm ready to help you with your questions. You can:\n• Ask me anything about programming, writing, or general topics\n• Use 'F' to browse files and 'M' to switch AI models\n• Use 'C' to copy code blocks when available\n• Press 'H' for detailed help and instructions\n• Press Ctrl+C to exit\n\nWhat would you like to know?"

	return Model{
		viewport:      vp,
		TextInput:     ti,
		llmClient:     llmClient,
		SelectedModel: selectedModel,
		History:       []string{"", welcomeMessage}, // Empty user message, then welcome
		renderer:      renderer,
		codeBlocks:    []string{},
		selectedCode:  0,
		showCodeHelp:  false,

		// Initialize chat templates
		chatTemplates: map[string]string{
			"Code Review":   "### Code Review Template\nReview the code provided below thoroughly, addressing the following aspects:\n1. **Readability**: Is the code easily understandable? Provide suggestions for improving readability if necessary.\n2. **Performance**: Identify any bottlenecks or optimizations that can be applied.\n3. **Best Practices**: Ensure the code adheres to language and industry best practices.\n4. **Errors and Bugs**: Highlight potential bugs or errors with recommendations for fixes.\n5. **Security**: Check for security vulnerabilities and suggest mitigations.\n\nPlease provide detailed comments or annotations within the code:\n\n\n```\n[Paste your code here] \n```",
			"Documentation": "### Documentation Template\nGenerate comprehensive documentation for the following code or API, including:\n1. **Overview**: A brief description of the functionality and purpose.\n2. **Usage**: Include code snippets demonstrating how to effectively use the code/API.\n3. **Input Parameters**: List all parameters, including types and descriptions.\n4. **Return Values**: Document return types and their meanings.\n5. **Examples**: Provide example use cases.\n6. **Notes**: Any additional information or caveats.\n\n\n```\n[Paste your code here]\n```",
			"Debugging":     "### Debugging Template\nHelp debug the code by diagnosing the issue described and suggesting solutions. Include:\n1. **Problem Description**: Elaborate on the encountered issue.\n2. **Expected Behavior**: Detail the expected outcome of the code.\n3. **Observed Behavior**: Describe what actually happens.\n4. **Error Messages**: Include any error messages or logs.\n5. **Analysis**: Provide an analysis of potential root causes.\n6. **Solutions**: Suggest fixes or alternative approaches.\n\nAdditional context or setup that may be useful:\n\n\n```\n[Paste your code here] \n```",
		},
		selectedTemplate: "",
	}
}

// extractCodeBlocks extracts code blocks from markdown text
func extractCodeBlocks(text string) []string {
	codeBlockRegex := regexp.MustCompile("```(?:[a-zA-Z0-9_+-]*\\n)?([\\s\\S]*?)```")
	matches := codeBlockRegex.FindAllStringSubmatch(text, -1)

	var codeBlocks []string
	for _, match := range matches {
		if len(match) > 1 {
			codeBlocks = append(codeBlocks, strings.TrimSpace(match[1]))
		}
	}
	return codeBlocks
}

// copySelectedCodeBlock copies the selected code block to clipboard
func (m *Model) copySelectedCodeBlock() error {
	if len(m.codeBlocks) == 0 || m.selectedCode >= len(m.codeBlocks) {
		return fmt.Errorf("no code block selected")
	}
	return clipboard.WriteAll(m.codeBlocks[m.selectedCode])
}

// SetModel updates the selected model without recreating the entire chat
func (m *Model) SetModel(selectedModel string) {
	m.SelectedModel = selectedModel
	// Keep existing history and UI state
}

// cycleTemplate cycles through the available chat templates.
func (m *Model) cycleTemplate() {
	templates := []string{"Code Review", "Documentation", "Debugging"}
	idx := -1
	for i, template := range templates {
		if template == m.selectedTemplate {
			idx = i
			break
		}
	}

	idx = (idx + 1) % len(templates)
	m.selectedTemplate = templates[idx]

	// Update the text input with the selected template
	m.TextInput.SetValue(m.chatTemplates[m.selectedTemplate])
	m.TextInput.CursorEnd()
}

// Reset clears the chat history while preserving the model and UI state
func (m *Model) Reset() {
	// Clear history but keep welcome message
	welcomeMessage := "Welcome to LamaCLI! 🦙✨\n\nI'm ready to help you with your questions. You can:\n• Ask me anything about programming, writing, or general topics\n• Use 'Alt+T' to switch between templates\n• Use 'F' to browse files and 'M' to switch AI models\n• Use 'C' to copy code blocks when available\n• Press 'H' for detailed help and instructions\n• Press Ctrl+C to exit\n\nWhat would you like to know?"
	m.History = []string{"", welcomeMessage}
	m.codeBlocks = []string{}
	m.selectedCode = 0
	m.showCodeHelp = false
	m.streaming = false
	m.err = nil
	m.renderViewport()
}

// Init is a command that can be run when the program starts.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model accordingly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Handle special keys BEFORE text input updates
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "alt+t":
			m.cycleTemplate()
			return m, nil
		default:
			if m.TextInput.Value() == "" {
				switch keyMsg.String() {
				case "C":
					if len(m.codeBlocks) > 0 {
						m.showCodeHelp = !m.showCodeHelp
						return m, nil
					}
				case "j", "down":
					if m.showCodeHelp && len(m.codeBlocks) > 0 {
						m.selectedCode = (m.selectedCode + 1) % len(m.codeBlocks)
						return m, nil
					}
				case "k", "up":
					if m.showCodeHelp && len(m.codeBlocks) > 0 {
						m.selectedCode = (m.selectedCode - 1 + len(m.codeBlocks)) % len(m.codeBlocks)
						return m, nil
					}
				}
			}
		}

		// Handle enter key for code copying
		if keyMsg.Type == tea.KeyEnter {
			if m.showCodeHelp && len(m.codeBlocks) > 0 && m.TextInput.Value() == "" {
				err := m.copySelectedCodeBlock()
				if err == nil {
					m.showCodeHelp = false
				}
				return m, nil
			}
		}
	}

	// Update TextInput and viewport after special key handling
	// Filter out potential problematic key sequences on macOS
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Skip certain problematic key sequences that might cause gibberish
		if keyMsg.Type == tea.KeyRunes {
			// Filter out non-printable characters that might cause display issues
			for _, r := range keyMsg.Runes {
				if r < 32 && r != 9 && r != 10 && r != 13 { // Allow tab, newline, carriage return
					// Skip this update if it contains problematic control characters
					return m, tea.Batch(cmds...)
				}
			}
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - styles.ChatBoxStyle().GetHorizontalFrameSize()
		m.viewport.Height = msg.Height - lipgloss.Height(m.TextInput.View()) - styles.ChatBoxStyle().GetVerticalFrameSize() - 2
		m.TextInput.Width = msg.Width - 8 // Adjust for padding and borders

		// Update renderer width
		if m.renderer != nil {
			m.renderer, _ = glamour.NewTermRenderer(
				glamour.WithStylePath("dark"),
				glamour.WithWordWrap(m.viewport.Width-4),
			)
		}

		if !m.ready {
			m.ready = true
			// Clean text input state on first ready event (macOS fix)
			if runtime.GOOS == "darwin" {
				currentValue := m.TextInput.Value()
				m.TextInput.SetValue("")
				m.TextInput.SetValue(currentValue)
			}
		}
		m.renderViewport()

	case llmResponseChunkMsg:
		if m.streaming {
			m.History[len(m.History)-1] += string(msg)
			m.renderViewport()
			m.viewport.GotoBottom()
			return m, readStreamCmd(m.responseChan)
		}

	case streamCompleteMsg:
		m.streaming = false
		m.responseChan = nil
		m.renderViewport()
		m.viewport.GotoBottom()
		// Auto-save session after response completion
		m.AutoSaveSession()

	case errMsg:
		m.err = msg.err
		m.streaming = false
		m.responseChan = nil

	case tea.KeyMsg:
		// Handle message sending
		switch msg.Type {
		case tea.KeyEnter:
			if m.streaming {
				return m, nil // Don't send new prompts while streaming
			}
			question := strings.TrimSpace(m.TextInput.Value())
			if question == "" {
				return m, nil
			}

			m.History = append(m.History, question)
			m.History = append(m.History, "") // Placeholder for LLM response
			m.TextInput.SetValue("")
			m.ContextFileName = "" // Clear the context file name
			m.renderViewport()
			m.viewport.GotoBottom()

			m.streaming = true
			m.err = nil // Clear previous errors
			m.responseChan = make(chan string)
			f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				// Handle error if log file can't be opened
				return m, nil
			}
			defer f.Close()
			f.WriteString(fmt.Sprintf("DEBUG: Calling GenerateResponseStream with model: %s\n", m.SelectedModel))
			go m.llmClient.GenerateResponseStream(m.SelectedModel, "You are a helpful assistant.", m.History[:len(m.History)-1], m.responseChan)
			return m, readStreamCmd(m.responseChan)

		case tea.KeyRunes:
			// Check for "@" to trigger file context selection
			if len(msg.Runes) == 1 && msg.Runes[0] == '@' {
				return m, func() tea.Msg {
					return FileContextRequestMsg{}
				}
			}

			// More aggressive filtering of problematic sequences on macOS
			if runtime.GOOS == "darwin" {
				// Check for ANSI escape sequences which often start with ESC ([)
				hasEscapeSequence := false
				for _, r := range msg.Runes {
					if r == 27 || r == 91 { // ESC or [
						hasEscapeSequence = true
						break
					}
				}
				if hasEscapeSequence {
					return m, nil
				}
			}

			// Filter out any remaining problematic sequences
			for _, r := range msg.Runes {
				if r < 32 && r != 9 && r != 10 && r != 13 {
					// Skip processing if control characters are present
					return m, nil
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// readStreamCmd waits for the next message from the stream.
func readStreamCmd(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		chunk, ok := <-ch
		if !ok {
			return streamCompleteMsg{}
		}
		if strings.HasPrefix(chunk, "Error: ") {
			return errMsg{err: fmt.Errorf("%s", strings.TrimPrefix(chunk, "Error: "))}
		}
		return llmResponseChunkMsg(chunk)
	}
}

func (m *Model) renderViewport() {
	var content strings.Builder

	// Clear existing code blocks
	m.codeBlocks = []string{}

	for i, line := range m.History {
		var styledLine string
		if i%2 == 0 {
			// User messages
			if line != "" { // Skip empty user messages (like welcome message prefix)
				userIcon := "👤"
				styledLine = styles.UserPromptStyle().Render(userIcon + " You: " + line)
			}
		} else {
			// LLM responses - render as markdown
			llmIcon := "🤖"
			if line != "" {
				// Extract code blocks before rendering
				codeBlocks := extractCodeBlocks(line)
				m.codeBlocks = append(m.codeBlocks, codeBlocks...)

				// Render markdown
				if m.renderer != nil {
					rendered, err := m.renderer.Render(line)
					if err == nil {
						styledLine = llmIcon + " LLM:\n" + rendered
					} else {
						// Fallback to plain text if markdown rendering fails
						styledLine = styles.LLMResponseStyle().Render(llmIcon + " LLM: " + line)
					}
				} else {
					styledLine = styles.LLMResponseStyle().Render(llmIcon + " LLM: " + line)
				}
			}
		}
		if styledLine != "" {
			content.WriteString(styledLine)
			content.WriteString("\n\n") // Add extra spacing between messages
		}
	}

	// Reset selected code block if we have fewer blocks now
	if m.selectedCode >= len(m.codeBlocks) {
		m.selectedCode = 0
	}

	m.viewport.SetContent(content.String())
}

// View returns the string representation of the UI.
func (m Model) View() string {
	var view strings.Builder

	// Enhanced header with prominent model indicator
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.TitleStyle().GetForeground()).
		Padding(0, 2).
		MarginBottom(1).
		Background(lipgloss.Color("235")).
		Width(m.width - 4)

	modelIcon := "🤖"
	statusIcon := ""
	if m.streaming {
		modelIcon = "⚡"
		statusIcon = " • 🔄 thinking..."
	} else {
		statusIcon = " • ✅ ready"
	}

	// Create a more prominent model indicator
	modelLabel := lipgloss.NewStyle().
		Foreground(styles.SubtleStyle().GetForeground()).
		Render("Current Model: ")

	modelName := lipgloss.NewStyle().
		Foreground(styles.TitleStyle().GetForeground()).
		Bold(true).
		Render(fmt.Sprintf("%s %s", modelIcon, m.SelectedModel))

	statusText := lipgloss.NewStyle().
		Foreground(styles.SubtleStyle().GetForeground()).
		Render(statusIcon)

	headerContent := lipgloss.JoinHorizontal(lipgloss.Left, modelLabel, modelName, statusText)
	header := headerStyle.Render(headerContent) + "\n"

	view.WriteString(header)

	// Enhanced chat box with better styling
	chatBoxHeight := m.height - lipgloss.Height(m.TextInput.View()) - lipgloss.Height(header) - 3
	enhancedChatBox := styles.ChatBoxStyle().
		Width(m.width).
		Height(chatBoxHeight).
		BorderForeground(styles.TitleStyle().GetForeground())

	view.WriteString(enhancedChatBox.Render(m.viewport.View()))
	view.WriteString("\n")

	// Enhanced input with prompt indicator
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, true, false, true).
		BorderForeground(styles.TitleStyle().GetForeground()).
		Padding(0, 1)

	promptIndicator := "💭"
	if m.streaming {
		promptIndicator = "⏳"
	}

	inputPrefix := lipgloss.NewStyle().
		Foreground(styles.PromptStyle().GetForeground()).
		Bold(true).
		Render(promptIndicator + " ")

	customInput := lipgloss.JoinHorizontal(lipgloss.Left, inputPrefix, m.TextInput.View())
	view.WriteString(inputStyle.Render(customInput))

	// Footer with error, status, or code block help
	if m.err != nil {
		errorFooter := lipgloss.NewStyle().
			Foreground(styles.ErrorStyle().GetForeground()).
			Bold(true).
			MarginTop(1).
			Render("❌ Error: " + m.err.Error())
		return lipgloss.JoinVertical(lipgloss.Left, view.String(), errorFooter)
	}

	// Show code block help if active
	if m.showCodeHelp && len(m.codeBlocks) > 0 {
		codeHelpStyle := lipgloss.NewStyle().
			Foreground(styles.StatusStyle().GetForeground()).
			Bold(true).
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.TitleStyle().GetForeground())

		codeHelpText := fmt.Sprintf(
			"📎 Code Blocks (%d/%d)\n↑/↓ or j/k: Navigate • Enter: Copy • C: Close",
			m.selectedCode+1, len(m.codeBlocks),
		)

		// Show preview of selected code block
		if m.selectedCode < len(m.codeBlocks) {
			preview := m.codeBlocks[m.selectedCode]
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			codeHelpText += "\n\n" + styles.SubtleStyle().Render("Preview: ") + preview
		}

		codeHelp := codeHelpStyle.Render(codeHelpText)
		return lipgloss.JoinVertical(lipgloss.Left, view.String(), codeHelp)
	}

	// Show code block indicator if blocks are available
	if len(m.codeBlocks) > 0 && !m.showCodeHelp {
		codeIndicator := lipgloss.NewStyle().
			Foreground(styles.SubtleStyle().GetForeground()).
			MarginTop(1).Render(fmt.Sprintf("📎 %d code block(s) available • Press C to copy", len(m.codeBlocks)))
		return lipgloss.JoinVertical(lipgloss.Left, view.String(), codeIndicator)
	}

	// Show attached file indicator
	if m.ContextFileName != "" {
		fileIndicator := lipgloss.NewStyle().
			Foreground(styles.SubtleStyle().GetForeground()).
			MarginTop(1).Render(fmt.Sprintf("📄 Attached: %s", m.ContextFileName))
		return lipgloss.JoinVertical(lipgloss.Left, view.String(), fileIndicator)
	}

	// Show template indicator if template is selected
	if m.selectedTemplate != "" {
		templateIndicator := lipgloss.NewStyle().
			Foreground(styles.SubtleStyle().GetForeground()).
			MarginTop(1).Render(fmt.Sprintf("📝 Template: %s • Press Alt+T to switch", m.selectedTemplate))
		return lipgloss.JoinVertical(lipgloss.Left, view.String(), templateIndicator)
	}

	return view.String()
}

// LoadFromSession loads a chat session into the current model
func (m *Model) LoadFromSession(session *chathistory.ChatSession) {
	m.History = make([]string, len(session.History))
	copy(m.History, session.History)
	m.SelectedModel = session.Model
	m.currentSession = session
	m.codeBlocks = []string{}
	m.selectedCode = 0
	m.showCodeHelp = false
	m.streaming = false
	m.err = nil
	m.renderViewport()
}

// SaveToSession saves the current chat to a session
func (m *Model) SaveToSession() (*chathistory.ChatSession, error) {
	historyManager, err := chathistory.NewChatHistoryManager()
	if err != nil {
		return nil, err
	}

	if m.currentSession == nil {
		// Create new session
		m.currentSession = &chathistory.ChatSession{
			Model:   m.SelectedModel,
			History: make([]string, len(m.History)),
		}
		copy(m.currentSession.History, m.History)
	} else {
		// Update existing session
		m.currentSession.Model = m.SelectedModel
		m.currentSession.History = make([]string, len(m.History))
		copy(m.currentSession.History, m.History)
	}

	if err := historyManager.SaveSession(m.currentSession); err != nil {
		return nil, err
	}

	return m.currentSession, nil
}

// AutoSaveSession automatically saves the session after each response
func (m *Model) AutoSaveSession() {
	// Only auto-save if we have meaningful conversation (more than just welcome message)
	if len(m.History) > 2 {
		go func() {
			m.SaveToSession()
		}()
	}
}

// GetCurrentSession returns the current session
func (m *Model) GetCurrentSession() *chathistory.ChatSession {
	return m.currentSession
}
