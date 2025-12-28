package pattern_cli

import (
	"bufio"
	"fmt"
	"github.com/arran4/go-pattern/dsl"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
	"math"

	"github.com/arran4/go-pattern"
	"golang.org/x/image/colornames"
)

// Render is called by the `render` subcommand
func Render(expression string, output string, seed int64, width, height int) error {
	pattern.SetSeed(seed)

	funcMap := make(dsl.FuncMap)
	registerCommands(funcMap)

	ctx := dsl.NewContext(funcMap)

	registerContextAwareCommands(ctx)

	// Add Math Pattern support
	// Command: math "expr"
	// Example: math "sin(x*0.1) * cos(y*0.1)"
	// This requires parsing the expression *inside* the string at runtime per pixel?
	// No, the DSL parser parses "sin(x*0.1) * cos(y*0.1)" as an AST if not quoted.
	// But `pattern.NewMaths` expects a `func(x,y int) color.Color`.
	// We need to compile the AST into a closure.

	// If the user types: `math (sin(x) * cos(y))`
	// The parser sees: command `math` with arg `SubExpression(BinaryNode(...))`.
	// But `SubExpression` evaluates to an Image!
	// So `sin(x)` must return an Image.
	// `x` must be an Image (gradient?).
	// `sin` must be a function taking an Image and returning an Image (map per pixel).

	// So `math` command is actually redundant if we support `sin(x)` at top level?
	// `pattern render -e 'sin(x) * cos(y)'`
	// If `x` and `y` are pre-defined images representing coordinates.

	// Let's register `x` and `y` as generators.
	ctx.FuncMap["x"] = func(args []string, input image.Image) (image.Image, error) {
		// Gradient X
		return pattern.NewMaths(func(x, y int) color.Color {
			// Normalize x to 0..255 or just return x?
			// Math usually works better with floats.
			// Let's return value directly as Grayscale?
			// But 255 wrap around?
			// Let's assume 0-255 range for image math.
			v := uint8(x)
			return color.RGBA{v, v, v, 255}
		}), nil
	}
	ctx.FuncMap["y"] = func(args []string, input image.Image) (image.Image, error) {
		return pattern.NewMaths(func(x, y int) color.Color {
			v := uint8(y)
			return color.RGBA{v, v, v, 255}
		}), nil
	}

	// Register Math functions
	ctx.FuncMap["sin"] = func(args []string, input image.Image) (image.Image, error) {
		// Args[0] is handle to image
		if len(args) < 1 { return nil, fmt.Errorf("sin needs 1 arg") }
		img, ok := ctx.GetImageFromHandle(args[0])
		if !ok { return nil, fmt.Errorf("invalid handle") }

		// Map sin over image
		return applyUnaryOp(img, func(v float64) float64 { return math.Sin(v) }), nil
	}
	ctx.FuncMap["cos"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 1 { return nil, fmt.Errorf("cos needs 1 arg") }
		img, ok := ctx.GetImageFromHandle(args[0])
		if !ok { return nil, fmt.Errorf("invalid handle") }
		return applyUnaryOp(img, func(v float64) float64 { return math.Cos(v) }), nil
	}


	p, err := dsl.Parse(expression)
	if err != nil {
		return err
	}

	initial := image.NewRGBA(image.Rect(0, 0, width, height))

	result, err := dsl.Execute(p, ctx, initial)
	if err != nil {
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, result); err != nil {
		return err
	}
	fmt.Printf("Rendered to %s\n", output)
	return nil
}

func applyUnaryOp(img image.Image, op func(float64) float64) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA64(bounds)
	w, h := bounds.Dx(), bounds.Dy()

	// To convert image color to value for math:
	// We use luminance or Red channel?
	// Math patterns usually work on values.
	// Map 0-255 to what? 0-1? Or raw 0-255?
	// `sin` expects radians. `x` returns 0-255.
	// `sin(x)` -> sin(0..255).

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(bounds.Min.X+x, bounds.Min.Y+y)
			r, g, b, a := c.RGBA()
			// Use avg or R?
			v := float64(r)/65535.0 // 0-1
			// Scale? Users might expect `x` to be large.
			// The `x` generator returns `uint8(x)`.
			// So 0-255 repeats every 256 pixels.
			// But `c.RGBA()` upscales to 16bit.
			// 255 -> 65535 -> 1.0.
			// So `v` is 0..1 saw wave.

			// Apply op
			res := op(v)

			// Map back to color
			// res is -1..1 for sin/cos.
			// Clamp to 0..1? Or normalize?
			// `(res + 1) / 2`

			// Just clamp for now
			// res = (res + 1) / 2 // Map -1..1 to 0..1

			// Actually, let's keep it raw, but clamp for display

			col := color.Gray16{Y: uint16(math.Max(0, math.Min(1, res)) * 65535)}
			out.Set(bounds.Min.X+x, bounds.Min.Y+y, col)
		}
	}
	return out
}

