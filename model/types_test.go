package model

import "testing"

func TestUserMessage(t *testing.T) {
	m := UserMessage{Content: "hi"}
	if m.Content != "hi" {
		t.Fatal("unexpected content")
	}
}

func TestAssistantMessage(t *testing.T) {
	msg := AssistantMessage{Content: []ContentBlock{TextBlock{Text: "hello"}}}
	if len(msg.Content) != 1 {
		t.Fatalf("expected 1 block")
	}
}

func TestOptionsDefaults(t *testing.T) {
	opts := Options{}
	if opts.MaxThinkingTokens != 0 {
		t.Fatalf("expected default 0, got %d", opts.MaxThinkingTokens)
	}
	if opts.PermissionMode != "" {
		t.Fatalf("expected empty permission mode")
	}
}
