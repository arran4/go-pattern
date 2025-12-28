package dsl_test

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/arran4/go-pattern/dsl"
	"github.com/arran4/go-pattern/dsl/mocks"
)

func TestEvaluatePipeline(t *testing.T) {
	fm := make(dsl.FuncMap)

	// Mock Source: returns a red 10x10 image
	fm["source"] = func(args []string, input image.Image) (image.Image, error) {
		return mocks.NewMockImage(10, 10, color.RGBA{255, 0, 0, 255}), nil
	}

	// Mock Filter: turns input green
	fm["greenify"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("input required")
		}
		return mocks.NewMockImage(input.Bounds().Dx(), input.Bounds().Dy(), color.RGBA{0, 255, 0, 255}), nil
	}

	ctx := dsl.NewContext(fm)
	input := `source | greenify`
	p, err := dsl.Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	res, err := dsl.Execute(p, ctx, nil)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	r, g, _, _ := res.At(0, 0).RGBA()
	if g == 0 || r != 0 {
		t.Errorf("Expected green image, got R=%d G=%d", r, g)
	}
}

func TestEvaluateArgs(t *testing.T) {
	fm := make(dsl.FuncMap)

	fm["check_arg"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) != 1 || args[0] != "expected" {
			return nil, fmt.Errorf("arg mismatch: %v", args)
		}
		return mocks.NewMockImage(1, 1, color.White), nil
	}

	fm["check_kv"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) != 1 || args[0] != "key=value" {
			return nil, fmt.Errorf("kv mismatch: %v", args)
		}
		return mocks.NewMockImage(1, 1, color.White), nil
	}

	ctx := dsl.NewContext(fm)

	if _, err := dsl.Execute(mustParse("check_arg expected"), ctx, nil); err != nil {
		t.Errorf("Simple arg failed: %v", err)
	}

	if _, err := dsl.Execute(mustParse("check_kv key=value"), ctx, nil); err != nil {
		t.Errorf("KV arg failed: %v", err)
	}
}

// We need a helper to wire context for tests
func setupContext(funcs map[string]func(*dsl.Context, []string, image.Image) (image.Image, error)) *dsl.Context {
	fm := make(dsl.FuncMap)
	// We create a circular dependency: FuncMap needs to call functions that need Context.
	// But Context needs FuncMap.
	// Solution: Create Context with empty FuncMap, then populate it.
	ctx := dsl.NewContext(fm)

	for name, fn := range funcs {
		// Capture fn and ctx
		f := fn // local copy
		fm[name] = func(args []string, input image.Image) (image.Image, error) {
			return f(ctx, args, input)
		}
	}
	return ctx
}

func TestEvaluateSubExpression(t *testing.T) {
	funcs := map[string]func(*dsl.Context, []string, image.Image) (image.Image, error){
		"source": func(ctx *dsl.Context, args []string, input image.Image) (image.Image, error) {
			return mocks.NewMockImage(5, 5, color.White), nil
		},
		"consumer": func(ctx *dsl.Context, args []string, input image.Image) (image.Image, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("arg missing")
			}
			handle := args[0]
			img, ok := ctx.GetImageFromHandle(handle)
			if !ok {
				return nil, fmt.Errorf("handle not found: %s", handle)
			}
			if img.Bounds().Dx() != 5 {
				return nil, fmt.Errorf("wrong image dimensions")
			}
			return img, nil
		},
	}

	ctx := setupContext(funcs)

	input := `consumer (source)`
	p := mustParse(input)
	_, err := dsl.Execute(p, ctx, nil)
	if err != nil {
		t.Fatalf("SubExpression execution failed: %v", err)
	}
}

func TestEvaluateBinaryOp(t *testing.T) {
	funcs := map[string]func(*dsl.Context, []string, image.Image) (image.Image, error){
		"src_a": func(ctx *dsl.Context, args []string, input image.Image) (image.Image, error) {
			return mocks.NewMockImage(1, 1, color.RGBA{100, 0, 0, 255}), nil
		},
		"src_b": func(ctx *dsl.Context, args []string, input image.Image) (image.Image, error) {
			return mocks.NewMockImage(1, 1, color.RGBA{0, 100, 0, 255}), nil
		},
		"op_xor": func(ctx *dsl.Context, args []string, input image.Image) (image.Image, error) {
			// Logic: input is left, args[0] handle is right
			rightHandle := args[0]
			right, ok := ctx.GetImageFromHandle(rightHandle)
			if !ok {
				return nil, fmt.Errorf("right handle missing")
			}
			// Return a mixed image
			r1, _, _, _ := input.At(0,0).RGBA()
			_, g2, _, _ := right.At(0,0).RGBA()
			return mocks.NewMockImage(1, 1, color.RGBA{uint8(r1), uint8(g2), 0, 255}), nil
		},
	}

	ctx := setupContext(funcs)
	input := `(src_a) ^ (src_b)`
	p := mustParse(input)
	res, err := dsl.Execute(p, ctx, nil)
	if err != nil {
		t.Fatalf("Binary op failed: %v", err)
	}

	r, g, _, _ := res.At(0, 0).RGBA()
	// MockImage returns strict colors?
	// NewMockImage takes color.Color.
	// RGBA{100,0,0,255} -> R=100.
	// RGBA struct uses uint8. Image.At().RGBA() returns uint32 (shifted 8 bits usually).
	// Let's check non-zero.
	if r == 0 || g == 0 {
		t.Errorf("Expected mixed color, got R=%d G=%d", r, g)
	}
}

func mustParse(input string) dsl.Node {
	p, err := dsl.Parse(input)
	if err != nil {
		panic(err)
	}
	return p
}