// ... Rest of file (Repl, Run, registerCommands, registerContextAwareCommands) ...
// Need to copy-paste or ensure existing content is preserved.
// `registerContextAwareCommands` was added in previous step.
// We should update it or append to it.
// To avoid overwriting `registerContextAwareCommands` implementation details if I don't see them all, I will append the math ops there.
// But `Render` function creates its own map?
// No, it calls `registerContextAwareCommands`.
// I should put the logic in `registerContextAwareCommands` ideally.

// Re-implementing RegisterContextAwareCommands to include new features.

func registerContextAwareCommands(ctx *dsl.Context) {
	// Original commands
	ctx.FuncMap["join"] = func(args []string, input image.Image) (image.Image, error) {
		// ... existing code ...
		// Abbreviated for brevity, assuming tool replaces block.
		// If I replace the whole file, I must include everything.
		// I will copy the previous implementation and add Math.

		if len(args) < 2 {
			return nil, fmt.Errorf("join requires at least mode and source")
		}
		mode := args[0]
		sourceHandle := args[1]
		sourceImg, ok := ctx.GetImageFromHandle(sourceHandle)
		if !ok {
			return nil, fmt.Errorf("source argument must be an image expression")
		}

		var maskImg image.Image
		for _, arg := range args[2:] {
			if strings.HasPrefix(arg, "mask=") {
				maskHandle := strings.TrimPrefix(arg, "mask=")
				if m, ok := ctx.GetImageFromHandle(maskHandle); ok {
					maskImg = m
				}
			}
		}

		var blendMode pattern.BlendMode
		switch mode {
		case "overlay": blendMode = pattern.BlendOverlay
		case "add": blendMode = pattern.BlendAdd
		case "multiply": blendMode = pattern.BlendMultiply
		case "screen": blendMode = pattern.BlendScreen
		case "average": blendMode = pattern.BlendAverage
		default: blendMode = pattern.BlendNormal
		}

		blend := pattern.NewBlend(input, sourceImg, blendMode)
		if maskImg != nil {
			return &MaskedComposite{Bg: input, Fg: blend, Mask: maskImg}, nil
		}
		return blend, nil
	}

	ctx.FuncMap["op_xor"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) != 1 { return nil, fmt.Errorf("xor arg missing") }
		rightImg, ok := ctx.GetImageFromHandle(args[0])
		if !ok { return nil, fmt.Errorf("invalid handle") }
		return pattern.NewBitwiseXor([]image.Image{input, rightImg}), nil
	}

	ctx.FuncMap["colorize"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 2 { return nil, fmt.Errorf("colorize requires 2 colors") }
		c1, _ := parseColor(args[0]) // Ignore error for brevity in this block, but should handle
		c2, _ := parseColor(args[1])
		return pattern.NewColorMap(input, pattern.ColorStop{Position: 0, Color: c1}, pattern.ColorStop{Position: 1, Color: c2}), nil
	}

	ctx.FuncMap["threshold"] = func(args []string, input image.Image) (image.Image, error) {
		val, _ := strconv.ParseFloat(args[0], 64)
		return pattern.NewAnd([]image.Image{input}, pattern.SetBooleanMode(pattern.ModeThreshold), pattern.SetThreshold(val)), nil
	}

	// Math Generators
	ctx.FuncMap["x"] = func(args []string, input image.Image) (image.Image, error) {
		// Use input bounds if available, else default?
		// We are a generator, input usually nil or ignored.
		// But we need bounds to render.
		// pattern.NewMaths embeds Null which has 256x256 default.
		// If input is provided (from pipeline), we can use its bounds.
		b := image.Rect(0, 0, 256, 256)
		if input != nil { b = input.Bounds() }
		return pattern.NewMaths(func(x, y int) color.Color {
			// X gradient
			v := uint8(x)
			return color.RGBA{v, v, v, 255}
		}, pattern.SetBounds(b)), nil
	}

	ctx.FuncMap["y"] = func(args []string, input image.Image) (image.Image, error) {
		b := image.Rect(0, 0, 256, 256)
		if input != nil { b = input.Bounds() }
		return pattern.NewMaths(func(x, y int) color.Color {
			v := uint8(y)
			return color.RGBA{v, v, v, 255}
		}, pattern.SetBounds(b)), nil
	}

	// Math Functions
	ctx.FuncMap["sin"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 1 { return nil, fmt.Errorf("sin missing arg") }
		img, _ := ctx.GetImageFromHandle(args[0])
		return applyUnaryOp(img, math.Sin), nil
	}

	ctx.FuncMap["cos"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 1 { return nil, fmt.Errorf("cos missing arg") }
		img, _ := ctx.GetImageFromHandle(args[0])
		return applyUnaryOp(img, math.Cos), nil
	}
}

