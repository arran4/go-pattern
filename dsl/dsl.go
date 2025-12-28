package dsl

// Legacy wrappers to maintain compatibility where possible,
// though we changed the internal structure significantly.

// Pipeline is now just an alias for Node to keep some type compatibility if needed,
// but really we return Node now.
// The old dsl.Pipeline was []Command. This breaks compatibility.
// We should check who uses dsl.Pipeline.
// pkg/pattern-cli/def.go uses it.

// Parse parses the input string into a Node.
func Parse(input string) (Node, error) {
	l := NewLexer(input)
	p := NewParser(l)
	return p.ParsePipeline()
}
