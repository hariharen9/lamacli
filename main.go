package main

import (
	"flag"
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
  /:/  /   /::\~\:\  \  /:/|:|__|__  /::\~\:\  \  /:/  \:\  \  /:/  /       /::\__
 /:/__/   /:/\:\ \:\__\/:/ |::::\__\/:/\:\ \:\__\/:/__/ \:\__\/:/__/     __/:/\/__/
 \:\  \   \/__\:\/:/  /\/__/~~/:/  /\/__\:\/:/  /\:\  \  \/__/\:\  \    /\/:/  /   
  \:\  \       \::/  /       /:/  /      \::/  /  \:\  \       \:\  \   \::/__/    
   \:\  \      /:/  /       /:/  /       /:/  /    \:\  \       \:\  \   \:\__\    
    \:\__\    /:/  /       /:/  /       /:/  /      \:\__\       \:\__\   \/__/    
     \/__/    \/__/        \/__/        \/__/        \/__/        \/__/         

`

func main() {
	// Define a command-line flag for the theme.
	theme := flag.String("theme", "dark", "Set the UI theme ('dark' or 'light')")
	flag.Parse()

	// Set the background color profile based on the theme flag.
	// This prevents lipgloss from querying the terminal, fixing issues on macOS.
	switch *theme {
	case "dark":
		lipgloss.SetHasDarkBackground(true)
	case "light":
		lipgloss.SetHasDarkBackground(false)
	default:
		fmt.Fprintf(os.Stderr, "Error: Invalid theme value '%s'. Please use 'dark' or 'light'.\n", *theme)
		os.Exit(1)
	}

	// After parsing our own flags, check if there are any remaining positional arguments.
	// If there are, we run in CLI mode. Otherwise, we start the interactive TUI.
	if flag.NArg() > 0 {
		// Reconstruct the arguments for the CLI command processor.
		// os.Args[0] is the program name, and flag.Args() contains the rest.
		cliArgs := append([]string{os.Args[0]}, flag.Args()...)
		if err := cli.ProcessCLICommand(cliArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// No positional arguments provided - start interactive mode.
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
