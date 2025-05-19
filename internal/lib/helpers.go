package lib

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func CopyFile(source, destination string, modificationTime time.Time, isOverwrite bool) error {
	if isOverwrite {
		_, _ = fmt.Fprintf(os.Stdout, "♻️ Overwriting file %s...\n", destination)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "📄 Copying file %s...\n", destination)
	}

	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer closeFile(in)

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer closeFile(out)

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	syncFile(out) //save file before settings access/modification time

	if err := os.Chtimes(destination, modificationTime, modificationTime); err != nil {
		return err
	}

	if isOverwrite {
		_, _ = fmt.Fprintf(os.Stdout, "✅ File overwritten: %s\n", destination)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "✅ File copied: %s\n", destination)
	}

	return nil
}

func RemoveFile(path string) error {
	_, _ = fmt.Fprintf(os.Stdout, "🗑️ Deleting file %s...\n", path)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete %s: %w", path, err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "✅ File deleted: %s\n", path)

	return nil
}

func CompareChecksum(path1, path2 string) (bool, error) {
	hash1, err := fileChecksum(path1)
	if err != nil {
		return false, err
	}

	hash2, err := fileChecksum(path2)
	if err != nil {
		return false, err
	}

	return hash1 == hash2, nil
}

func fileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer closeFile(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func syncFile(f *os.File) {
	if err := f.Sync(); err != nil {
		log.Printf("failed to sync file: %v", err)
	}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Printf("failed to close file: %v", err)
	}
}
