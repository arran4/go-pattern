package pattern

import (
	"image/png"
	"os"
)

var NullOutputFilename = "null.png"

func ExampleNewNull() {
	i := NewNull()
	f, err := os.Create(NullOutputFilename)
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

const NullOrder = 0
