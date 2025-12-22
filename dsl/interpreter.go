package dsl

import (
	"fmt"
	"image"
)

type CommandFunc func(args []string, input image.Image) (image.Image, error)
type FuncMap map[string]CommandFunc

func (p Pipeline) Execute(fm FuncMap, initial image.Image) (image.Image, error) {
	var img image.Image = initial
	var err error
	for _, cmd := range p {
		fn, ok := fm[cmd.Name]
		if !ok {
			return nil, fmt.Errorf("unknown command: %s", cmd.Name)
		}
		img, err = fn(cmd.Args, img)
		if err != nil {
			return nil, fmt.Errorf("command %s failed: %w", cmd.Name, err)
		}
	}
	return img, nil
}
