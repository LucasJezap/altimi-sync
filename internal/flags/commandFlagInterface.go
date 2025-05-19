package flags

//this might be helpful if the flags have complicated use-case/definition,
//then we may add extra functions for processing etc.,
//in this command it's a little overkill

type CommandFlagInterface interface {
	Init()
	IsSet() bool
}
