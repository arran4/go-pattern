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

type PatternCtor struct {
	Name string
	Func *ast.FuncDecl
}

func main() {
	flags := flag.NewFlagSet("bootstrap", flag.ExitOnError)
	fn := "readme.md"
	flags.StringVar(&fn, "filename", fn, "output filename")
	cliOut := "pattern_cli/init_gen.go"
	flags.StringVar(&cliOut, "cli-out", cliOut, "CLI output filename")
	err := flags.Parse(os.Args)
	if err != nil {
		flags.Usage()
		return
	}
	if !flags.Parsed() {
		flags.Usage()
		return
	}

	patterns, ctors, err := discoverPatterns(".")
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

	if err := generateCLI(cliOut, patterns, ctors); err != nil {
		log.Fatalf("Error generating CLI: %v", err)
	}
}

func discoverPatterns(root string) ([]PatternDemo, map[string]*ast.FuncDecl, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, root, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	var patterns []PatternDemo
	ctors := make(map[string]*ast.FuncDecl)

	for _, pkg := range pkgs {
		for filename, f := range pkg.Files {
			// Find constructors in all files (excluding _test and _example)
			if !strings.HasSuffix(filename, "_test.go") && !strings.HasSuffix(filename, "_example.go") {
				ast.Inspect(f, func(n ast.Node) bool {
					fn, ok := n.(*ast.FuncDecl)
					if !ok {
						return true
					}
					if strings.HasPrefix(fn.Name.Name, "New") {
						ctors[fn.Name.Name] = fn
					}
					return true
				})
			}

			// We only care about _example.go files for metadata
			if !strings.HasSuffix(filename, "_example.go") {
				continue
			}

			// To extract source code properly, we need to read the file content
			fileContent, err := os.ReadFile(filename)
			if err != nil {
				return nil, nil, err
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
				usage = strings.Trim(usage, "\n")

				pd := PatternDemo{
					Name:          name + " Pattern",
					GoUsageSample: usage,
				}

				if fn.Doc != nil {
					pd.Description = strings.TrimSpace(fn.Doc.Text())
				}

				pd.OutputFilename = findStringVar(f, name+"OutputFilename")
				pd.ZoomLevels = findIntSliceVar(f, fileContent, fset, name+"ZoomLevels")
				pd.Order = findIntConst(f, name+"Order")
				pd.BaseLabel = findStringConst(f, name+"BaseLabel")

				if gen, ok := pattern.GlobalGenerators[name]; ok {
					pd.Generator = gen
				} else {
					log.Printf("Warning: No generator found for %s", name)
				}

				if refsFunc, ok := pattern.GlobalReferences[name]; ok {
					refMap, order := refsFunc()
					for _, label := range order {
						if g, ok := refMap[label]; ok {
							pd.Inputs = append(pd.Inputs, LabelledGenerator{
								Label: label,
								Generator: g,
							})
						}
					}
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

	return patterns, ctors, nil
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

func loadFontFace() (font.Face, error) {
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})
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

	face, err := loadFontFace()
	if err != nil {
		log.Fatalf("failed to load font: %v", err)
	}

	// Measure widths
	var itemWidths []int
	totalW := padding // Initial padding
	for _, it := range items {
		d := &font.Drawer{Face: face}
		w := d.MeasureString(it.label).Ceil()

		iw := sz
		if w > sz {
			iw = w // expand cell if text is wider
		}
		itemWidths = append(itemWidths, iw)
		totalW += iw + padding
	}

	totalH := sz + 2*padding + labelHeight

	dst := image.NewRGBA(image.Rect(0, 0, totalW, totalH))
	white := image.NewUniform(color.White)
	draw.Draw(dst, dst.Bounds(), white, image.Point{}, draw.Src) // background

	currentX := padding
	for i, it := range items {
		iw := itemWidths[i]

		// Draw Label
		d := &font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(color.Black),
			Face: face,
		}
		labelW := d.MeasureString(it.label).Ceil()
		d.Dot = fixed.P(currentX+(iw-labelW)/2, padding+20) // Center align
		d.DrawString(it.label)

		// Draw Image (Centered in iw)
		imgX := currentX + (iw - sz)/2
		r := image.Rect(imgX, padding+labelHeight, imgX+sz, padding+labelHeight+sz)
		draw.Draw(dst, r, it.img, b.Min, draw.Src)

		currentX += iw + padding
	}

	return dst
}

func generateCLI(filename string, patterns []PatternDemo, ctors map[string]*ast.FuncDecl) error {
	var sb strings.Builder
	sb.WriteString("// Code generated by cmd/bootstrap; DO NOT EDIT.\n")
	sb.WriteString("package pattern_cli\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"image\"\n")
	sb.WriteString("\t\"strconv\"\n")
	sb.WriteString("\t\"github.com/arran4/go-pattern\"\n")
	sb.WriteString("\t\"github.com/arran4/go-pattern/dsl\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("func init() {\n")
	sb.WriteString("\tRegisterGeneratedCommands = func(fm dsl.FuncMap) {\n")

	for _, p := range patterns {
		rawName := strings.TrimSuffix(p.Name, " Pattern")
		cmdName := strings.ToLower(rawName)
		ctorName := "New" + rawName

		// Handle EdgeDetect -> NewEdgeDetect mismatch if needed (case sensitivity)
		// Assuming pattern.Name is derived from "ExampleNew<Name>"

		ctor, ok := ctors[ctorName]
		if !ok {
			// Try case insensitive match if direct match fails?
			for n, c := range ctors {
				if strings.EqualFold(n, ctorName) {
					ctor = c
					ok = true
					break
				}
			}
		}

		if !ok {
			sb.WriteString(fmt.Sprintf("\t\t// Warning: Constructor %s not found for %s\n", ctorName, cmdName))
			continue
		}

		sb.WriteString(fmt.Sprintf("\t\tfm[\"%s\"] = func(args []string, input image.Image) (image.Image, error) {\n", cmdName))

		// Build arguments
		var callArgs []string
		argIdx := 0
		var paramCheck strings.Builder

		// Filter relevant params (skip options for now, handle Image, Color, int/float)
		if ctor.Type.Params != nil {
			for _, field := range ctor.Type.Params.List {
				typeStr := getType(field.Type)

				// Handle multiple names for same type (e.g. x, y int)
				names := field.Names
				if len(names) == 0 {
					// Anonymous param, treat as 1
					names = []*ast.Ident{{Name: "_"}}
				}

				for range names {
					if typeStr == "image.Image" {
						// Usually input image.
						// Heuristic: If it's the first param, it *might* be the input.
						// Or if it's named 'input', 'img', 'src'.
						// For now, let's assume if it takes an image, it's the pipeline input.
						// But if we have multiple images? (e.g. blend?)
						// CLI usually pipes one image.
						// If constructor needs image, pass 'input'.
						// If 'input' is nil, we should check.

						paramCheck.WriteString("\t\t\tif input == nil {\n\t\t\t\treturn nil, fmt.Errorf(\"input image required\")\n\t\t\t}\n")
						callArgs = append(callArgs, "input")

					} else if typeStr == "color.Color" {
						paramCheck.WriteString(fmt.Sprintf("\t\t\tif len(args) <= %d {\n\t\t\t\treturn nil, fmt.Errorf(\"argument %d (color) missing\")\n\t\t\t}\n", argIdx, argIdx+1))
						paramCheck.WriteString(fmt.Sprintf("\t\t\tc%d, err := parseColor(args[%d])\n", argIdx, argIdx))
						paramCheck.WriteString("\t\t\tif err != nil {\n\t\t\t\treturn nil, err\n\t\t\t}\n")
						callArgs = append(callArgs, fmt.Sprintf("c%d", argIdx))
						argIdx++

					} else if typeStr == "int" {
						paramCheck.WriteString(fmt.Sprintf("\t\t\tif len(args) <= %d {\n\t\t\t\treturn nil, fmt.Errorf(\"argument %d (int) missing\")\n\t\t\t}\n", argIdx, argIdx+1))
						paramCheck.WriteString(fmt.Sprintf("\t\t\ti%d, err := strconv.Atoi(args[%d])\n", argIdx, argIdx))
						paramCheck.WriteString("\t\t\tif err != nil {\n\t\t\t\treturn nil, fmt.Errorf(\"invalid int argument: %%w\", err)\n\t\t\t}\n")
						callArgs = append(callArgs, fmt.Sprintf("i%d", argIdx))
						argIdx++

					} else if typeStr == "float64" {
						paramCheck.WriteString(fmt.Sprintf("\t\t\tif len(args) <= %d {\n\t\t\t\treturn nil, fmt.Errorf(\"argument %d (float) missing\")\n\t\t\t}\n", argIdx, argIdx+1))
						paramCheck.WriteString(fmt.Sprintf("\t\t\tf%d, err := strconv.ParseFloat(args[%d], 64)\n", argIdx, argIdx))
						paramCheck.WriteString("\t\t\tif err != nil {\n\t\t\t\treturn nil, fmt.Errorf(\"invalid float argument: %%w\", err)\n\t\t\t}\n")
						callArgs = append(callArgs, fmt.Sprintf("f%d", argIdx))
						argIdx++

					} else if strings.HasPrefix(typeStr, "...") || strings.HasPrefix(typeStr, "[]") {
						// Ignored variadic/slice for now
						// callArgs = append(callArgs, "nil") // or empty slice?
						// But if it's variadic options ...Option, we can just omit if it's the last arg.
						// If it's []image.Point (Voronoi), we can't easily support it via simple CLI yet.
						// Skip for now or print warning?

						// If variadic '...', Go allows passing nothing.
						if strings.HasPrefix(typeStr, "...") {
							// do nothing
						} else {
							sb.WriteString(fmt.Sprintf("\t\t\t// Unsupported arg type: %s\n", typeStr))
							sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"command %s has unsupported arguments\")\n", cmdName))
							callArgs = nil // abort
							break
						}
					} else {
						// Unsupported type (e.g. structs, maps)
						sb.WriteString(fmt.Sprintf("\t\t\t// Unsupported arg type: %s\n", typeStr))
						sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"command %s has unsupported arguments\")\n", cmdName))
						callArgs = nil // abort
						break
					}
				}
				if callArgs == nil { break }
			}
		}

		if callArgs == nil && (ctor.Type.Params == nil || len(ctor.Type.Params.List) == 0) {
			callArgs = []string{}
		}

		if callArgs != nil {
			sb.WriteString(paramCheck.String())
			// Must return (image.Image, error). Assuming pattern constructor returns image.Image or (image.Image, error)?
			// Most pattern constructors return just image.Image (e.g. NewChecker).
			// Let's check. pattern.NewChecker returns *Checker (which implements Image).
			// So we need to wrap in `return ..., nil`.
			sb.WriteString(fmt.Sprintf("\t\t\treturn pattern.%s(%s), nil\n", ctor.Name.Name, strings.Join(callArgs, ", ")))
		} else {
			// Stub to avoid missing return error
			sb.WriteString("\t\t\treturn nil, fmt.Errorf(\"not implemented yet\")\n")
		}

		sb.WriteString("\t\t}\n")
	}

	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")
	sb.WriteString("var RegisterGeneratedCommands func(dsl.FuncMap)\n")

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func getType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return getType(t.X) + "." + t.Sel.Name
	case *ast.Ellipsis:
		return "..." + getType(t.Elt)
	case *ast.StarExpr:
		return "*" + getType(t.X)
	case *ast.ArrayType:
		return "[]" + getType(t.Elt)
	default:
		return fmt.Sprintf("%T", t)
	}
}
