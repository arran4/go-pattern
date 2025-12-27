package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/arran4/go-pattern/pkg/pattern-cli"
)

var _ Cmd = (*runCmd)(nil)

type runCmd struct {
	*RootCmd
	Flags *flag.FlagSet

	pipeline string

	SubCommands map[string]Cmd
}

func (c *runCmd) Usage() {
	err := executeUsage(os.Stderr, "run_usage.txt", c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating usage: %s\n", err)
	}
}

func (c *runCmd) Execute(args []string) error {
	if len(args) > 0 {
		if cmd, ok := c.SubCommands[args[0]]; ok {
			return cmd.Execute(args[1:])
		}
	}
	err := c.Flags.Parse(args)
	if err != nil {
		return NewUserError(err, fmt.Sprintf("flag parse error %s", err.Error()))
	}
	pattern_cli.Run(c.pipeline)
	return nil
}

func (c *RootCmd) NewrunCmd() *runCmd {
	set := flag.NewFlagSet("run", flag.ContinueOnError)
	v := &runCmd{
		RootCmd:     c,
		Flags:       set,
		SubCommands: make(map[string]Cmd),
	}

	set.StringVar(&v.pipeline, "pipeline", "", "TODO: Add usage text")

	set.Usage = v.Usage

	return v
}
