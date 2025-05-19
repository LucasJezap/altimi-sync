package main

import (
	"altimi-sync/internal/cmd"
)

func main() {
	command := cmd.NewCommand()
	command.Run()
}
