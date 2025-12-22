package main

import (
	_ "embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
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

	Inputs       []LabelledGenerator
	Transformers []LabelledTransformer

	References []LabelledGenerator
	Steps      []LabelledGenerator
	BaseLabel  string
	ZoomLevels []int
	Order      int
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

	patterns, err := discoverPatterns(".")
	if err != nil {
		log.Fatalf("Failed to discover patterns: %v", err)
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
		Patterns:    patterns,
	}
	sz := image.Rect(0, 0, 255, 255)
	for _, p := range patterns {
		DrawDemoPattern(&p, sz)
	}
	err = readmeTemplate.Execute(f, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
	log.Printf("Generated %s successfully\n", fn)
}

func discoverPatterns(root string) ([]PatternDemo, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, root, nil, 0)
	if err != nil {
		return nil, err
	}

	var patterns []PatternDemo

	for _, pkg := range pkgs {
		for filename, f := range pkg.Files {
			// We only care about _example.go files for metadata
			if !strings.HasSuffix(filename, "_example.go") {
				continue
			}

			// To extract source code properly, we need to read the file content
			fileContent, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}

			ast.Inspect(f, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}
				if !strings.HasPrefix(fn.Name.Name, "ExampleNew") {
					return true
				}

				name := strings.TrimPrefix(fn.Name.Name, "ExampleNew")

				// Extract Usage Sample
				start := fset.Position(fn.Body.Lbrace).Offset + 1
				end := fset.Position(fn.Body.Rbrace).Offset
				usage := string(fileContent[start:end])
				// usage = strings.TrimSpace(usage) // Keep indentation or adjust?
				// The original code had tabs. Let's try to dedent if needed,
				// but simplistic extraction is likely fine if formatted.
				// Actually, we should probably strip leading/trailing newlines.
				usage = strings.Trim(usage, "\n")

				pd := PatternDemo{
					Name:          name + " Pattern", // Convention from hardcoded list
					GoUsageSample: usage,
					// Description:   "", // Description was hardcoded. We might need to extract it from comments?
					// For now, let's leave description empty or infer?
					// The hardcoded list had descriptions. The prompt doesn't explicitly say where to get description.
					// "Metadata is extracted... Usage... OutputFilename... ZoomLevels... Order... Custom Generator".
					// It missed Description. I'll check if I can get it from doc comments.
				}

				if fn.Doc != nil {
					pd.Description = strings.TrimSpace(fn.Doc.Text())
				}

				// Look up configuration in the AST of the same file
				pd.OutputFilename = findStringVar(f, name+"OutputFilename")
				pd.ZoomLevels = findIntSliceVar(f, fileContent, fset, name+"ZoomLevels")
				pd.Order = findIntConst(f, name+"Order")
				pd.BaseLabel = findStringConst(f, name+"BaseLabel")

				// Look up Generator in Registry
				if gen, ok := pattern.GlobalGenerators[name]; ok {
					pd.Generator = gen
				} else {
					log.Printf("Warning: No generator found for %s", name)
				}

				// Look up References in Registry
				if refsFunc, ok := pattern.GlobalReferences[name]; ok {
					refMap, order := refsFunc()
					// Map to Inputs/Transformers is tricky because the structure in main.go
					// was specific to Transposed which used Inputs/Transformers.
					// However, the `Generate` function in main.go handles `Inputs`, `References`, `Steps`.
					// `Transposed` example in main.go used `Inputs` and `Transformers`.
					// If `BootstrapTransposedReferences` returns a map, we can put them in `Inputs` or `References`.
					// Let's use `References` for general map items?
					// But `Transposed` logic in `Generate` (main.go) used `Inputs[0]` as base for Transformers.

					// Let's try to adapt.
					// If the pattern has references, we can populate `Inputs` or `References`.
					// `Transposed` had `Original` (Input) and `Transposed` (Transformer output?).
					// But here `refMap` gives us generators.
					// If `BootstrapTransposedReferences` returns "Original" -> func and "Transposed" -> func.
					// We can put them into `Inputs`?

					for _, label := range order {
						if g, ok := refMap[label]; ok {
							pd.Inputs = append(pd.Inputs, LabelledGenerator{
								Label: label,
								Generator: g,
							})
						}
					}

					// If "Transposed" logic specifically needs `Transformers`, we might need more metadata.
					// But the requirement says: "References for patterns are provided by Bootstrap<Name>References functions which return a map of generators and a slice of labels for ordering."
					// This matches `Inputs` or `References` in `PatternDemo` struct usage in `Generate` method:
					// // 1. References
					// for _, input := range p.Inputs { ... }
					// for _, ref := range p.References { ... }
					// So putting them in `Inputs` seems fine to display them.
				}

				patterns = append(patterns, pd)
				return true
			})
		}
	}

	// Sort by Order
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Order < patterns[j].Order
	})

	return patterns, nil
}

