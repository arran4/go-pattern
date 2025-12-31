package dsl

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
)

type FuncMap map[string]func([]string, image.Image) (image.Image, error)

type Context struct {
	FuncMap    FuncMap
	TempImages map[string]image.Image
	Counter    int
}

func NewContext(fm FuncMap) *Context {
	return &Context{
		FuncMap:    fm,
		TempImages: make(map[string]image.Image),
		Counter:    0,
	}
}

func (ctx *Context) RegisterTempImage(img image.Image) string {
	id := fmt.Sprintf("@img:%d", ctx.Counter)
	ctx.Counter++
	ctx.TempImages[id] = img
	return id
}

func (ctx *Context) GetImageFromHandle(handle string) (image.Image, bool) {
	if strings.HasPrefix(handle, "@img:") {
		img, ok := ctx.TempImages[handle]
		return img, ok
	}
	return nil, false
}

// Execute executes the node with the given context and initial input.
func (n *CommandNode) Execute(ctx *Context, input image.Image) (image.Image, error) {
	if fn, ok := ctx.FuncMap[n.Name]; ok {
		var args []string
		for _, arg := range n.Args {
			if kv, ok := arg.(*KeyValueNode); ok {
				valStr, err := resolveArg(kv.Value, ctx, input)
				if err != nil {
					return nil, err
				}
				args = append(args, fmt.Sprintf("%s=%s", kv.Key, valStr))
			} else {
				valStr, err := resolveArg(arg, ctx, input)
				if err != nil {
					return nil, err
				}
				args = append(args, valStr)
			}
		}
		return fn(args, input)
	}
	return nil, fmt.Errorf("command not found: %s", n.Name)
}

func resolveArg(arg ArgNode, ctx *Context, input image.Image) (string, error) {
	switch a := arg.(type) {
	case *LiteralNode:
		return a.Value, nil
	case *SubExpressionNode:
		res, err := Execute(a.Node, ctx, input)
		if err != nil {
			return "", err
		}
		return ctx.RegisterTempImage(res), nil
	case *KeyValueNode:
		return "", fmt.Errorf("unexpected key-value in resolveArg")
	}
	return "", fmt.Errorf("unknown arg type")
}

func (n *PipelineNode) Execute(ctx *Context, input image.Image) (image.Image, error) {
	current := input
	var err error
	for _, node := range n.Nodes {
		current, err = Execute(node, ctx, current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

func (n *BinaryNode) Execute(ctx *Context, input image.Image) (image.Image, error) {
	left, err := Execute(n.Left, ctx, input)
	if err != nil {
		return nil, err
	}
	right, err := Execute(n.Right, ctx, input)
	if err != nil {
		return nil, err
	}

	switch n.Operator {
	case "^":
		if fn, ok := ctx.FuncMap["op_xor"]; ok {
             handle := ctx.RegisterTempImage(right)
             return fn([]string{handle}, left)
		}
		return nil, fmt.Errorf("operator ^ (op_xor) not implemented")
	case "+", "-", "*", "/", "%":
		return ExecuteImageMath(left, right, n.Operator)
	default:
		return nil, fmt.Errorf("unknown operator %s", n.Operator)
	}
}

func (n *GroupNode) Execute(ctx *Context, input image.Image) (image.Image, error) {
	return Execute(n.Inner, ctx, input)
}

func Execute(node Node, ctx *Context, input image.Image) (image.Image, error) {
	switch n := node.(type) {
	case *CommandNode:
		return n.Execute(ctx, input)
	case *PipelineNode:
		return n.Execute(ctx, input)
	case *BinaryNode:
		return n.Execute(ctx, input)
	case *GroupNode:
		return n.Execute(ctx, input)
	case *LiteralNode:
		return CreateConstantImage(n.Value, input.Bounds()), nil
	}
	return nil, fmt.Errorf("unknown node type: %T", node)
}

// Helpers for Image Math

func ExecuteImageMath(a, b image.Image, op string) (image.Image, error) {
	bounds := a.Bounds()
	out := image.NewRGBA64(bounds)

	w, h := bounds.Dx(), bounds.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c1 := a.At(bounds.Min.X+x, bounds.Min.Y+y)
			c2 := b.At(bounds.Min.X+x, bounds.Min.Y+y)
			r1, g1, b1, a1 := colorToFloats(c1)
			r2, g2, b2, a2 := colorToFloats(c2)

			var r, g, bVal, aVal float64

			switch op {
			case "+":
				r, g, bVal, aVal = r1+r2, g1+g2, b1+b2, a1+a2
			case "-":
				r, g, bVal, aVal = r1-r2, g1-g2, b1-b2, a1-a2
			case "*":
				r, g, bVal, aVal = r1*r2, g1*g2, b1*b2, a1*a2
			case "/":
				r = divSafe(r1, r2)
				g = divSafe(g1, g2)
				bVal = divSafe(b1, b2)
				aVal = divSafe(a1, a2)
			case "%":
				r = modSafe(r1, r2)
				g = modSafe(g1, g2)
				bVal = modSafe(b1, b2)
				aVal = modSafe(a1, a2)
			}

			out.Set(bounds.Min.X+x, bounds.Min.Y+y, floatsToColor(r, g, bVal, aVal))
		}
	}
	return out, nil
}

func colorToFloats(c color.Color) (float64, float64, float64, float64) {
	r, g, b, a := c.RGBA()
	return float64(r)/65535.0, float64(g)/65535.0, float64(b)/65535.0, float64(a)/65535.0
}

func floatsToColor(r, g, b, a float64) color.Color {
	clamp := func(v float64) uint16 {
		if v < 0 { return 0 }
		if v > 1 { return 65535 }
		return uint16(v * 65535)
	}
	return color.RGBA64{clamp(r), clamp(g), clamp(b), clamp(a)}
}

func divSafe(a, b float64) float64 {
	if b == 0 { return 0 }
	return a / b
}

func modSafe(a, b float64) float64 {
	if b == 0 { return 0 }
	return math.Mod(a, b)
}

func CreateConstantImage(valStr string, bounds image.Rectangle) image.Image {
	var v float64
	if _, err := fmt.Sscanf(valStr, "%f", &v); err != nil {
		// Log or fallback? Assuming 0 on failure or just try?
		// Sscanf might fail for non-numbers (e.g. "red").
		// But Parser identifies LiteralNode.
		// If literal is "red", %f fails.
		// We should probably check if it's a color?
		// But here we are in Eval assuming it's a number for math context if used as operand.
		// If it's used as arg to command, CommandNode handles it.
		// This branch is only hit if LiteralNode appears at top level of expression (e.g. `0.5`).
		v = 0
	}
	c := floatsToColor(v, v, v, 1.0)
	return &image.Uniform{C: c}
}
