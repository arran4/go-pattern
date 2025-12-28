package dsl

import (
	"fmt"
	"image"
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
	fn, ok := ctx.FuncMap[n.Name]
	if !ok {
		return nil, fmt.Errorf("command not found: %s", n.Name)
	}

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
	}
	return nil, fmt.Errorf("unknown node type")
}

// Helper wrapper for legacy calls (if any) or simplified usage
// But dsl.Execute signature changed.
// Need to update callers.
