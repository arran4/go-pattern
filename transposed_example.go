package pattern

import (
	"image"
	"image/png"
	"os"
)

var TransposedOutputFilename = "transposed.png"

func ExampleNewTransposed() {
	i := NewTransposed(NewDemoNull(), 10, 10)
	f, err := os.Create(TransposedOutputFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()
	if err = png.Encode(f, i); err != nil {
		panic(err)
	}
}

func BootstrapTransposed(ops ...func(any)) image.Image {
	return NewTransposed(NewSimpleZoom(NewDemoChecker(ops...), 10, ops...), 5, 5, ops...)
}

func BootstrapTransposedReferences() (map[string]func(ops ...func(any)) image.Image, []string) {
	return map[string]func(ops ...func(any)) image.Image{
		"Original": func(ops ...func(any)) image.Image {
			return NewSimpleZoom(NewDemoChecker(ops...), 10, ops...)
		},
	}, []string{"Original"}
}

const TransposedOrder = 3
