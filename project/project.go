package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProjectType represents the type of a project.
type ProjectType string

const (
	GoProject      ProjectType = "Go"
	NodeProject    ProjectType = "Node.js"
	PythonProject  ProjectType = "Python"
	UnknownProject ProjectType = "Unknown"
)

// Project represents a project with its type and root path.
type Project struct {
	Type    ProjectType
	Root    string
	Context string // Context for LLM summarization
	Summary string // LLM generated summary
}

// ScanForProject scans a directory for project files to determine the project type
// and collects context for LLM summarization.
func ScanForProject(rootPath string) (*Project, error) {
	project := &Project{Root: rootPath}
	var contextBuilder strings.Builder

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip dotfiles and common ignored directories
		if strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == "dist" || info.Name() == "build" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		switch info.Name() {
		case "go.mod":
			project.Type = GoProject
			content, _ := os.ReadFile(path)
			contextBuilder.WriteString(fmt.Sprintf("\n--- %s ---\n%s", info.Name(), string(content)))
		case "package.json":
			project.Type = NodeProject
			content, _ := os.ReadFile(path)
			contextBuilder.WriteString(fmt.Sprintf("\n--- %s ---\n%s", info.Name(), string(content)))
		case "requirements.txt":
			project.Type = PythonProject
			content, _ := os.ReadFile(path)
			contextBuilder.WriteString(fmt.Sprintf("\n--- %s ---\n%s", info.Name(), string(content)))
		}

		// Add file names to context (up to a certain limit to avoid overwhelming LLM)
		if !info.IsDir() && contextBuilder.Len() < 5000 { // Arbitrary limit for now
			contextBuilder.WriteString(fmt.Sprintf("\nFile: %s", info.Name()))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if project.Type == "" {
		project.Type = UnknownProject
	}

	project.Context = contextBuilder.String()

	return project, nil
}
