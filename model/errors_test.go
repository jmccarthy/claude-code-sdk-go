package model

import "testing"

func TestProcessError(t *testing.T) {
	err := &ProcessError{Msg: "fail", ExitCode: 1, Stderr: "bad"}
	if err.Error() == "" {
		t.Fatal("missing error message")
	}
	if err.ExitCode != 1 || err.Stderr != "bad" {
		t.Fatal("fields not set")
	}
}

func TestCLINotFoundError(t *testing.T) {
	err := &CLINotFoundError{Msg: "missing"}
	if err.Error() != "missing" {
		t.Fatalf("unexpected message: %s", err.Error())
	}
}
