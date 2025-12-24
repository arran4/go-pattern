package pattern

import (
	"image"
)

func NewDemoGridRows(ops ...func(any)) image.Image {
	gopher := NewScale(NewGopher(), ScaleFactor(0.25))

	args := []any{
		Row(Cell(gopher), Cell(gopher)),
		Row(Cell(gopher), Cell(gopher)),
	}
	for _, op := range ops {
		args = append(args, op)
	}

	return NewGrid(args...)
}
