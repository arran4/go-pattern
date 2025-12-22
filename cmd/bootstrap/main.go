package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"text/template"

	pattern "github.com/arran4/go-pattern"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed "readme.md.gotmpl"
var readmeTemplateRaw []byte

type PatternDemo struct {
	Name           string
	OutputFilename string
	Description    string
	GoUsageSample  string

	Generator func(image.Rectangle) image.Image

	References []LabelledGenerator
	Steps      []LabelledGenerator
	BaseLabel  string
	ZoomLevels []int
}

var Patterns = CreatePatternList()

func CreatePatternList() []*PatternDemo {
	return []*PatternDemo{
		{
			Name:          "Null Pattern",
			Description:   "Undefined RGBA colour.",
			GoUsageSample: "i := pattern.NewNull()\n\tf, err := os.Create(\"null.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Generator: func(bounds image.Rectangle) image.Image {
				return pattern.NewDemoNull(pattern.SetBounds(bounds))
			},
			OutputFilename: "null.png",
		},
		{
			Name:          "Checker Pattern",
			Description:   "Alternates between two colors in a checkerboard fashion.",
			GoUsageSample: "i := pattern.NewChecker(color.Black, color.White)\n\tf, err := os.Create(\"checker.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Generator: func(b image.Rectangle) image.Image {
				return pattern.NewDemoChecker(pattern.SetBounds(b))
			},
			ZoomLevels:     []int{2, 4},
			OutputFilename: "checker.png",
		},
		{
			Name:          "Simple Zoom Pattern",
			Description:   "Zooms in on an underlying image.",
			GoUsageSample: "i := pattern.NewSimpleZoom(pattern.NewChecker(color.Black, color.White), 2)\n\tf, err := os.Create(\"simplezoom.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Generator: func(b image.Rectangle) image.Image {
				return pattern.NewDemoChecker(pattern.SetBounds(b))
			},
			ZoomLevels:     []int{2, 4},
			OutputFilename: "simplezoom.png",
		},
		{
			Name:          "Transposed Pattern",
			Description:   "Transposes the X and Y coordinates of an underlying image.",
			GoUsageSample: "i := pattern.NewTransposed(pattern.NewDemoNull(), 10, 10)\n\tf, err := os.Create(\"transposed.png\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer func() {\n\t\tif e := f.Close(); e != nil {\n\t\t\tpanic(e)\n\t\t}\n\t}()\n\tif err = png.Encode(f, i); err != nil {\n\t\tpanic(err)\n\t}",
			Generator: func(b image.Rectangle) image.Image {
				// Base: simple zoom of a checker (5x)
				baseImg := pattern.NewSimpleZoom(pattern.NewDemoChecker(pattern.SetBounds(b)), 10, pattern.SetBounds(b))
				// Transposed
				return pattern.NewTransposed(baseImg, 5, 5, pattern.SetBounds(b))
			},
			References: []LabelledGenerator{
				{
					Label: "Original",
					Generator: func(b image.Rectangle) image.Image {
						return pattern.NewSimpleZoom(pattern.NewDemoChecker(pattern.SetBounds(b)), 5, pattern.SetBounds(b))
					},
				},
			},
			BaseLabel:      "Transposed",
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
	i := addBorder(pattern.Generate())
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

func addBorder(img image.Image) image.Image {
	b := img.Bounds()
	borderWidth := 5
	nb := image.Rect(0, 0, b.Dx()+2*borderWidth, b.Dy()+2*borderWidth)
	dst := image.NewRGBA(nb)
	draw.Draw(dst, nb, image.NewUniform(color.Black), image.Point{}, draw.Src)
	tr := image.Rect(borderWidth, borderWidth, nb.Dx()-borderWidth, nb.Dy()-borderWidth)
	draw.Draw(dst, tr, img, b.Min, draw.Src)
	return dst
}

func drawLabel(img draw.Image, label string, x, y int) {
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("failed to create face: %v", err)
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(label)
}

type LabelledGenerator struct {
	Label     string
	Generator func(image.Rectangle) image.Image
}

func (p *PatternDemo) Generate() image.Image {
	sz := 150
	b := image.Rect(0, 0, sz, sz)
	padding := 10
	labelHeight := 30

	type item struct {
		img   image.Image
		label string
	}
	var items []item

	// 1. References
	for _, ref := range p.References {
		items = append(items, item{ref.Generator(b), ref.Label})
	}

	// 2. Steps
	for _, step := range p.Steps {
		items = append(items, item{step.Generator(b), step.Label})
	}

	// 3. Base 1x
	baseLabel := p.BaseLabel
	if baseLabel == "" {
		baseLabel = "1x"
	}
	items = append(items, item{p.Generator(b), baseLabel})

	// 4. Zooms
	for _, z := range p.ZoomLevels {
		img := pattern.NewSimpleZoom(p.Generator(b), z, pattern.SetBounds(b))
		items = append(items, item{img, fmt.Sprintf("%dx", z)})
	}

	// Layout
	n := len(items)
	totalW := n*sz + (n+1)*padding
	totalH := sz + 2*padding + labelHeight

	dst := image.NewRGBA(image.Rect(0, 0, totalW, totalH))
	white := image.NewUniform(color.White)
	draw.Draw(dst, dst.Bounds(), white, image.Point{}, draw.Src) // background

	for i, it := range items {
		xOffset := padding + i*(sz+padding)
		drawLabel(dst, it.label, xOffset+5, padding+20)
		r := image.Rect(xOffset, padding+labelHeight, xOffset+sz, padding+labelHeight+sz)
		draw.Draw(dst, r, it.img, b.Min, draw.Src)
	}

	return dst
}
