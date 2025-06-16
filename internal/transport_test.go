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
