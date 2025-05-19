package flags

import (
	"flag"
	"os"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestHelpFlagShort(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-h"}
	help := &HelpFlag{}
	help.Init()
	flag.Parse()

	if !help.IsSet() {
		t.Error("Expected help flag to be set with -h")
	}
}

func TestHelpFlagLong(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--help"}
	help := &HelpFlag{}
	help.Init()
	flag.Parse()

	if !help.IsSet() {
		t.Error("Expected help flag to be set with --help")
	}
}

func TestDeleteMissingFlagShort(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "-d"}
	del := &DeleteMissingFlag{}
	del.Init()
	flag.Parse()

	if !del.IsSet() {
		t.Error("Expected delete-missing flag to be set with -d")
	}
}

// Test DeleteMissingFlag with --delete-missing
func TestDeleteMissingFlagLong(t *testing.T) {
	resetFlags()
	os.Args = []string{"cmd", "--delete-missing"}
	del := &DeleteMissingFlag{}
	del.Init()
	flag.Parse()

	if !del.IsSet() {
		t.Error("Expected delete-missing flag to be set with --delete-missing")
	}
}
