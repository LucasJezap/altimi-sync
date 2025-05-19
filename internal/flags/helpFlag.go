package flags

import (
	"flag"
	"fmt"
	"os"
)

type HelpFlag struct {
	set bool
}

func (f *HelpFlag) Init() {
	flag.BoolVar(&f.set, "h", false, "")
	flag.BoolVar(&f.set, "help", false, "")
}

func (f *HelpFlag) IsSet() bool {
	return f.set
}

func PrintHelpMessage() {
	_, _ = fmt.Fprintln(os.Stdout, `Usage:
	  altimi-sync [options] <source> <destination>
	
	Description:
	  Synchronizes files from a source directory to a destination directory.
	  Only files that differ are copied. Optionally, files in the destination
	  that do not exist in the source can be deleted.
	
	Positional Arguments:
	  source         Path to the source directory
	  destination    Path to the destination directory
	
	Options:
	  -d, --delete-missing   Delete files in the destination that are not present in the source
	  -h, --help             Show this help message and exit
	
	Examples:
	  altimi-sync ./data ./backup
	  altimi-sync -d ./source ./target`)
}
