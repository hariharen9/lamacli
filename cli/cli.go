package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hariharen9/lamacli/fileops"
	"github.com/hariharen9/lamacli/llm"
)

// Command represents a CLI command type.
type Command string

const (
	CommandAsk     Command = "ask"
	CommandSuggest Command = "suggest"
	CommandExplain Command = "explain"
	CommandConfig  Command = "config"
	CommandVersion Command = "version"
	CommandHelp    Command = "help"
)

// CommandOptions holds options for CLI commands
type CommandOptions struct {
	Model        string
	Context      string
	Include      string
	SystemPrompt string
}

// Version information
const (
	Version = "0.4.0"
)

// ProcessCLICommand processes CLI commands with flags and arguments
func ProcessCLICommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("no command specified. Use 'lamacli help' for usage information")
	}

	cmdStr := args[1]
	command := parseCommand(cmdStr)

	switch command {
	case CommandHelp:
		printHelp()
		return nil
	case CommandVersion:
		printVersion()
		return nil
	case CommandConfig:
		return handleConfigCommand(args[2:])
	case CommandAsk, CommandSuggest, CommandExplain:
		return handleLLMCommand(command, args[2:])
	default:
		return fmt.Errorf("unknown command '%s'. Use 'lamacli help' for usage information", cmdStr)
	}
}

// parseCommand converts string to Command type, handling aliases
func parseCommand(cmdStr string) Command {
	switch strings.ToLower(cmdStr) {
	case "ask", "a":
		return CommandAsk
	case "suggest", "s":
		return CommandSuggest
	case "explain", "e":
		return CommandExplain
	case "config", "c":
		return CommandConfig
	case "version", "v":
		return CommandVersion
	case "help", "h":
		return CommandHelp
	default:
		return Command(cmdStr)
	}
}

// handleLLMCommand processes ask, suggest, and explain commands
func handleLLMCommand(command Command, args []string) error {
	// Initialize Ollama client
	llmClient, err := llm.NewOllamaClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Ollama client: %w", err)
	}

	// Parse flags and options
	options, prompt, err := parseCommandFlags(args)
	if err != nil {
		return err
	}

	if prompt == "" {
		return fmt.Errorf("prompt is required for %s command", command)
	}

	// Get default model if not specified
	model := options.Model
	if model == "" {
		model = getDefaultModel(llmClient)
	}

	// Build context if specified
	contextContent := ""
	if options.Context != "" {
		contextContent, err = buildContext(options.Context, options.Include)
		if err != nil {
			return fmt.Errorf("failed to build context: %w", err)
		}
	}

	// Prepare system prompt based on command
	systemPrompt := buildSystemPrompt(command, options.SystemPrompt)

	// Combine prompt with context
	finalPrompt := prompt
	if contextContent != "" {
		finalPrompt = fmt.Sprintf("%s\n\nContext:\n%s", prompt, contextContent)
	}

// Show simple loading indicator that works on all terminals
	fmt.Print("Thinking")
	
	// Create a done channel to coordinate the loading indicator
	loadingDone := make(chan bool)
	
	// Start simple dot-based loading animation
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-loadingDone:
				return
			case <-ticker.C:
				fmt.Print(".")
			}
		}
	}()
	
	// Generate response
	response, err := llmClient.GenerateResponse(model, finalPrompt, systemPrompt)
	
	// Stop the loading animation
	loadingDone <- true
	// Print a newline to finish the loading line
	fmt.Println()

	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	// Print response with appropriate formatting
	printFormattedResponse(command, response, model)
	return nil
}

// parseCommandFlags parses command line flags and returns options and prompt
func parseCommandFlags(args []string) (*CommandOptions, string, error) {
	flags := flag.NewFlagSet("lamacli", flag.ContinueOnError)
	flags.Usage = func() {} // Suppress default usage

	options := &CommandOptions{}
	flags.StringVar(&options.Model, "model", "", "Override default model")
	flags.StringVar(&options.Context, "context", "", "Include directory context")
	flags.StringVar(&options.Include, "include", "", "File pattern to include in context")
	flags.StringVar(&options.SystemPrompt, "system", "", "Custom system prompt")

	err := flags.Parse(args)
	if err != nil {
		return nil, "", err
	}

	// Get the remaining arguments as the prompt
	remainingArgs := flags.Args()
	prompt := strings.Join(remainingArgs, " ")

	return options, prompt, nil
}

