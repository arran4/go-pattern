package dsl

import (
	"testing"
)

func TestParseAndString(t *testing.T) {
	cases := []string{
		"checkers black white | zoom 10 | save out.png",
		"null | transposed 10 20",
		"  checkers   red   blue    |   zoom 5  ",
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
			t.Errorf("Length mismatch: %d vs %d", len(p), len(p2))
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
