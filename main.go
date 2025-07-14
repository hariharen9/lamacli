package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hariharen9/lamacli/cli"
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
	// Check if command line arguments are provided for CLI mode
	if len(os.Args) > 1 {
		// Handle CLI commands
		if err := cli.ProcessCLICommand(os.Args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// No arguments provided - start interactive mode
	startInteractiveMode()
}

// startInteractiveMode initializes and runs the interactive TUI
func startInteractiveMode() {
	// Set terminal to raw mode on macOS to help prevent escape sequence issues
	if runtime.GOOS == "darwin" {
		fmt.Print("\033c") // Clear terminal to reset state
	}
	
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
