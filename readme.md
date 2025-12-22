# go-pattern

`go-pattern` is a collection of pre-created `image.Image` implementations. It provides a variety of ready-to-use patterns that implement the standard Go image interface.

These patterns are designed to be:
- **Ready to use**: Instantly available as standard `image.Image` objects.
- **Composable**: Easily combined (e.g., zooming, transposing) to form complex visual structures.
- **Standard**: Fully compatible with any Go library that accepts `image.Image`.

## Patterns


### Null Pattern

Undefined RGBA colour.

![Null Pattern](null.png)

```go
i := pattern.NewNull()
	f, err := os.Create("null.png")
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
i := pattern.NewChecker(color.Black, color.White)
	f, err := os.Create("checker.png")
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


### Simple Zoom Pattern

Zooms in on an underlying image.

![Simple Zoom Pattern](simplezoom.png)

```go
i := pattern.NewSimpleZoom(pattern.NewChecker(color.Black, color.White), 2)
	f, err := os.Create("simplezoom.png")
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
i := pattern.NewTransposed(pattern.NewDemoNull(), 10, 10)
	f, err := os.Create("transposed.png")
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



## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.
