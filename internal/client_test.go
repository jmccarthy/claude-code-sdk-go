package internal

import (
	"context"
	"testing"

	"github.com/anthropics/claude-code-sdk-go/model"
)

type stubTransport struct {
	msgs         []map[string]any
	connected    bool
	disconnected bool
	sent         bool
}

func (s *stubTransport) Connect(ctx context.Context) error {
	s.connected = true
	return nil
}

func (s *stubTransport) Disconnect() error {
	s.disconnected = true
	return nil
}

func (s *stubTransport) SendRequest(ctx context.Context, prompt string, opts *model.Options) error {
	s.sent = true
	return nil
}

func (s *stubTransport) ReceiveMessages(ctx context.Context) (<-chan map[string]any, error) {
	ch := make(chan map[string]any)
	go func() {
		for _, m := range s.msgs {
			ch <- m
		}
		close(ch)
	}()
	return ch, nil
}

func TestParseMessageUser(t *testing.T) {
	data := map[string]any{
		"type": "user",
		"message": map[string]any{
			"content": "hi",
		},
	}
	m := parseMessage(data)
	u, ok := m.(model.UserMessage)
	if !ok || u.Content != "hi" {
		t.Fatalf("unexpected message: %#v", m)
	}
}

func TestParseMessageAssistant(t *testing.T) {
	data := map[string]any{
		"type": "assistant",
		"message": map[string]any{
			"content": []any{
				map[string]any{"type": "text", "text": "hello"},
				map[string]any{"type": "tool_use", "id": "1", "name": "Read", "input": map[string]any{"file": "a.txt"}},
			},
		},
	}
	m := parseMessage(data)
	a, ok := m.(model.AssistantMessage)
	if !ok || len(a.Content) != 2 {
		t.Fatalf("unexpected message: %#v", m)
	}
	if tb, ok := a.Content[0].(model.TextBlock); !ok || tb.Text != "hello" {
		t.Fatalf("unexpected first block: %#v", a.Content[0])
	}
	if ub, ok := a.Content[1].(model.ToolUseBlock); !ok || ub.Name != "Read" {
		t.Fatalf("unexpected second block: %#v", a.Content[1])
	}
}

func TestParseMessageResult(t *testing.T) {
	data := map[string]any{
		"type":      "result",
		"subtype":   "success",
		"cost_usd":  0.1,
		"num_turns": 1,
	}
	m := parseMessage(data)
	r, ok := m.(model.ResultMessage)
	if !ok || r.Subtype != "success" || r.CostUSD != 0.1 {
		t.Fatalf("unexpected message: %#v", m)
	}
}

func TestClientQuery(t *testing.T) {
	st := &stubTransport{msgs: []map[string]any{
		{"type": "user", "message": map[string]any{"content": "hi"}},
		{"type": "result", "subtype": "done"},
	}}
	c := &Client{Transport: st}
	ch, errCh, err := c.Query(context.Background(), "hi", nil)
	if err != nil {
		t.Fatal(err)
	}
	var msgs []model.Message
	for m := range ch {
		msgs = append(msgs, m)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if !st.connected || !st.disconnected || !st.sent {
		t.Fatalf("transport lifecycle not called")
	}
	if errVal := <-errCh; errVal != nil {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestClientQueryErrorPropagation(t *testing.T) {
	perr := &model.ProcessError{Msg: "fail", ExitCode: 1}
	st := &stubTransport{msgs: []map[string]any{
		{"error": perr},
	}}
	c := &Client{Transport: st}
	ch, errCh, err := c.Query(context.Background(), "hi", nil)
	if err != nil {
		t.Fatal(err)
	}
	for range ch {
	}
	if e := <-errCh; e != perr {
		t.Fatalf("expected error propagation")
	}
}
