package lib

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCopyFile(t *testing.T) {
	restore := suppressOutput(t)
	defer restore()

	srcDir := t.TempDir()
	dstDir := t.TempDir()

	srcFile := filepath.Join(srcDir, "file.txt")
	dstFile := filepath.Join(dstDir, "file.txt")

	content := []byte("hello ALTIMI")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	modTime := time.Now()

	err := CopyFile(srcFile, dstFile, modTime, false)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("Copied file contents mismatch: got %s, want %s", data, content)
	}

	info, _ := os.Stat(dstFile)
	if !info.ModTime().Round(time.Second).Equal(modTime.Round(time.Second)) {
		t.Errorf("Modification time mismatch: got %v, want %v", info.ModTime(), modTime)
	}
}

func TestRemoveFile(t *testing.T) {
	restore := suppressOutput(t)
	defer restore()

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "unused-altimi.txt")

	if err := os.WriteFile(file, []byte("byebyebyebye"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := RemoveFile(file); err != nil {
		t.Errorf("Expected file to be deleted, got error: %v", err)
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Errorf("Expected file to not exist, but it does")
	}
}

func TestCompareChecksum(t *testing.T) {
	restore := suppressOutput(t)
	defer restore()

	dir := t.TempDir()
	f1 := filepath.Join(dir, "altimi1.txt")
	f2 := filepath.Join(dir, "altimi2.txt")
	f3 := filepath.Join(dir, "altimi3.txt")

	_ = os.WriteFile(f1, []byte("abc"), 0644)
	_ = os.WriteFile(f2, []byte("abc"), 0644)
	_ = os.WriteFile(f3, []byte("xyz"), 0644)

	match, err := CompareChecksum(f1, f2)
	if err != nil || !match {
		t.Errorf("Expected matching checksums, got %v, err: %v", match, err)
	}

	match, err = CompareChecksum(f1, f3)
	if err != nil || match {
		t.Errorf("Expected different checksums, got %v, err: %v", match, err)
	}

	_, err = CompareChecksum(f1, "not-exist.txt")
	if err == nil {
		t.Error("Expected error comparing to missing file, got nil")
	}
}

func suppressOutput(t *testing.T) (restore func()) {
	t.Helper()

	null, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatalf("failed to open os.DevNull: %v", err)
	}

	stdout := os.Stdout
	stderr := os.Stderr

	os.Stdout = null
	os.Stderr = null

	return func() {
		os.Stdout = stdout
		os.Stderr = stderr
		_ = null.Close()
	}
}
