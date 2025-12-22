package pattern

import "image"

var (
	GlobalGenerators = make(map[string]func(image.Rectangle) image.Image)
	GlobalReferences = make(map[string]func() (map[string]func(image.Rectangle) image.Image, []string))
)

func RegisterGenerator(name string, gen func(image.Rectangle) image.Image) {
	GlobalGenerators[name] = gen
}

func RegisterReferences(name string, refs func() (map[string]func(image.Rectangle) image.Image, []string)) {
	GlobalReferences[name] = refs
}
