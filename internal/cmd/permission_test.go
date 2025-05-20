package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCopySourcePermissionDenied(t *testing.T) {
	resetFlags()

	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	sourceFile := filepath.Join(sourceDir, "altimi.txt")
	err := os.WriteFile(sourceFile, []byte("hello altimi"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	err = os.Chmod(sourceFile, 000)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}

	os.Args = []string{"altimi-sync", sourceDir, targetDir}
	cmd := NewCommand()

	output := captureOutput(func() {
		cmd.Run()
	})

	if !strings.Contains(output, "🚨 [ACCESS DENIED]") {
		t.Errorf("Expected permission denied error, got: %s", output)
	}
}

func TestRunCopyTargetPermissionDenied(t *testing.T) {
	resetFlags()

	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	sourceFile := filepath.Join(sourceDir, "altimi.txt")
	err := os.WriteFile(sourceFile, []byte("hello altimi"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	targetFile := filepath.Join(targetDir, "altimi.txt")
	err = os.WriteFile(targetFile, []byte("hello altimi"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	err = os.Chmod(targetFile, 0444)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}

	os.Args = []string{"altimi-sync", sourceDir, targetDir}
	cmd := NewCommand()

	output := captureOutput(func() {
		cmd.Run()
	})

	if !strings.Contains(output, "🚨 [ACCESS DENIED]") {
		t.Errorf("Expected permission denied error, got: %s", output)
	}
}

func TestRunDeletePermissionDenied(t *testing.T) {
	resetFlags()

	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	sourceFile := filepath.Join(sourceDir, "altimi.txt")
	err := os.WriteFile(sourceFile, []byte("hello altimi"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	targetFile := filepath.Join(targetDir, "altimi.txt")
	err = os.WriteFile(targetFile, []byte("hello altimi"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	err = os.Chmod(targetFile, 0444)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}

	os.Args = []string{"altimi-sync", "-d", sourceDir, targetDir}
	cmd := NewCommand()

	output := captureOutput(func() {
		cmd.Run()
	})

	if !strings.Contains(output, "🚨 [ACCESS DENIED]") {
		t.Errorf("Expected permission denied error, got: %s", output)
	}
}

func TestRunDirectoriesPermissionDenied(t *testing.T) {
	resetFlags()

	sourceDir := "/root/protected-dirA"
	targetDir := "/root/protected-dirB"

	os.Args = []string{"altimi-sync", sourceDir, targetDir}
	cmd := NewCommand()

	output := captureOutput(func() {
		cmd.Run()
	})

	if !strings.Contains(output, "🚨 [ERROR]") {
		t.Errorf("Expected permission denied error when creating target directory, got: %s", output)
	}
}
