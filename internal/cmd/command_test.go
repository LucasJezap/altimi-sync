package cmd

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = stdout

	out, _ := io.ReadAll(r)
	return string(out)
}

func TestNewCommandMissingArgs(t *testing.T) {
	resetFlags()
	os.Args = []string{"altimi-sync"}

	called := false
	code := -1

	originalExit := exitFunc
	exitFunc = func(c int) {
		called = true
		code = c
		panic("exit")
	}
	defer func() {
		exitFunc = originalExit
		if r := recover(); r == nil {
			t.Errorf("Expected panic from exitFunc")
		}
		if !called {
			t.Errorf("Expected exitFunc to be called")
		}
		if code != 1 {
			t.Errorf("Expected exit code 1, got %d", code)
		}
	}()

	NewCommand()
}

func TestNewCommandHelpFlag(t *testing.T) {
	resetFlags()
	os.Args = []string{"altimi-sync", "-h"}

	called := false
	code := -1

	originalExit := exitFunc
	exitFunc = func(c int) {
		called = true
		code = c
		panic("exit")
	}
	defer func() {
		exitFunc = originalExit
		if r := recover(); r == nil {
			t.Errorf("Expected panic from exitFunc")
		}
		if !called {
			t.Errorf("Expected exitFunc to be called")
		}
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
	}()

	NewCommand()
}

func TestNewCommandValidArgs(t *testing.T) {
	resetFlags()
	os.Args = []string{"altimi-sync", "src", "dst"}

	c := NewCommand()

	if c == nil {
		t.Fatal("Expected command to be created")
	}
	if c.sourceDirectory != "src" || c.targetDirectory != "dst" {
		t.Errorf("Incorrect source or destination parsed")
	}
}

func TestCommandRunSyncFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	srcFile := filepath.Join(src, "altimi.txt")
	err := os.WriteFile(srcFile, []byte("hello altimi"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	resetFlags()
	os.Args = []string{"altimi-sync", src, dst}
	c := NewCommand()
	c.Run()

	dstFile := filepath.Join(dst, "altimi.txt")
	_, err = os.Stat(dstFile)
	if err != nil {
		t.Fatalf("Expected file to be copied, got error: %v", err)
	}
}

func TestRunDeleteMissing(t *testing.T) {
	resetFlags()

	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	targetFile := filepath.Join(targetDir, "altimi.txt")
	_ = os.WriteFile(targetFile, []byte("unused"), 0644)

	os.Args = []string{"altimi-sync", "-d", sourceDir, targetDir}
	cmd := NewCommand()

	output := captureOutput(func() {
		cmd.Run()
	})

	if strings.Contains(output, "File deleted") == false {
		t.Errorf("Expected deletion log, got: %s", output)
	}

	if _, err := os.Stat(targetFile); !os.IsNotExist(err) {
		t.Errorf("Expected file to be deleted")
	}
}
