package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestSuccessfulExecuteHelp(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--help")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error when executing root command with --help, received %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "ReqCorder") {
		t.Error("Expected output to contain 'ReqCorder'")
	}
}

func TestSuccessfulVersionCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error when executing version command, received %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "ReqCorder version") {
		t.Error("Expected version output to contain 'ReqCorder version'")
	}
}

func TestSuccessfulHelpCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "help")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error when executing help command, received %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Usage of ReqCorder") {
		t.Error("Expected help output to contain usage information")
	}
}

func TestFailedInvalidCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "invalid-command")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err == nil {
		t.Error("Expected error when executing invalid command")
	}
}
