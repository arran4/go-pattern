# go-pattern

`go-pattern` is a collection of pre-created `image.Image` implementations. It provides a variety of ready-to-use patterns that implement the standard Go image interface.

These patterns are designed to be:
- **Ready to use**: Instantly available as standard `image.Image` objects.
- **Composable**: Easily combined (e.g., zooming, transposing) to form complex visual structures.
- **Standard**: Fully compatible with any Go library that accepts `image.Image`.

## Patterns


### Null Pattern

Returns a transparent color for all pixels.

![Null Pattern](null.png)

```go
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
```


### Checker Pattern

Alternates between two colors in a checkerboard fashion.

![Checker Pattern](checker.png)

```go
i := NewChecker(color.Black, color.White)
	f, err := os.Create(CheckerOutputFilename)
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
```


### SimpleZoom Pattern

Zooms in on an underlying image.

![SimpleZoom Pattern](simplezoom.png)

```go
i := NewSimpleZoom(NewChecker(color.Black, color.White), 2)
	f, err := os.Create(SimpleZoomOutputFilename)
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
```


### Transposed Pattern

Transposes the X and Y coordinates of an underlying image.

![Transposed Pattern](transposed.png)

```go
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
```


