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

func TestBuildCommandExtras(t *testing.T) {
	opts := &model.Options{MaxThinkingTokens: 100, MCPTools: []string{"a", "b"}}
	cmd := buildCommand("/bin/claude", "p", opts)
	expect := []string{"/bin/claude", "--output-format", "stream-json", "--verbose", "--max-thinking-tokens", "100", "--mcp-tools", "a,b", "--print", "p"}
	if !reflect.DeepEqual(cmd, expect) {
		t.Fatalf("unexpected cmd: %v", cmd)
	}
}
