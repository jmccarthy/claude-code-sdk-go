package internal

import (
	"reflect"
	"testing"

	"github.com/anthropics/claude-code-sdk-go/model"
)

func TestBuildCommandBasic(t *testing.T) {
	opts := &model.Options{SystemPrompt: "hi"}
	cmd := buildCommand("/usr/bin/claude", "hello", opts)
	expect := []string{"/usr/bin/claude", "--output-format", "stream-json", "--verbose", "--system-prompt", "hi", "--print", "hello"}
	if !reflect.DeepEqual(cmd, expect) {
		t.Fatalf("unexpected cmd: %v", cmd)
	}
}

func TestBuildCommandWithExtras(t *testing.T) {
	// Test with other valid options (not the hallucinated max-thinking-tokens/mcpTools)
	opts := &model.Options{MaxTurns: 5, Model: "claude-3-sonnet"}
	cmd := buildCommand("/usr/bin/claude", "hello", opts)
	expect := []string{"/usr/bin/claude", "--output-format", "stream-json", "--verbose", "--max-turns", "5", "--model", "claude-3-sonnet", "--print", "hello"}
	if !reflect.DeepEqual(cmd, expect) {
		t.Fatalf("unexpected cmd: %v", cmd)
	}
}
