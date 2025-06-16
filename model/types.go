package model

// PermissionMode defines how tools may execute.
type PermissionMode string

const (
	PermissionDefault     PermissionMode = "default"
	PermissionAcceptEdits PermissionMode = "acceptEdits"
	PermissionBypass      PermissionMode = "bypassPermissions"
)

// Options configures a query to Claude.
type Options struct {
	AllowedTools             []string
	MaxThinkingTokens        int
	SystemPrompt             string
	AppendSystemPrompt       string
	MCPTools                 []string
	MCPServers               map[string]MCPServerConfig
	PermissionMode           PermissionMode
	ContinueConversation     bool
	Resume                   string
	MaxTurns                 int
	DisallowedTools          []string
	Model                    string
	PermissionPromptToolName string
	Cwd                      string
}

// MCPServerConfig represents configuration for an MCP server.
type MCPServerConfig struct {
	URL    string
	APIKey string
}

// Message is a generic response from Claude.
type Message interface{}
