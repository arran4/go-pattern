package main

import (
	"fmt"
	"image"

	"github.com/arran4/go-pattern"
)

func main() {
    b := image.Rect(10, 10, 20, 20)
    gens := pattern.AvailableGenerators()

    fmt.Printf("Testing %d generators for non-zero origin bounds behavior...\n", len(gens))
    failures := make(map[string]string)

    for _, name := range gens {
        genFn := pattern.GlobalGenerators[name]
        img := genFn(b)

        if img == nil {
            failures[name] = "Returned nil"
            continue
        }

        func(name string, img image.Image) {
            defer func() {
                if r := recover(); r != nil {
                    failures[name] = fmt.Sprintf("Panics on At(): %v", r)
                }
            }()

            bounds := img.Bounds()
            if bounds.Min.X != 10 || bounds.Min.Y != 10 || bounds.Max.X != 20 || bounds.Max.Y != 20 {
                failures[name] = fmt.Sprintf("Got: %v", bounds)
            }
        }(name, img)
    }

    if len(failures) > 0 {
        fmt.Printf("%d generators failed.\n", len(failures))
        for name, reason := range failures {
            fmt.Printf("- %s: %s\n", name, reason)
        }
    } else {
        fmt.Println("All tested pattern generators survived non-zero bounds calls.")
    }
}
