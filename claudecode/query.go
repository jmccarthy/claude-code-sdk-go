package claudecode

import (
	"context"
	"os"

	"github.com/anthropics/claude-code-sdk-go/internal"
	"github.com/anthropics/claude-code-sdk-go/model"
)

// Query sends a prompt to Claude Code and returns a stream of Messages.
func Query(ctx context.Context, prompt string, opts *model.Options) (<-chan model.Message, <-chan error, error) {
	os.Setenv("CLAUDE_CODE_ENTRYPOINT", "sdk-go")
	transport := &internal.SubprocessCLITransport{}
	client := &internal.Client{Transport: transport}
	return client.Query(ctx, prompt, opts)
}
