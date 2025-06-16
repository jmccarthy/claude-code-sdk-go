package model

import "fmt"

// CLIConnectionError is returned when the CLI cannot be started or communicated with.
type CLIConnectionError struct{ Msg string }

func (e *CLIConnectionError) Error() string { return e.Msg }

// CLINotFoundError is returned when the claude CLI binary cannot be located.
type CLINotFoundError struct{ Msg string }

func (e *CLINotFoundError) Error() string { return e.Msg }

// ProcessError is returned when the CLI exits with a non-zero status.
type ProcessError struct {
	Msg      string
	ExitCode int
	Stderr   string
}

func (e *ProcessError) Error() string { return e.Msg }

// CLIJSONDecodeError is returned when a line from the CLI fails to decode as JSON.
type CLIJSONDecodeError struct {
	Line string
	Err  error
}

func (e *CLIJSONDecodeError) Error() string {
	return fmt.Sprintf("failed to decode JSON: %v", e.Err)
}
