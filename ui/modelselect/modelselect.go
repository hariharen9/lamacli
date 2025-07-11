package modelselect

import (
	"fmt"
	"log"
	"os"

	"lamacli/llm"
	"lamacli/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type modelSelectedMsg struct {
	model string
}

// Model represents the state of the model selection UI.
type Model struct {
	llmClient     *llm.OllamaClient
	SelectedModel string
	form          *huh.Form
}

// New creates a new model selection model.
func New(llmClient *llm.OllamaClient) (*Model, error) {
	models, err := llmClient.ListModels()
	if err != nil {
		return nil, fmt.Errorf("failed to list Ollama models: %w", err)
	}

	options := make([]huh.Option[string], len(models))
	for i, model := range models {
		options[i] = huh.NewOption(model, model)
	}

	var selectedModel string
	ms := &Model{
		llmClient: llmClient,
		SelectedModel: selectedModel,
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("selectedModel").
				Title("Select an Ollama Model").
				Options(options...).
				Value(&ms.SelectedModel),
		),
	).WithTheme(huh.ThemeBase16())

	ms.form = form

	return ms, nil
}

// Init is a command that can be run when the program starts.
func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

// Update handles messages and updates the model accordingly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for escape key before passing to form
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "escape" {
			// Return a special message to signal escape
			return m, func() tea.Msg {
				return tea.KeyMsg{Type: tea.KeyEscape}
			}
		}
	}
	
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)
	
	// Update the SelectedModel when the form is completed
	if m.form.State == huh.StateCompleted {
		selectedVal := m.form.Get("selectedModel")
		if selectedStr, ok := selectedVal.(string); ok {
			m.SelectedModel = selectedStr
			f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				log.SetOutput(f)
				log.Printf("DEBUG: Model selected: %s", selectedStr)
				f.Close()
			}
		}
	}
	
	return m, cmd
}

// FormCompleted returns true if the user has submitted the form.
func (m *Model) FormCompleted() bool {
	return m.form.State == huh.StateCompleted
}

// View returns the string representation of the UI.
func (m Model) View() string {
	return styles.AppStyle.Render(m.form.View())
}

// GetSelectedModel returns the currently selected model from the form.
func (m *Model) GetSelectedModel() string {
	selectedVal := m.form.Get("selectedModel")
	if selectedStr, ok := selectedVal.(string); ok {
		f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			log.SetOutput(f)
			log.Printf("DEBUG: GetSelectedModel returning from form: %s", selectedStr)
			f.Close()
		}
		return selectedStr
	}
	// If form value is not available, return the model stored in the struct
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(f)
		log.Printf("DEBUG: GetSelectedModel returning from struct: %s", m.SelectedModel)
		f.Close()
	}
	return m.SelectedModel
}

// SetSelectedModel sets the selected model in the form
func (m *Model) SetSelectedModel(model string) {
	m.SelectedModel = model
	// The form will use the SelectedModel field as its value since it's bound with Value(&ms.SelectedModel)
}
