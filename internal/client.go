package internal

import (
	"context"

	"github.com/anthropics/claude-code-sdk-go/model"
)

// Client orchestrates communication with the Claude CLI.
type Client struct {
	Transport Transport
}

func (c *Client) Query(ctx context.Context, prompt string, opts *model.Options) (<-chan model.Message, error) {
	if err := c.Transport.Connect(ctx); err != nil {
		return nil, err
	}

	rawCh, err := c.Transport.ReceiveMessages(ctx)
	if err != nil {
		return nil, err
	}

	out := make(chan model.Message)
	go func() {
		defer close(out)
		for range rawCh {
			// TODO: parse into Message structs
			var m model.Message
			out <- m
		}
	}()
	return out, nil
}
