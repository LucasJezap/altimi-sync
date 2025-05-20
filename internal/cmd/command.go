package cmd

import (
	"altimi-sync/internal/flags"
	"altimi-sync/internal/lib"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var exitFunc = os.Exit

type Command struct {
	sourceDirectory   string
	targetDirectory   string
	helpFlag          flags.CommandFlagInterface
	deleteMissingFlag flags.CommandFlagInterface
}

func NewCommand() *Command {
	command := &Command{}

	command.helpFlag = &flags.HelpFlag{}
	command.helpFlag.Init()
	command.deleteMissingFlag = &flags.DeleteMissingFlag{}
	command.deleteMissingFlag.Init()

	flag.Parse()

	if command.helpFlag.IsSet() {
		flags.PrintHelpMessage()
		exitFunc(0)
	}

	args := flag.Args()
	if len(args) != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "Error: source and destination must be specified")
		exitFunc(1)
	}

	command.sourceDirectory = args[0]
	command.targetDirectory = args[1]

	return command
}

func (c *Command) Run() {
	//sync files from A to B
	_ = filepath.Walk(c.sourceDirectory, func(sourcePath string, info os.FileInfo, err error) error {
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ACCESS DENIED] %s: %v\033[0m\n", sourcePath, err)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ERROR] Accessing %s: %v\033[0m\n", sourcePath, err)
			}
			return nil
		}

		// files in subdirectories will be processed later in Walk
		if info.IsDir() {
			return nil
		}

		// get absolute path for destination file
		relativePath, _ := filepath.Rel(c.sourceDirectory, sourcePath)
		targetPath := filepath.Join(c.targetDirectory, relativePath)

		// first we check whether file sizes and modification times match (quick)
		// if they do - we proceed to check file checksums to be 100% sure (slow)
		// in that way we might omit second check
		needsSync := true
		isOverwrite := true
		destinationInfo, err := os.Stat(targetPath)
		fmt.Println(err)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ACCESS DENIED] %s: %v\033[0m\n", sourcePath, err)
				return nil
			}
			isOverwrite = false
		} else {
			if info.Size() == destinationInfo.Size() && info.ModTime() == destinationInfo.ModTime() {
				checksumMatches, err := lib.CompareChecksum(sourcePath, targetPath)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ERROR] Comparing checksums for %s and %s: %v\033[0m\n", sourcePath, targetPath, err)
					return nil
				}
				if checksumMatches {
					needsSync = false
					isOverwrite = false
				}
			}
		}

		if needsSync {
			_ = os.MkdirAll(filepath.Dir(targetPath), 0755)
			if err := lib.CopyFile(sourcePath, targetPath, info.ModTime(), isOverwrite); err != nil {
				if errors.Is(err, os.ErrPermission) {
					_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ACCESS DENIED] %s: %v\033[0m\n", targetPath, err)
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "\033[31m❌ [ERROR] Failed to copy %s to %s: %v\033[0m\n", sourcePath, targetPath, err)
				}
				return nil
			}
		}

		return nil
	})

	//sync files in B that do not exist in A
	if c.deleteMissingFlag.IsSet() {
		_ = filepath.Walk(c.targetDirectory, func(targetPath string, info os.FileInfo, err error) error {
			if err != nil {
				if errors.Is(err, os.ErrPermission) {
					_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ACCESS DENIED] %s: %v\033[0m\n", targetPath, err)
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ERROR] Accessing %s: %v\033[0m\n", targetPath, err)
				}
				return nil
			}

			// files in subdirectories will be processed later in Walk
			if info.IsDir() {
				return nil
			}

			// get absolute path for source file
			relativePath, _ := filepath.Rel(c.targetDirectory, targetPath)
			sourcePath := filepath.Join(c.sourceDirectory, relativePath)

			// check if source file exists, if not then delete
			_, err = os.Stat(sourcePath)
			if err != nil {
				if errors.Is(err, os.ErrPermission) {
					_, _ = fmt.Fprintf(os.Stderr, "\033[31m🚨 [ACCESS DENIED] %s: %v\033[0m\n", sourcePath, err)
					return nil
				} else if os.IsExist(err) {
					return nil
				}

				if err := lib.RemoveFile(targetPath); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "\033[31m❌ [ERROR] Failed to remove %s: %v\033[0m\n", targetPath, err)
				}
			}

			return nil
		})
	}

	_, _ = fmt.Fprintln(os.Stdout, "✅ Sync complete")
}
