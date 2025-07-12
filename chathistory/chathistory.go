package chathistory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ChatSession represents a saved chat session
type ChatSession struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Model     string    `json:"model"`
	History   []string  `json:"history"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatHistoryManager manages chat history persistence
type ChatHistoryManager struct {
	historyDir string
}

// NewChatHistoryManager creates a new chat history manager
func NewChatHistoryManager() (*ChatHistoryManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	historyDir := filepath.Join(homeDir, ".lamacli", "chat_history")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create chat history directory: %w", err)
	}
	
	return &ChatHistoryManager{
		historyDir: historyDir,
	}, nil
}

// SaveSession saves a chat session to disk
func (chm *ChatHistoryManager) SaveSession(session *ChatSession) error {
	if session.ID == "" {
		session.ID = generateSessionID()
	}
	
	// Generate title from first user message if not set
	if session.Title == "" {
		session.Title = chm.generateSessionTitle(session.History)
	}
	
	session.UpdatedAt = time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = session.UpdatedAt
	}
	
	filename := fmt.Sprintf("%s.json", session.ID)
	filepath := filepath.Join(chm.historyDir, filename)
	
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}
	
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}
	
	return nil
}

// LoadSession loads a chat session from disk
func (chm *ChatHistoryManager) LoadSession(sessionID string) (*ChatSession, error) {
	filename := fmt.Sprintf("%s.json", sessionID)
	filepath := filepath.Join(chm.historyDir, filename)
	
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}
	
	var session ChatSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}
	
	return &session, nil
}

// ListSessions returns all available chat sessions, sorted by update time (newest first)
func (chm *ChatHistoryManager) ListSessions() ([]*ChatSession, error) {
	files, err := os.ReadDir(chm.historyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read history directory: %w", err)
	}
	
	var sessions []*ChatSession
	
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			sessionID := strings.TrimSuffix(file.Name(), ".json")
			session, err := chm.LoadSession(sessionID)
			if err != nil {
				// Skip corrupted files
				continue
			}
			sessions = append(sessions, session)
		}
	}
	
	// Sort by update time (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})
	
	return sessions, nil
}

// DeleteSession deletes a chat session from disk
func (chm *ChatHistoryManager) DeleteSession(sessionID string) error {
	filename := fmt.Sprintf("%s.json", sessionID)
	filepath := filepath.Join(chm.historyDir, filename)
	
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete session file: %w", err)
	}
	
	return nil
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().Unix())
}

// generateSessionTitle generates a title from the chat history
func (chm *ChatHistoryManager) generateSessionTitle(history []string) string {
	// Find the first non-empty user message
	for i := 0; i < len(history); i += 2 {
		if history[i] != "" {
			// Take first 50 characters as title
			title := strings.TrimSpace(history[i])
			if len(title) > 50 {
				title = title[:47] + "..."
			}
			return title
		}
	}
	
	// Fallback title
	return fmt.Sprintf("Chat Session %s", time.Now().Format("Jan 2, 2006"))
}

// GetSessionSummary returns a brief summary of the session
func (session *ChatSession) GetSessionSummary() string {
	messageCount := len(session.History) / 2
	timeAgo := time.Since(session.UpdatedAt)
	
	var timeStr string
	if timeAgo < time.Hour {
		timeStr = fmt.Sprintf("%d minutes ago", int(timeAgo.Minutes()))
	} else if timeAgo < 24*time.Hour {
		timeStr = fmt.Sprintf("%d hours ago", int(timeAgo.Hours()))
	} else {
		timeStr = fmt.Sprintf("%d days ago", int(timeAgo.Hours()/24))
	}
	
	return fmt.Sprintf("%s (%d messages, %s)", session.Title, messageCount, timeStr)
}
