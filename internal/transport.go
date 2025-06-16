package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anthropics/claude-code-sdk-go/model"
)

// Transport defines methods for communicating with the Claude CLI.
type Transport interface {
	Connect(ctx context.Context) error
	Disconnect() error
	SendRequest(ctx context.Context, prompt string, opts *model.Options) error
	ReceiveMessages(ctx context.Context) (<-chan map[string]any, error)
}

// SubprocessCLITransport implements Transport using os/exec.
type SubprocessCLITransport struct {
	Prompt  string
	Options *model.Options
	CLIPath string
	Cwd     string

	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr io.ReadCloser
}

// SendRequest stores the prompt and options for the next connection.
func (t *SubprocessCLITransport) SendRequest(ctx context.Context, prompt string, opts *model.Options) error {
	t.Prompt = prompt
	t.Options = opts
	if opts != nil {
		t.Cwd = opts.Cwd
		if opts.CLIPath != "" {
			t.CLIPath = opts.CLIPath
		}
	}
	return nil
}

// Connect starts the CLI process with the provided options.
func (t *SubprocessCLITransport) Connect(ctx context.Context) error {
	if t.CLIPath == "" {
		var err error
		t.CLIPath, err = findCLI()
		if err != nil {
			return err
		}
	}

	cmdArgs := buildCommand(t.CLIPath, t.Prompt, t.Options)
	t.cmd = exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	t.cmd.Env = append(os.Environ(), "CLAUDE_CODE_ENTRYPOINT=sdk-go")
	if t.Cwd != "" {
		t.cmd.Dir = t.Cwd
	}

	stdout, err := t.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := t.cmd.StderrPipe()
	if err != nil {
		return err
	}
	t.stdout = stdout
	t.stderr = stderr

	if err := t.cmd.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return &model.CLINotFoundError{Msg: fmt.Sprintf("Claude Code not found at: %s", t.CLIPath)}
		}
		return &model.CLIConnectionError{Msg: err.Error()}
	}
	return nil
}

// Disconnect terminates the running process.
func (t *SubprocessCLITransport) Disconnect() error {
	if t.cmd != nil && t.cmd.Process != nil {
		_ = t.cmd.Process.Kill()
	}
	return nil
}

// ReceiveMessages streams JSON objects from stdout until the process exits.
func (t *SubprocessCLITransport) ReceiveMessages(ctx context.Context) (<-chan map[string]any, error) {
	if t.cmd == nil || t.stdout == nil {
		return nil, &model.CLIConnectionError{Msg: "not connected"}
	}

	ch := make(chan map[string]any)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(t.stdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				if len(line) > 0 && (line[0] == '{' || line[0] == '[') {
					ch <- map[string]any{"error": &model.CLIJSONDecodeError{Line: string(line), Err: err}}
					continue
				}
				continue
			}
			ch <- m
		}
		if err := t.cmd.Wait(); err != nil {
			exitErr := &exec.ExitError{}
			if errors.As(err, &exitErr) {
				msg := &model.ProcessError{Msg: "CLI process failed", ExitCode: exitErr.ExitCode()}
				if b, _ := io.ReadAll(t.stderr); len(b) > 0 {
					msg.Stderr = string(b)
				}
				ch <- map[string]any{"error": msg}
			}
		}
	}()
	return ch, nil
}

func findCLI() (string, error) {
	if p, err := exec.LookPath("claude"); err == nil {
		return p, nil
	}
	homeDir, _ := os.UserHomeDir()
	locations := []string{
		filepath.Join(homeDir, ".npm-global/bin/claude"),
		"/usr/local/bin/claude",
		filepath.Join(homeDir, ".local/bin/claude"),
		filepath.Join(homeDir, "node_modules/.bin/claude"),
		filepath.Join(homeDir, ".yarn/bin/claude"),
	}
	for _, p := range locations {
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p, nil
		}
	}
	if _, err := exec.LookPath("node"); err != nil {
		return "", &model.CLINotFoundError{Msg: "Claude Code requires Node.js. Install Node.js and then run `npm install -g @anthropic-ai/claude-code`"}
	}
	return "", &model.CLINotFoundError{Msg: "Claude Code not found. Install with `npm install -g @anthropic-ai/claude-code`"}
}

func buildCommand(cliPath, prompt string, opts *model.Options) []string {
	args := []string{cliPath, "--output-format", "stream-json", "--verbose"}

	if opts == nil {
		opts = &model.Options{}
	}
	if opts.SystemPrompt != "" {
		args = append(args, "--system-prompt", opts.SystemPrompt)
	}
	if opts.AppendSystemPrompt != "" {
		args = append(args, "--append-system-prompt", opts.AppendSystemPrompt)
	}
	if len(opts.AllowedTools) > 0 {
		args = append(args, "--allowedTools", join(opts.AllowedTools))
	}
	// Note: max-thinking-tokens and mcpTools CLI parameters don't exist in claude CLI
	// These were hallucinated - removed to match Python SDK exactly
	if opts.MaxTurns > 0 {
		args = append(args, "--max-turns", fmt.Sprint(opts.MaxTurns))
	}
	if len(opts.DisallowedTools) > 0 {
		args = append(args, "--disallowedTools", join(opts.DisallowedTools))
	}
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}
	if opts.PermissionPromptToolName != "" {
		args = append(args, "--permission-prompt-tool", opts.PermissionPromptToolName)
	}
	if opts.PermissionMode != "" {
		args = append(args, "--permission-mode", string(opts.PermissionMode))
	}
	if opts.ContinueConversation {
		args = append(args, "--continue")
	}
	if opts.Resume != "" {
		args = append(args, "--resume", opts.Resume)
	}
	if len(opts.MCPServers) > 0 {
		// simple JSON encoding
		b, _ := json.Marshal(map[string]any{"mcpServers": opts.MCPServers})
		args = append(args, "--mcp-config", string(b))
	}

	args = append(args, "--print", prompt)
	return args
}

func join(list []string) string {
	return strings.Join(list, ",")
}
