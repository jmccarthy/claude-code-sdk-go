package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
)

// Transport defines methods for communicating with the Claude CLI.
type Transport interface {
	Connect(ctx context.Context) error
	Disconnect() error
	ReceiveMessages(ctx context.Context) (<-chan map[string]any, error)
}

// SubprocessCLITransport implements Transport using os/exec.
type SubprocessCLITransport struct {
	Prompt  string
	Options []string

	cmd    *exec.Cmd
	stdout io.ReadCloser
}

// ErrCLINotFound is returned when the Claude CLI cannot be located.
var ErrCLINotFound = errors.New("claude CLI not found")

func (t *SubprocessCLITransport) Connect(ctx context.Context) error {
	t.cmd = exec.CommandContext(ctx, "claude", append([]string{"query", t.Prompt}, t.Options...)...)
	stdout, err := t.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	t.stdout = stdout
	if err := t.cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (t *SubprocessCLITransport) Disconnect() error {
	if t.cmd != nil && t.cmd.Process != nil {
		return t.cmd.Process.Kill()
	}
	return nil
}

func (t *SubprocessCLITransport) ReceiveMessages(ctx context.Context) (<-chan map[string]any, error) {
	ch := make(chan map[string]any)
	scanner := bufio.NewScanner(t.stdout)
	go func() {
		defer close(ch)
		for scanner.Scan() {
			line := scanner.Bytes()
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				// ignore errors for now
				continue
			}
			ch <- m
		}
	}()
	return ch, nil
}
