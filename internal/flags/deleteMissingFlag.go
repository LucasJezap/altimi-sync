package flags

import (
	"flag"
)

type DeleteMissingFlag struct {
	set bool
}

func (f *DeleteMissingFlag) Init() {
	flag.BoolVar(&f.set, "d", false, "")
	flag.BoolVar(&f.set, "delete-missing", false, "")
}

func (f *DeleteMissingFlag) IsSet() bool {
	return f.set
}
