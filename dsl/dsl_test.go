package dsl

import (
	"fmt"
	"image"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `cmd arg1 arg2 | next key=val (sub cmd) ^ (other)`
	l := NewLexer(input)

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "cmd"},
		{IDENT, "arg1"},
		{IDENT, "arg2"},
		{PIPE, "|"},
		{IDENT, "next"},
		{IDENT, "key"},
		{EQUALS, "="},
		{IDENT, "val"},
		{LPAREN, "("},
		{IDENT, "sub"},
		{IDENT, "cmd"},
		{RPAREN, ")"},
		{CARET, "^"},
		{LPAREN, "("},
		{IDENT, "other"},
		{RPAREN, ")"},
		{EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestParserPipeline(t *testing.T) {
	input := `cmd1 arg1 | cmd2 arg2`
	l := NewLexer(input)
	p := NewParser(l)
	node, err := p.ParsePipeline()
	if err != nil {
		t.Fatalf("ParsePipeline error: %v", err)
	}

	pipe, ok := node.(*PipelineNode)
	if !ok {
		t.Fatalf("Expected PipelineNode, got %T", node)
	}
	if len(pipe.Nodes) != 2 {
		t.Fatalf("Expected 2 nodes, got %d", len(pipe.Nodes))
	}

	cmd1, ok := pipe.Nodes[0].(*CommandNode)
	if !ok || cmd1.Name != "cmd1" {
		t.Errorf("First node not cmd1")
	}
	cmd2, ok := pipe.Nodes[1].(*CommandNode)
	if !ok || cmd2.Name != "cmd2" {
		t.Errorf("Second node not cmd2")
	}
}

func TestParserGroupAndOp(t *testing.T) {
	input := `(a) ^ (b)`
	l := NewLexer(input)
	p := NewParser(l)
	node, err := p.ParsePipeline()
	if err != nil {
		t.Fatalf("ParsePipeline error: %v", err)
	}

	bin, ok := node.(*BinaryNode)
	if !ok {
		t.Fatalf("Expected BinaryNode, got %T", node)
	}
	if bin.Operator != "^" {
		t.Errorf("Expected ^, got %s", bin.Operator)
	}

	// Left should be GroupNode(CommandNode(a))
	leftGroup, ok := bin.Left.(*GroupNode)
	if !ok {
		t.Fatalf("Left not GroupNode")
	}
	leftCmd, ok := leftGroup.Inner.(*CommandNode)
	if !ok || leftCmd.Name != "a" {
		t.Errorf("Left inner not command a")
	}
}

func TestEvaluator(t *testing.T) {
	fm := make(FuncMap)
	fm["mkimg"] = func(args []string, input image.Image) (image.Image, error) {
		return image.NewRGBA(image.Rect(0, 0, 10, 10)), nil
	}
	fm["check"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("input is nil")
		}
		return input, nil
	}
	fm["op_xor"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("xor needs 1 arg")
		}
		return image.NewRGBA(image.Rect(0, 0, 10, 10)), nil
	}

	ctx := NewContext(fm)

	// Test basic pipeline
	input := `mkimg | check`
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	_, err = Execute(p, ctx, nil)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	// Test XOR
	input2 := `(mkimg) ^ (mkimg)`
	p2, err := Parse(input2)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	_, err = Execute(p2, ctx, nil)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
}

func TestJoinArgs(t *testing.T) {
	fm := make(FuncMap)
	fm["gen"] = func(args []string, input image.Image) (image.Image, error) {
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), nil
	}
	fm["join"] = func(args []string, input image.Image) (image.Image, error) {
		// Expect args[0] = mode, args[1] = handle
		if len(args) < 2 {
			return nil, fmt.Errorf("missing args")
		}
		if args[0] != "overlay" {
			return nil, fmt.Errorf("wrong mode")
		}
		return nil, nil
	}

	ctx := NewContext(fm)
	input := `join overlay (gen)`
	p, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	_, err = Execute(p, ctx, nil)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
}
