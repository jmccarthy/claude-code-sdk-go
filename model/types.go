package model

// PermissionMode defines how tools may execute.
type PermissionMode string

const (
	PermissionDefault     PermissionMode = "default"
	PermissionAcceptEdits PermissionMode = "acceptEdits"
	PermissionBypass      PermissionMode = "bypassPermissions"
)

// ContentBlock represents a chunk of assistant output.
type ContentBlock interface{ isContentBlock() }

// TextBlock is plain text returned from Claude.
type TextBlock struct {
	Text string
}

func (TextBlock) isContentBlock() {}

// ToolUseBlock requests invocation of a tool with arguments.
type ToolUseBlock struct {
	ID    string
	Name  string
	Input map[string]any
}

func (ToolUseBlock) isContentBlock() {}

// ToolResultBlock returns the output of a previously requested tool.
type ToolResultBlock struct {
	ToolUseID string
	Content   any
	IsError   bool
}

func (ToolResultBlock) isContentBlock() {}

// UserMessage originated from the caller.
type UserMessage struct {
	Content string
}

func (UserMessage) isMessage() {}

// AssistantMessage is streaming content from Claude.
type AssistantMessage struct {
	Content []ContentBlock
}

func (AssistantMessage) isMessage() {}

// SystemMessage conveys internal metadata.
type SystemMessage struct {
	Subtype string
	Data    map[string]any
}

func (SystemMessage) isMessage() {}

// ResultMessage summarizes the final session result.
type ResultMessage struct {
	Subtype       string
	CostUSD       float64
	DurationMS    int
	DurationAPIMS int
	IsError       bool
	NumTurns      int
	SessionID     string
	TotalCostUSD  float64
	Usage         map[string]any
	Result        string
}

func (ResultMessage) isMessage() {}

// Message represents any message returned from the CLI.
type Message interface{ isMessage() }

// Options configures a query to Claude.
type Options struct {
	AllowedTools             []string
	MaxThinkingTokens        int
	SystemPrompt             string
	AppendSystemPrompt       string
	MCPTools                 []string
	MCPServers               map[string]MCPServerConfig
	CLIPath                  string
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
