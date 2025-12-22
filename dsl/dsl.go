package dsl

import (
	"strings"
)

type Command struct {
	Name string
	Args []string
}

type Pipeline []Command

func Parse(input string) (Pipeline, error) {
	var p Pipeline
	parts := strings.Split(input, "|")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) == 0 {
			continue
		}
		cmd := Command{
			Name: fields[0],
			Args: fields[1:],
		}
		p = append(p, cmd)
	}
	return p, nil
}

func (p Pipeline) String() string {
	var parts []string
	for _, cmd := range p {
		s := cmd.Name
		if len(cmd.Args) > 0 {
			s += " " + strings.Join(cmd.Args, " ")
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, " | ")
}
