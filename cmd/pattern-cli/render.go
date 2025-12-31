package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arran4/go-pattern/pkg/pattern-cli"
)

var _ Cmd = (*renderCmd)(nil)

type renderCmd struct {
	*RootCmd
	Flags *flag.FlagSet

	expression string
	output     string
	seed       int64
	size       string
	channels   string

	SubCommands map[string]Cmd
}

func (c *renderCmd) Usage() {
	// Simple usage for now
	fmt.Fprintf(os.Stderr, "Usage: %s render [options]\n", os.Args[0])
	c.Flags.PrintDefaults()
}

func (c *renderCmd) Execute(args []string) error {
	err := c.Flags.Parse(args)
	if err != nil {
		return NewUserError(err, fmt.Sprintf("flag parse error %s", err.Error()))
	}

	if c.expression == "" {
		return fmt.Errorf("expression (-e) is required")
	}

	// Parsing size "512x512"
	width, height := 256, 256
	if c.size != "" {
		parts := strings.Split(c.size, "x")
		if len(parts) == 2 {
			if _, err := fmt.Sscanf(parts[0], "%d", &width); err != nil {
				return fmt.Errorf("invalid width in size: %v", err)
			}
			if _, err := fmt.Sscanf(parts[1], "%d", &height); err != nil {
				return fmt.Errorf("invalid height in size: %v", err)
			}
		}
	}

	// Call into pkg/pattern-cli
	return pattern_cli.Render(c.expression, c.output, c.seed, width, height)
}

func (c *RootCmd) NewrenderCmd() *renderCmd {
	set := flag.NewFlagSet("render", flag.ContinueOnError)
	v := &renderCmd{
		RootCmd:     c,
		Flags:       set,
		SubCommands: make(map[string]Cmd),
	}

	set.StringVar(&v.expression, "e", "", "Pattern expression to render")
	set.StringVar(&v.output, "o", "out.png", "Output filename")
	set.Int64Var(&v.seed, "seed", 0, "Random seed")
	set.StringVar(&v.size, "size", "256x256", "Output size (WxH)")
	set.StringVar(&v.channels, "channels", "", "Channels (not fully implemented)")

	set.Usage = v.Usage

	return v
}
