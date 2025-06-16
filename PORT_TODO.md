# Porting Plan: Python Claude Code SDK to Go

This document outlines the steps to reimplement the Python SDK found in this
repository as an idiomatic Go package. The goal is feature parity with the
current Python version while embracing Go conventions for structure, error
handling and concurrency.

## 1. Repository Setup

1. **Initialize a module** ✅
   - `go mod init github.com/anthropics/claude-code-sdk-go`
   - Use Go 1.20+ for generics and context features.
2. **Directory layout** ✅
   - `cmd/` (optional): example binaries or demos. ✅
   - `claudecode/`: root package containing the public API. ✅
   - `internal/`: subpackages for client and transport implementations similar to `src/claude_code_sdk/_internal` in Python. ✅
   - `examples/`: rewrite Python examples using Go. ✅
   - `test/`: unit tests mirroring behaviour of the current `tests/` folder.

## 2. API Surface ✅

Replicate the high level `query` function which streams `Message` values.
In Go this could be:

```go
func Query(ctx context.Context, prompt string, opts *Options) (<-chan Message, error)
```

`Message` is a sum type represented by separate structs. The function returns a
read-only channel to allow consumers to range over the stream.

## 3. Types ✅

Translate all dataclasses in `src/claude_code_sdk/types.py` to Go structs.
Use string constants for enums such as `PermissionMode`.

```go
// PermissionMode defines how tools may execute.
type PermissionMode string

const (
    PermissionDefault        PermissionMode = "default"
    PermissionAcceptEdits    PermissionMode = "acceptEdits"
    PermissionBypass         PermissionMode = "bypassPermissions"
)
```

Define message and content block structures:

```go
type TextBlock struct { Text string }

type ToolUseBlock struct {
    ID   string
    Name string
    Input map[string]any
}

type ToolResultBlock struct {
    ToolUseID string
    Content   any // string or []map[string]any
    IsError   bool
}

// Union of blocks
// interface{ TextBlock(), ToolUseBlock(), ... } can be represented with an
// interface and type assertions.
```

Similarly translate `UserMessage`, `AssistantMessage`, `SystemMessage`, and
`ResultMessage` structs.
Options become:

```go
type Options struct {
    AllowedTools       []string
    MaxThinkingTokens  int
    SystemPrompt       string
    AppendSystemPrompt string
    MCPTools           []string
    MCPServers         map[string]MCPServerConfig
    PermissionMode     PermissionMode
    ContinueConversation bool
    Resume             string
    MaxTurns           int
    DisallowedTools    []string
    Model              string
    PermissionPromptToolName string
    Cwd                string
}
```

## 4. Errors ✅

Create custom error types analogous to `_errors.py`. Each should implement the
`error` interface. Example:

```go
type CLIConnectionError struct{ Msg string }
func (e *CLIConnectionError) Error() string { return e.Msg }
```

Include errors for `CLINotFoundError`, `ProcessError` (with fields for exit code
and stderr), and `CLIJSONDecodeError`.

## 5. Transport Layer ✅

Recreate the `Transport` interface and the subprocess implementation using
`os/exec`.

- `Connect()` starts the CLI process with the proper command-line arguments.
- `Disconnect()` terminates the process respecting context cancellation.
- `SendRequest()` is unused but keep for interface completeness.
- `ReceiveMessages()` returns a channel of `map[string]any` decoded from the
  CLI JSON stream. Use `bufio.Scanner` to read lines from `stdout`.
- Search for the `claude` binary using locations from `_find_cli` in the Python
  code. When not found, return `CLINotFoundError` with helpful instructions.

## 6. Internal Client ✅

Implement a client similar to `_internal/client.py`:

1. Construct a `SubprocessCLITransport` with the prompt and options.
2. Connect, then range over the received JSON messages.
3. Map raw data to strongly typed `Message` values (mirroring `_parse_message`).
4. Ensure all resources are cleaned up with `defer transport.Disconnect()`.

**Context & Resource Management:**
- Propagate context through all operations for proper cancellation.
- Use `context.WithCancel()` to ensure subprocess cleanup on context cancellation.
- Implement proper channel closing patterns to prevent goroutine leaks.
- Consider using `sync.WaitGroup` or `errgroup.Group` for coordinating cleanup.

## 7. Public API ✅

Expose a package-level `Query` function which sets `CLAUDE_CODE_ENTRYPOINT`
environment variable to `sdk-go` and delegates to the internal client.

Consumers should be able to use:

```go
ctx := context.Background()
ch, err := claudecode.Query(ctx, "What is 2 + 2?", nil)
for msg := range ch {
    // handle messages
}
```

## 8. Testing ✅

Rewrite Python tests using Go's `testing` package. Focus on:

- Command construction for the transport.
- Parsing of CLI JSON into message structs.
- Error cases (CLI not found, process errors, JSON decode failures).
- **TODO:** Go `Query` currently drops these errors. Implement propagation so callers
  receive `CLINotFoundError`, `ProcessError` and `CLIJSONDecodeError` like the
  Python SDK, and add tests for these cases.
- High level `Query` behaviour with mocked transport (use interfaces and test
  doubles).
- Context cancellation behavior and timeout handling.
- Resource cleanup (ensure processes are terminated and channels closed).

Add Go-specific testing:
- Race condition detection using `go test -race`.
- Table-driven tests for comprehensive option validation.

Aim for parity with the current `tests/` to maintain confidence during the port.

## 9. Examples & Documentation

- Convert `examples/quick_start.py` to `examples/quick_start.go` showing basic usage and options. ✅
- Update `README.md` to include Go installation instructions and usage snippets. ✅

## 10. CI/CD & Development Workflow

Set up automated testing and quality checks:

- **GitHub Actions**: Configure workflows for testing on multiple Go versions and platforms. ✅
- **Code Quality**: Integrate `golangci-lint` for comprehensive linting. ✅
- **Coverage**: Set up code coverage reporting and enforcement. ✅
- **Dependabot**: Enable automatic dependency updates. ✅
- **Release Automation**: Consider using `goreleaser` for automated releases. ✅

## 11. Future Enhancements

- Consider exposing a synchronous API for simple use cases in addition to the
  streaming channel form.
- Provide context-aware cancellation so callers can abort long-running queries.
- Explore packaging the CLI binary with the Go SDK for easier installation.
- **Observability**: Add structured logging with configurable levels.

---
This plan should serve as a starting point for reimplementing the Claude Code
SDK in Go while following idiomatic patterns for package layout, error handling
and concurrency.
