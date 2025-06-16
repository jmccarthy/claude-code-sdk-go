package internal

import (
	"reflect"
	"testing"

	"github.com/anthropics/claude-code-sdk-go/model"
)

func TestBuildCommandBasic(t *testing.T) {
	opts := &model.Options{SystemPrompt: "hi"}
	cmd := buildCommand("/usr/bin/claude", "hello", opts)
	expect := []string{"/usr/bin/claude", "--output-format", "stream-json", "--verbose", "--system-prompt", "hi", "--max-thinking-tokens", "8000", "--print", "hello"}
	if !reflect.DeepEqual(cmd, expect) {
		t.Fatalf("unexpected cmd: %v", cmd)
	}
}

func TestBuildCommandWithExtras(t *testing.T) {
	opts := &model.Options{MaxThinkingTokens: 9000, MCPTools: []string{"foo", "bar"}}
	cmd := buildCommand("/usr/bin/claude", "hello", opts)
	expect := []string{"/usr/bin/claude", "--output-format", "stream-json", "--verbose", "--max-thinking-tokens", "9000", "--mcpTools", "foo,bar", "--print", "hello"}
	if !reflect.DeepEqual(cmd, expect) {
		t.Fatalf("unexpected cmd: %v", cmd)
	}
}
