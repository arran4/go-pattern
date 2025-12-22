package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type PatternInfo struct {
	Name            string // e.g., "Checker"
	FullName        string // e.g., "Checker Pattern"
	Description     string
	Usage           string
	OutputFilename  string
	ZoomLevels      string // e.g., "[]int{2, 4}" (raw Go code)
	GeneratorFunc   string // e.g., "pattern.NewDemoChecker"
	ReferencesFunc  string
	Order           int
	StructName      string
}

func main() {
	// Parse the parent directory (root of repo)
	fset := token.NewFileSet()
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	var files []*ast.File
	var filenames []string

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
			file, err := parser.ParseFile(fset, entry.Name(), nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}
			if file.Name.Name == "pattern" {
				files = append(files, file)
				filenames = append(filenames, entry.Name())
			}
		}
	}

	patterns := make(map[string]*PatternInfo)

	// 1. Scan for ExampleNew<Name> in _example.go files
	for i, file := range files {
		filename := filenames[i]
		if !strings.HasSuffix(filename, "_example.go") {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if strings.HasPrefix(fn.Name.Name, "ExampleNew") {
				name := strings.TrimPrefix(fn.Name.Name, "ExampleNew")
				if _, exists := patterns[name]; !exists {
					patterns[name] = &PatternInfo{
						Name:       name,
						FullName:   name + " Pattern",
						StructName: name,
					}
				}
				p := patterns[name]

				// Extract Usage from body
				// Since we parsed with ParseComments, we can use ast.Inspect or printer.
				// But we are in "cmd/bootstrap", accessing ".." files.
				// We can read the file section.
				fBytes, err := os.ReadFile(filename)
				if err == nil {
					start := fset.Position(fn.Body.Lbrace).Offset + 1
					end := fset.Position(fn.Body.Rbrace).Offset
					p.Usage = strings.TrimSpace(string(fBytes[start:end]))
				}
			}
		}
	}

	// 2. Scan for Metadata in _example.go files
	for i, file := range files {
		filename := filenames[i]
		if !strings.HasSuffix(filename, "_example.go") {
			continue
		}
		// Look for Variables and Functions
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					vSpec := spec.(*ast.ValueSpec)
					for i, nameIdent := range vSpec.Names {
						name := nameIdent.Name
						// <Name>OutputFilename
						if strings.HasSuffix(name, "OutputFilename") {
							pName := strings.TrimSuffix(name, "OutputFilename")
							if p, ok := patterns[pName]; ok {
								// Extract value
								if i < len(vSpec.Values) {
									if lit, ok := vSpec.Values[i].(*ast.BasicLit); ok {
										p.OutputFilename = strings.Trim(lit.Value, "\"")
									}
								}
							}
						}
						// <Name>ZoomLevels
						if strings.HasSuffix(name, "ZoomLevels") {
							pName := strings.TrimSuffix(name, "ZoomLevels")
							if p, ok := patterns[pName]; ok {
								// Extract value as code
								if i < len(vSpec.Values) {
									// We want the code representation of the composite lit
									v := vSpec.Values[i]
									fBytes, _ := os.ReadFile(filename)
									start := fset.Position(v.Pos()).Offset
									end := fset.Position(v.End()).Offset
									p.ZoomLevels = string(fBytes[start:end])
								}
							}
						}
					}
				}
			}
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
				for _, spec := range genDecl.Specs {
					vSpec := spec.(*ast.ValueSpec)
					for i, nameIdent := range vSpec.Names {
						name := nameIdent.Name
						if strings.HasSuffix(name, "Order") {
							pName := strings.TrimSuffix(name, "Order")
							if p, ok := patterns[pName]; ok {
								if i < len(vSpec.Values) {
									if lit, ok := vSpec.Values[i].(*ast.BasicLit); ok {
										if val, err := strconv.Atoi(lit.Value); err == nil {
											p.Order = val
										}
									}
								}
							}
						}
					}
				}
			}
			if fn, ok := decl.(*ast.FuncDecl); ok {
				// Bootstrap<Name>
				if strings.HasPrefix(fn.Name.Name, "Bootstrap") {
					suffix := strings.TrimPrefix(fn.Name.Name, "Bootstrap")
					// Check if it's <Name> or <Name>References
					if strings.HasSuffix(suffix, "References") {
						pName := strings.TrimSuffix(suffix, "References")
						if p, ok := patterns[pName]; ok {
							p.ReferencesFunc = "pattern." + fn.Name.Name
						}
					} else {
						// It is <Name>
						pName := suffix
						if p, ok := patterns[pName]; ok {
							p.GeneratorFunc = "pattern." + fn.Name.Name
						}
					}
				}
			}
		}
	}

	// 3. Scan for Description in regular .go files (New<Name> or Type <Name>)
	for i, file := range files {
		filename := filenames[i]
		if strings.HasSuffix(filename, "_example.go") || strings.HasSuffix(filename, "_test.go") {
			continue
		}

		// Use go/doc to get comments easily
		// But go/doc works on package level.
		// Let's just inspect AST.

		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					tSpec := spec.(*ast.TypeSpec)
					if p, ok := patterns[tSpec.Name.Name]; ok {
						if genDecl.Doc != nil {
							p.Description = strings.TrimSpace(genDecl.Doc.Text())
						}
					}
				}
			}
			// Also check New<Name> for description if Type doesn't have it (optional, but Type usually has it)
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if strings.HasPrefix(fn.Name.Name, "New") {
					name := strings.TrimPrefix(fn.Name.Name, "New")
					if p, ok := patterns[name]; ok && p.Description == "" {
						if fn.Doc != nil {
							p.Description = strings.TrimSpace(fn.Doc.Text())
						}
					}
				}
			}
		}
	}

	// 4. Cleanup and Defaults
	var sortedPatterns []*PatternInfo
	for _, p := range patterns {
		// Clean Description: "Checker is a..." -> "Alternates between..."
		// The current bootstrap has "Alternates between..."
		// The doc says "Checker is a pattern that alternates..."
		// I'll leave it as is, or try to strip "X is a pattern that ".
		if strings.HasPrefix(p.Description, p.Name+" is a pattern that ") {
			p.Description = strings.ToUpper(p.Description[len(p.Name)+19:len(p.Name)+20]) + p.Description[len(p.Name)+20:]
		}

		// Default Generator
		if p.GeneratorFunc == "" {
			// Check if NewDemo<Name> exists?
			// I assume it does if not Bootstrapped.
			p.GeneratorFunc = fmt.Sprintf("pattern.NewDemo%s", p.Name)
		}

		// Order?
		// User mentioned "order priority".
		// I haven't implemented scanning for "Order" constant yet.
		// I will rely on Name sorting for stability for now, unless I see Explicit Order field.

		sortedPatterns = append(sortedPatterns, p)
	}

	sort.Slice(sortedPatterns, func(i, j int) bool {
		if sortedPatterns[i].Order != sortedPatterns[j].Order {
			return sortedPatterns[i].Order < sortedPatterns[j].Order
		}
		return sortedPatterns[i].Name < sortedPatterns[j].Name
	})

	// Generate the code
	tmpl := template.Must(template.New("gen").Parse(genTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, sortedPatterns); err != nil {
		log.Fatal(err)
	}

	genFile := "cmd/bootstrap/gen_readme.go"
	if err := os.WriteFile(genFile, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(genFile)

	// Run the generated code
	cmd := exec.Command("go", "run", genFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

const genTemplate = `package main

import (
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
    _ "embed"
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

type LabelledGenerator struct {
	Label     string
	Generator func(image.Rectangle) image.Image
}

func main() {
	patterns := []*PatternDemo{
{{- range . }}
		{
			Name:          "{{.FullName}}",
			Description:   "{{.Description}}",
			GoUsageSample: {{printf "%q" .Usage}},
			OutputFilename: "{{.OutputFilename}}",
			Generator: func(b image.Rectangle) image.Image {
				return {{.GeneratorFunc}}(pattern.SetBounds(b))
			},
            {{- if .ZoomLevels }}
			ZoomLevels: {{.ZoomLevels}},
            {{- end }}
            {{- if .ReferencesFunc }}
            References: func() []LabelledGenerator {
                m, order := {{.ReferencesFunc}}()
                var res []LabelledGenerator
                for _, k := range order {
                    gen := m[k]
                    res = append(res, LabelledGenerator{
                        Label: k,
                        Generator: func(b image.Rectangle) image.Image {
                            return gen(pattern.SetBounds(b))
                        },
                    })
                }
                return res
            }(),
            {{- end }}
		},
{{- end }}
	}

	readmeTemplate, err := template.New("readme.md").Parse(string(readmeTemplateRaw))
	if err != nil {
		panic(err)
	}
	f, err := os.Create("readme.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := struct {
		ProjectName string
		Patterns    []PatternDemo
	}{
		ProjectName: "go-pattern",
        Patterns: make([]PatternDemo, 0),
	}
	sz := image.Rect(0, 0, 255, 255)
	for _, p := range patterns {
        data.Patterns = append(data.Patterns, *p)
		DrawDemoPattern(p, sz)
	}
	err = readmeTemplate.Execute(f, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
    log.Println("Generated readme.md successfully")
}

func DrawDemoPattern(pattern *PatternDemo, size image.Rectangle) {
	i := addBorder(pattern.Generate())
	f, err := os.Create(pattern.OutputFilename)
	if err != nil {
		log.Fatalf("Error creating i file: %v", err)
	}
	defer f.Close()
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
`