func findStringVar(f *ast.File, name string) string {
	var val string
	ast.Inspect(f, func(n ast.Node) bool {
		if gen, ok := n.(*ast.GenDecl); ok && gen.Tok == token.VAR {
			for _, spec := range gen.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for i, ident := range vs.Names {
						if ident.Name == name {
							if len(vs.Values) > i {
								if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
									val, _ = strconv.Unquote(lit.Value)
								}
							}
							return false
						}
					}
				}
			}
		}
		return true
	})
	return val
}

func findIntConst(f *ast.File, name string) int {
	var val int
	ast.Inspect(f, func(n ast.Node) bool {
		if gen, ok := n.(*ast.GenDecl); ok && gen.Tok == token.CONST {
			for _, spec := range gen.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for i, ident := range vs.Names {
						if ident.Name == name {
							if len(vs.Values) > i {
								if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.INT {
									v, _ := strconv.Atoi(lit.Value)
									val = v
								}
							}
							return false
						}
					}
				}
			}
		}
		return true
	})
	return val
}

func findStringConst(f *ast.File, name string) string {
	var val string
	ast.Inspect(f, func(n ast.Node) bool {
		if gen, ok := n.(*ast.GenDecl); ok && gen.Tok == token.CONST {
			for _, spec := range gen.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for i, ident := range vs.Names {
						if ident.Name == name {
							if len(vs.Values) > i {
								if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
									val, _ = strconv.Unquote(lit.Value)
								}
							}
							return false
						}
					}
				}
			}
		}
		return true
	})
	return val
}


func findIntSliceVar(f *ast.File, content []byte, fset *token.FileSet, name string) []int {
	var nums []int
	ast.Inspect(f, func(n ast.Node) bool {
		if gen, ok := n.(*ast.GenDecl); ok && gen.Tok == token.VAR {
			for _, spec := range gen.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for i, ident := range vs.Names {
						if ident.Name == name {
							if len(vs.Values) > i {
								if cl, ok := vs.Values[i].(*ast.CompositeLit); ok {
									for _, elt := range cl.Elts {
										if lit, ok := elt.(*ast.BasicLit); ok && lit.Kind == token.INT {
											v, _ := strconv.Atoi(lit.Value)
											nums = append(nums, v)
										}
									}
								}
							}
							return false
						}
					}
				}
			}
		}
		return true
	})
	return nums
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
	if img == nil {
		// Create a placeholder if image is nil (e.g. if generator returned nil)
		img = image.NewRGBA(image.Rect(0, 0, 150, 150))
	}
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

type LabelledTransformer struct {
	Label       string
	Transformer func(image.Image, image.Rectangle) image.Image
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

	// 1. References / Inputs
	for _, input := range p.Inputs {
		items = append(items, item{input.Generator(b), input.Label})
	}
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

	var baseImg image.Image

	if p.Generator != nil {
		genImg := p.Generator(b)
		if genImg != nil {
			baseImg = genImg
			items = append(items, item{baseImg, baseLabel})
		}
	}
	// The original logic had: if generator != nil { ... } else if len(Inputs) > 0 { ... }
	// In the original Transposed example, Generator was nil? No, `main.go` logic used `Inputs` and `Transformers`.
	// Since we are refactoring, we might lose the `Transformers` logic unless we reconstructed it.
	// But `BootstrapTransposedReferences` returned a map of generators.
	// So we can just display those generators.
	// For Transposed, we now have "Original" and "Transposed" as generators in `Inputs`.
	// So we don't need `Transformers` logic if the generators already do the transformation.
	// The `BootstrapTransposedReferences` I wrote earlier creates a `NewTransposed` from scratch.
	// So displaying them via `Inputs` is sufficient.

	// 4. Zooms
	if baseImg != nil {
		for _, z := range p.ZoomLevels {
			img := pattern.NewSimpleZoom(baseImg, z, pattern.SetBounds(b))
			items = append(items, item{img, fmt.Sprintf("%dx", z)})
		}
	}

	// Layout
	n := len(items)
	if n == 0 {
		return image.NewRGBA(b)
	}

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