// getDefaultModel gets the first available model as default
func getDefaultModel(llmClient *llm.OllamaClient) string {
	models, err := llmClient.ListModels()
	if err != nil || len(models) == 0 {
		return "llama3.2:3b" // fallback
	}
	return models[0]
}

// buildContext builds context from directory and file patterns
func buildContext(contextPath, includePattern string) (string, error) {
	var contextBuilder strings.Builder

	// Handle relative paths
	if contextPath == "." {
		contextPath, _ = os.Getwd()
	}

	// Walk through directory
	err := filepath.Walk(contextPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Filter by include pattern if specified
		if includePattern != "" {
			matched, _ := filepath.Match(includePattern, info.Name())
			if !matched {
				return nil
			}
		}

		// Read file content
		content, err := fileops.ReadFile(path)
		if err != nil {
			return nil // Skip files that can't be read
		}

		// Add to context
		relPath, _ := filepath.Rel(contextPath, path)
		contextBuilder.WriteString(fmt.Sprintf("\n--- File: %s ---\n%s\n", relPath, string(content)))

		// Limit context size to prevent overwhelming the LLM
		if contextBuilder.Len() > 10000 {
			return filepath.SkipDir
		}

		return nil
	})

	return contextBuilder.String(), err
}

// buildSystemPrompt creates appropriate system prompt based on command
func buildSystemPrompt(command Command, customPrompt string) string {
	if customPrompt != "" {
		return customPrompt
	}

	switch command {
	case CommandSuggest:
		return "You are a helpful command-line assistant. When asked to suggest a command, provide concise, practical command-line solutions. Include brief explanations when helpful."
	case CommandExplain:
		return "You are a helpful technical assistant. When asked to explain a command, provide clear, detailed explanations of what the command does, its options, and usage examples."
	case CommandAsk:
		return "You are a helpful assistant. Provide clear, accurate, and helpful responses to questions."
	default:
		return "You are a helpful assistant."
	}
}

// printFormattedResponse prints the response with appropriate formatting
func printFormattedResponse(command Command, response, model string) {
	switch command {
	case CommandSuggest:
		fmt.Printf("\n🔮 Suggested Command (using %s):\n%s\n\n", model, response)
	case CommandExplain:
		fmt.Printf("\n📖 Command Explanation (using %s):\n%s\n\n", model, response)
	case CommandAsk:
		fmt.Printf("\n💭 Response (using %s):\n%s\n\n", model, response)
	default:
		fmt.Printf("\n%s\n\n", response)
	}
}

// handleConfigCommand handles configuration management
func handleConfigCommand(args []string) error {
	// For now, just print available models
	llmClient, err := llm.NewOllamaClient()
	if err != nil {
		return fmt.Errorf("failed to initialize Ollama client: %w", err)
	}

	models, err := llmClient.ListModels()
	if err != nil {
		return fmt.Errorf("failed to list models: %w", err)
	}

	fmt.Println("\n🤖 Available Models:")
	for i, model := range models {
		if i == 0 {
			fmt.Printf("  • %s (default)\n", model)
		} else {
			fmt.Printf("  • %s\n", model)
		}
	}
	fmt.Println()

	return nil
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("LamaCLI version %s\n", Version)
}

// printHelp prints comprehensive help information
func printHelp() {
	fmt.Println(`
🦙 LamaCLI - Your Terminal AI Assistant

USAGE:
  lamacli [command] [options] "<prompt>"
  lamacli                              # Start interactive mode

COMMANDS:
  ask, a      Ask a question
  suggest, s  Get command suggestions  
  explain, e  Explain a command
  config, c   Show configuration (available models)
  version, v  Show version information
  help, h     Show this help message

OPTIONS:
  --model     Override default model (e.g., --model=llama3.2:1b)
  --context   Include directory context (e.g., --context=.)
  --include   File pattern for context (e.g., --include=*.md)
  --system    Custom system prompt

EXAMPLES:
  lamacli ask "How do I list files in Linux?"
  lamacli a --model=qwen2.5-coder:1.5b "Explain async/await"
  
  lamacli suggest "find large files"
  lamacli s --model=llama3.2:1b "git workflow for teams"
  
  lamacli explain "find . -name '*.go' -exec grep -l 'func main' {} \;"
  lamacli e --model=qwen2.5-coder "docker compose up -d"
  
  lamacli ask --context=. --include="*.md" "Summarize this project"
  lamacli config
  lamacli version

NOTE: Run 'lamacli' without arguments to start the interactive mode.
`)
}
