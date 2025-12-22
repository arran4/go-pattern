package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/arran4/go-pattern/pattern_cli"
)

var _ Cmd = (*replCmd)(nil)

type replCmd struct {
	*RootCmd
	Flags *flag.FlagSet

	SubCommands map[string]Cmd
}

func (c *replCmd) Usage() {
	err := executeUsage(os.Stderr, "repl_usage.txt", c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating usage: %s\n", err)
	}
}

func (c *replCmd) Execute(args []string) error {
	if len(args) > 0 {
		if cmd, ok := c.SubCommands[args[0]]; ok {
			return cmd.Execute(args[1:])
		}
	}
	err := c.Flags.Parse(args)
	if err != nil {
		return NewUserError(err, fmt.Sprintf("flag parse error %s", err.Error()))
	}
	pattern_cli.Repl()
	return nil
}

func (c *RootCmd) NewreplCmd() *replCmd {
	set := flag.NewFlagSet("repl", flag.ContinueOnError)
	v := &replCmd{
		RootCmd:     c,
		Flags:       set,
		SubCommands: make(map[string]Cmd),
	}

	set.Usage = v.Usage

	return v
}
