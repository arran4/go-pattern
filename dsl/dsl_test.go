package dsl

import (
	"testing"
)

func TestParseAndString(t *testing.T) {
	cases := []string{
		"checkers black white | zoom 10 | save out.png",
		"null | transposed 10 20",
		"  checkers   red   blue    |   zoom 5  ",
		"",
		"   ",
		"|",
		"cmd |",
		"| cmd",
		"cmd1 arg1 | cmd2 arg2 arg3",
		"cmd1 | cmd2 | cmd3",
	}

	for _, c := range cases {
		p, err := Parse(c)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", c, err)
			continue
		}

		generated := p.String()

		// Parse the generated string again to ensure cycle consistency
		p2, err := Parse(generated)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", generated, err)
			continue
		}

		if len(p) != len(p2) {
			t.Errorf("Length mismatch for %q: %d vs %d", c, len(p), len(p2))
		}

		for i := range p {
			if p[i].Name != p2[i].Name {
				t.Errorf("Command name mismatch at %d: %q vs %q", i, p[i].Name, p2[i].Name)
			}
			if len(p[i].Args) != len(p2[i].Args) {
				t.Errorf("Args length mismatch at %d: %d vs %d", i, len(p[i].Args), len(p2[i].Args))
			}
			for j := range p[i].Args {
				if p[i].Args[j] != p2[i].Args[j] {
					t.Errorf("Arg mismatch at %d,%d: %q vs %q", i, j, p[i].Args[j], p2[i].Args[j])
				}
			}
		}
	}
}

func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected int // number of commands
	}{
		{"", 0},
		{"   ", 0},
		{"|", 0},
		{" | ", 0},
		{"cmd", 1},
		{"cmd arg", 1},
		{"cmd | cmd", 2},
		{"cmd |", 1},
		{"| cmd", 1},
		{"cmd || cmd", 2}, // empty command in middle should be ignored
	}

	for _, tt := range tests {
		p, err := Parse(tt.input)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if len(p) != tt.expected {
			t.Errorf("Parse(%q) expected %d commands, got %d", tt.input, tt.expected, len(p))
		}
	}
}
