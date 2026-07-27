package pattern

import (
	"image"
	"sort"
)

var (
	// GlobalGenerators is a map of registered generator functions.
	// Keys are the names of the patterns.
	GlobalGenerators = make(map[string]func(image.Rectangle) image.Image)

	// GlobalReferences is a map of registered reference generators.
	// Keys are the names of the patterns.
	GlobalReferences = make(map[string]func() (map[string]func(image.Rectangle) image.Image, []string))
)

// RegisterGenerator registers a pattern generator function with the given name.
func RegisterGenerator(name string, gen func(image.Rectangle) image.Image) {
	GlobalGenerators[name] = gen
}

// RegisterReferences registers a pattern reference generator function with the given name.
func RegisterReferences(name string, refs func() (map[string]func(image.Rectangle) image.Image, []string)) {
	GlobalReferences[name] = refs
}

// AvailableGenerators returns a sorted slice of available pattern generator names.
func AvailableGenerators() []string {
	keys := make([]string, 0, len(GlobalGenerators))
	for k := range GlobalGenerators {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AvailableReferences returns a sorted slice of available pattern reference generator names.
func AvailableReferences() []string {
	keys := make([]string, 0, len(GlobalReferences))
	for k := range GlobalReferences {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
