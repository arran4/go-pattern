package main

import (
	_ "embed"
	"flag"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"text/template"

	pattern "github.com/arran4/go-pattern"
)

//go:embed "readme.md.gotmpl"
var readmeTemplateRaw []byte

type PatternDemo struct {
	Name           string
	Create         func(bounds image.Rectangle) image.Image
	OutputFilename string
	Description    string
	GoUsageSample  string
}

var Patterns = CreatePatternList()

func CreatePatternList() []*PatternDemo {
	return []*PatternDemo{
		{
			Name:          "Null Pattern",
			Description:   "Undefined RGBA colour.",
			GoUsageSample: "i := pattern.NewNull()\n\tf, err := os.Create(\"null.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Create: func(bounds image.Rectangle) image.Image {
				return pattern.NewDemoNull(pattern.SetBounds(bounds))
			},
			OutputFilename: "null.png",
		},
		{
			Name:          "Checker Pattern",
			Description:   "Alternates between two colors in a checkerboard fashion.",
			GoUsageSample: "i := pattern.NewChecker(color.Black, color.White)\n\tf, err := os.Create(\"checker.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Create: func(bounds image.Rectangle) image.Image {
				return pattern.NewDemoChecker(pattern.SetBounds(bounds))
			},
			OutputFilename: "checker.png",
		},
		{
			Name:          "Simple Zoom Pattern",
			Description:   "Zooms in on an underlying image.",
			GoUsageSample: "i := pattern.NewSimpleZoom(pattern.NewChecker(color.Black, color.White), 2)\n\tf, err := os.Create(\"simplezoom.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Create: func(bounds image.Rectangle) image.Image {
				return pattern.NewDemoSimpleZoom(pattern.NewDemoChecker(pattern.SetBounds(bounds)), pattern.SetBounds(bounds))
			},
			OutputFilename: "simplezoom.png",
		},
		{
			Name:          "Transposed Pattern",
			Description:   "Transposes the X and Y coordinates of an underlying image.",
			GoUsageSample: "i := pattern.NewTransposed(pattern.NewDemoNull(), 10, 10)\n\tf, err := os.Create(\"transposed.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Create: func(bounds image.Rectangle) image.Image {
				return pattern.NewDemoTransposed(pattern.SetBounds(bounds))
			},
			OutputFilename: "transposed.png",
		},
	}
}

func main() {
	flags := flag.NewFlagSet("bootstrap", flag.ExitOnError)
	fn := "readme.md"
	flags.StringVar(&fn, "filename", fn, "output filename")
	err := flags.Parse(os.Args)
	if err != nil {
		flags.Usage()
		return
	}
	if !flags.Parsed() {
		flags.Usage()
		return
	}
	readmeTemplate, err := template.New("readme.md").Parse(string(readmeTemplateRaw))
	if err != nil {
		panic(err)
	}
	f, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	data := struct {
		ProjectName string
		Patterns    []PatternDemo
	}{
		ProjectName: "go-pattern",
	}
	sz := image.Rect(0, 0, 255, 255)
	for _, pattern := range Patterns {
		data.Patterns = append(data.Patterns, *pattern)
		DrawDemoPattern(pattern, sz)
	}
	err = readmeTemplate.Execute(f, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	log.Printf("Generated %s successfully\n", fn)
}

func DrawDemoPattern(pattern *PatternDemo, size image.Rectangle) {
	i := pattern.Create(size)
	f, err := os.Create(pattern.OutputFilename)
	if err != nil {
		log.Fatalf("Error creating i file: %v", err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Fatalf("Error closing i file: %v", e)
		}
	}()
	e := path.Ext(pattern.OutputFilename)
	switch e {
	case ".png":
		err = png.Encode(f, i)
	case ".jpeg", ".jpg":
		err = jpeg.Encode(f, i, nil)
	case ".gif":
		err = gif.Encode(f, i, nil)
	default:
		log.Fatalf("Unknown i format: %s", e)
	}
	if err != nil {
		log.Fatalf("Error encoding i: %v", err)
	}
	log.Printf("Generated i %s successfully\n", pattern.OutputFilename)
}
