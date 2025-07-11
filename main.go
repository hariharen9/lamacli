package main

import (
	"fmt"
	"os"

	"github.com/hariharen9/lamacli/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const lamaPortrait = `

      ___      ___          ___          ___          ___          ___          
     /\__\    /\  \        /\__\        /\  \        /\  \        /\__\      ___   
    /:/  /   /::\  \      /::|  |      /::\  \      /::\  \      /:/  /     /\  \  
   /:/  /   /:/\:\  \    /:|:|  |     /:/\:\  \    /:/\:\  \    /:/  /      \:\  \ 
  /:/  /   /::\~\:\  \  /:/|:|__|__  /::\~\:\  \  /:/  \:\  \  /:/  /       /::\__\
 /:/__/   /:/\:\ \:\__\/:/ |::::\__\/:/\:\ \:\__\/:/__/ \:\__\/:/__/     __/:/\/__/
 \:\  \   \/__\:\/:/  /\/__/~~/:/  /\/__\:\/:/  /\:\  \  \/__/\:\  \    /\/:/  /   
  \:\  \       \::/  /       /:/  /      \::/  /  \:\  \       \:\  \   \::/__/    
   \:\  \      /:/  /       /:/  /       /:/  /    \:\  \       \:\  \   \:\__\    
    \:\__\    /:/  /       /:/  /       /:/  /      \:\__\       \:\__\   \/__/    
     \/__/    \/__/        \/__/        \/__/        \/__/        \/__/         

`

func main() {
	initialModel := ui.InitialModel()
	if initialModel.Err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F28482")). // A reddish color for errors
			Align(lipgloss.Center).
			Bold(true)

		message := fmt.Sprintf("Error: Please ensure Ollama is running and accessible with atleast one model pulled.")

		fmt.Println(errorStyle.Render(lamaPortrait))
		fmt.Println(errorStyle.Render(message))
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