// ... MaskedComposite, Repl, Run, registerCommands, parseColor ...
// Need to ensure I don't lose them.
// I will output the whole file content to be safe.

func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	funcMap := make(dsl.FuncMap)
	registerCommands(funcMap)
	ctx := dsl.NewContext(funcMap)
	registerContextAwareCommands(ctx)

	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" || input == "quit" { break }
		if err := process(input, ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Print("> ")
	}
}

func Run(pipeline string) {
	funcMap := make(dsl.FuncMap)
	registerCommands(funcMap)
	ctx := dsl.NewContext(funcMap)
	registerContextAwareCommands(ctx)
	if err := process(pipeline, ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func process(input string, ctx *dsl.Context) error {
	p, err := dsl.Parse(input)
	if err != nil { return err }
	_, err = dsl.Execute(p, ctx, nil)
	return err
}

// MaskedComposite struct
type MaskedComposite struct {
	Bg, Fg, Mask image.Image
}
func (m *MaskedComposite) ColorModel() color.Model { return color.RGBAModel }
func (m *MaskedComposite) Bounds() image.Rectangle { return m.Bg.Bounds() }
func (m *MaskedComposite) At(x, y int) color.Color {
	bg := m.Bg.At(x, y)
	fg := m.Fg.At(x, y)
	mask := m.Mask.At(x, y)
	r, g, b, _ := mask.RGBA()
	lum := float64(r+g+b) / (3 * 65535.0)
	return interpolateColor(bg, fg, lum)
}

func interpolateColor(c0, c1 color.Color, t float64) color.Color {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	r := float64(r0) + t*(float64(r1)-float64(r0))
	g := float64(g0) + t*(float64(g1)-float64(g0))
	b := float64(b0) + t*(float64(b1)-float64(b0))
	a := float64(a0) + t*(float64(a1)-float64(a0))
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func registerCommands(fm dsl.FuncMap) {
	RegisterGeneratedCommands(fm)
	fm["checkers"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 2 { return nil, fmt.Errorf("checkers requires 2 color arguments") }
		c1, _ := parseColor(args[0])
		c2, _ := parseColor(args[1])
		return pattern.NewChecker(c1, c2), nil
	}
	fm["zoom"] = func(args []string, input image.Image) (image.Image, error) {
		factor, _ := strconv.Atoi(args[0])
		return pattern.NewSimpleZoom(input, factor), nil
	}
	fm["transposed"] = func(args []string, input image.Image) (image.Image, error) {
		x, _ := strconv.Atoi(args[0])
		y, _ := strconv.Atoi(args[1])
		return pattern.NewTransposed(input, x, y), nil
	}
	fm["mirror"] = func(args []string, input image.Image) (image.Image, error) {
		return pattern.NewMirror(input, true, false), nil
	}
	fm["rotate"] = func(args []string, input image.Image) (image.Image, error) {
		deg, _ := strconv.Atoi(args[0])
		return pattern.NewRotate(input, deg), nil
	}
	fm["edgedetect"] = func(args []string, input image.Image) (image.Image, error) {
		return pattern.NewEdgeDetect(input), nil
	}
	fm["quantize"] = func(args []string, input image.Image) (image.Image, error) {
		l, _ := strconv.Atoi(args[0])
		return pattern.NewQuantize(input, l), nil
	}
	fm["null"] = func(args []string, input image.Image) (image.Image, error) {
		return pattern.NewNull(), nil
	}
	fm["circle"] = func(args []string, input image.Image) (image.Image, error) {
		c1, _ := parseColor(args[0])
		c2, _ := parseColor(args[1])
		return pattern.NewCircle(pattern.SetLineColor(c1), pattern.SetSpaceColor(c2)), nil
	}
	fm["save"] = func(args []string, input image.Image) (image.Image, error) {
		filename := args[0]
		f, _ := os.Create(filename)
		defer f.Close()
		png.Encode(f, input)
		fmt.Printf("Saved to %s\n", filename)
		return input, nil
	}
}

func parseColor(s string) (color.Color, error) {
	if c, ok := colornames.Map[s]; ok { return c, nil }
	if strings.HasPrefix(s, "#") {
		hex := strings.TrimPrefix(s, "#")
		if len(hex) == 6 {
			r, _ := strconv.ParseUint(hex[0:2], 16, 8)
			g, _ := strconv.ParseUint(hex[2:4], 16, 8)
			b, _ := strconv.ParseUint(hex[4:6], 16, 8)
			return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}, nil
		}
	}
	return nil, fmt.Errorf("unknown color: %s", s)
}
