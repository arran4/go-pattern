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

// Repl is a subcommand `pattern-cli repl`
func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	funcMap := make(dsl.FuncMap)
	registerCommands(funcMap)
	ctx := dsl.NewContext(funcMap)
	registerContextAwareCommands(ctx)

	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" || input == "quit" {
			break
		}
		if err := process(input, ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Print("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}
}

// Run is a subcommand `pattern-cli run`
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
	if err != nil {
		return err
	}
	_, err = dsl.Execute(p, ctx, nil)
	return err
}

func registerContextAwareCommands(ctx *dsl.Context) {
	// These commands need access to ctx to resolve handles
	ctx.FuncMap["join"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("join requires at least mode and source")
		}
		mode := args[0]
		// source
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
				} else {
					return nil, fmt.Errorf("invalid mask handle: %s", maskHandle)
				}
			}
		}

		var blendMode pattern.BlendMode
		switch mode {
		case "overlay":
			blendMode = pattern.BlendOverlay
		case "add":
			blendMode = pattern.BlendAdd
		case "multiply":
			blendMode = pattern.BlendMultiply
		case "screen":
			blendMode = pattern.BlendScreen
		case "average":
			blendMode = pattern.BlendAverage
		default:
			blendMode = pattern.BlendNormal
		}

		blend := pattern.NewBlend(input, sourceImg, blendMode)

		if maskImg != nil {
			return &MaskedComposite{
				Bg: input,
				Fg: blend,
				Mask: maskImg,
			}, nil
		}

		return blend, nil
	}

	ctx.FuncMap["op_xor"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("xor operator requires 1 argument (right image)")
		}
		rightImg, ok := ctx.GetImageFromHandle(args[0])
		if !ok {
			return nil, fmt.Errorf("invalid image handle for xor: %s", args[0])
		}
		return pattern.NewBitwiseXor([]image.Image{input, rightImg}), nil
	}

	ctx.FuncMap["colorize"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("colorize requires 2 colors")
		}
		c1, err := parseColor(args[0])
		if err != nil {
			return nil, err
		}
		c2, err := parseColor(args[1])
		if err != nil {
			return nil, err
		}
		return pattern.NewColorMap(input, pattern.ColorStop{Position: 0, Color: c1}, pattern.ColorStop{Position: 1, Color: c2}), nil
	}

	ctx.FuncMap["threshold"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("threshold requires a float value")
		}
		val, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return nil, err
		}
		return pattern.NewAnd([]image.Image{input}, pattern.SetBooleanMode(pattern.ModeThreshold), pattern.SetThreshold(val)), nil
	}
}

// Simple Masked Composite Implementation
type MaskedComposite struct {
	Bg, Fg, Mask image.Image
}

func (m *MaskedComposite) ColorModel() color.Model {
	return color.RGBAModel
}

func (m *MaskedComposite) Bounds() image.Rectangle {
	return m.Bg.Bounds()
}

func (m *MaskedComposite) At(x, y int) color.Color {
	bg := m.Bg.At(x, y)
	fg := m.Fg.At(x, y)
	mask := m.Mask.At(x, y)

	// Get mask alpha/luminance
	r, g, b, _ := mask.RGBA()
	// Average for intensity
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

	return color.RGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: uint16(a),
	}
}

func registerCommands(fm dsl.FuncMap) {
	RegisterGeneratedCommands(fm)
	fm["checkers"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("checkers requires 2 color arguments")
		}
		c1, err := parseColor(args[0])
		if err != nil {
			return nil, err
		}
		c2, err := parseColor(args[1])
		if err != nil {
			return nil, err
		}
		return pattern.NewChecker(c1, c2), nil
	}

	fm["zoom"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("zoom requires an input image")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("zoom requires a factor argument")
		}
		factor, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, fmt.Errorf("invalid zoom factor: %v", err)
		}
		return pattern.NewSimpleZoom(input, factor), nil
	}

	fm["transposed"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("transposed requires an input image")
		}
		x, y := 0, 0
		var err error
		if len(args) >= 1 {
			x, err = strconv.Atoi(args[0])
			if err != nil {
				return nil, fmt.Errorf("invalid x offset: %v", err)
			}
		}
		if len(args) >= 2 {
			y, err = strconv.Atoi(args[1])
			if err != nil {
				return nil, fmt.Errorf("invalid y offset: %v", err)
			}
		}
		return pattern.NewTransposed(input, x, y), nil
	}

	fm["mirror"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("mirror requires an input image")
		}
		horizontal := false
		vertical := false
		if len(args) > 0 {
			switch args[0] {
			case "h":
				horizontal = true
			case "v":
				vertical = true
			case "hv", "vh":
				horizontal = true
				vertical = true
			default:
				return nil, fmt.Errorf("mirror argument must be 'h', 'v', or 'hv'")
			}
		} else {
			horizontal = true
		}
		return pattern.NewMirror(input, horizontal, vertical), nil
	}

	fm["rotate"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("rotate requires an input image")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("rotate requires degrees (90, 180, 270)")
		}
		deg, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, fmt.Errorf("invalid degrees: %v", err)
		}
		return pattern.NewRotate(input, deg), nil
  }
	fm["edgedetect"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("edgedetect requires an input image")
		}
		return pattern.NewEdgeDetect(input), nil
  }
	fm["quantize"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("quantize requires an input image")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("quantize requires a levels argument")
		}
		levels, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, fmt.Errorf("invalid levels: %v", err)
		}
		return pattern.NewQuantize(input, levels), nil
	}

	fm["null"] = func(args []string, input image.Image) (image.Image, error) {
		return pattern.NewNull(), nil
	}

	fm["circle"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("circle requires 2 color arguments (line, space)")
		}
		c1, err := parseColor(args[0])
		if err != nil {
			return nil, err
		}
		c2, err := parseColor(args[1])
		if err != nil {
			return nil, err
		}
		return pattern.NewCircle(pattern.SetLineColor(c1), pattern.SetSpaceColor(c2)), nil
	}

	fm["save"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("save requires an input image")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("save requires a filename argument")
		}
		filename := args[0]
		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if strings.HasSuffix(filename, ".png") {
			if err := png.Encode(f, input); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsupported file format: %s", filename)
		}
		fmt.Printf("Saved to %s\n", filename)
		return input, nil
	}
}

func parseColor(s string) (color.Color, error) {
	if c, ok := colornames.Map[s]; ok {
		return c, nil
	}
	if strings.HasPrefix(s, "#") {
		// Quick hex parser
		hex := strings.TrimPrefix(s, "#")
		if len(hex) == 6 {
			r, _ := strconv.ParseUint(hex[0:2], 16, 8)
			g, _ := strconv.ParseUint(hex[2:4], 16, 8)
			b, _ := strconv.ParseUint(hex[4:6], 16, 8)
			return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
		}
	}
	return nil, fmt.Errorf("unknown color: %s", s)
}
