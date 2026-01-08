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
	"unicode"

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

	if err := generateCLIInit(patterns, "pkg/pattern-cli/init_gen.go"); err != nil {
		log.Fatalf("Error generating CLI init: %v", err)
	}
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
				usage = strings.Trim(usage, "\n")

				pd := PatternDemo{
					Name:          name + " Pattern",
					GoUsageSample: usage,
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
					for _, label := range order {
						if g, ok := refMap[label]; ok {
							pd.Inputs = append(pd.Inputs, LabelledGenerator{
								Label:     label,
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
		if patterns[i].Order != patterns[j].Order {
			return patterns[i].Order < patterns[j].Order
		}
		return patterns[i].Name < patterns[j].Name
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
		log.Fatalf("Error creating image file: %v", err)
	}
	defer func() {
		if e := f.Close(); e != nil {
			log.Fatalf("Error closing image file: %v", e)
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
		log.Fatalf("Unknown image format: %s", e)
	}
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}
	log.Printf("Generated image %s successfully\n", pattern.OutputFilename)
}

func addBorder(img image.Image) image.Image {
	if img == nil {
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

	for _, input := range p.Inputs {
		items = append(items, item{input.Generator(b), input.Label})
	}
	for _, ref := range p.References {
		items = append(items, item{ref.Generator(b), ref.Label})
	}

	for _, step := range p.Steps {
		items = append(items, item{step.Generator(b), step.Label})
	}

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

	if baseImg != nil {
		for _, z := range p.ZoomLevels {
			img := pattern.NewSimpleZoom(baseImg, z, pattern.SetBounds(b))
			items = append(items, item{img, fmt.Sprintf("%dx", z)})
		}
	}

	n := len(items)
	if n == 0 {
		return image.NewRGBA(b)
	}

	face, err := loadFontFace()
	if err != nil {
		log.Fatalf("failed to load font: %v", err)
	}

	var itemWidths []int
	totalW := padding
	for _, it := range items {
		d := &font.Drawer{Face: face}
		w := d.MeasureString(it.label).Ceil()

		iw := sz
		if w > sz {
			iw = w
		}
		itemWidths = append(itemWidths, iw)
		totalW += iw + padding
	}

	totalH := sz + 2*padding + labelHeight

	dst := image.NewRGBA(image.Rect(0, 0, totalW, totalH))
	white := image.NewUniform(color.White)
	draw.Draw(dst, dst.Bounds(), white, image.Point{}, draw.Src)

	currentX := padding
	for i, it := range items {
		iw := itemWidths[i]

		d := &font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(color.Black),
			Face: face,
		}
		labelW := d.MeasureString(it.label).Ceil()
		d.Dot = fixed.P(currentX+(iw-labelW)/2, padding+20)
		d.DrawString(it.label)

		imgX := currentX + (iw - sz)/2
		r := image.Rect(imgX, padding+labelHeight, imgX+sz, padding+labelHeight+sz)
		draw.Draw(dst, r, it.img, b.Min, draw.Src)

		currentX += iw + padding
	}

	return dst
}

func generateCLIInit(demos []PatternDemo, outfile string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		return err
	}

	type Command struct {
		Name      string
		FuncName  string
		Args      []string
		TakesInput bool
	}
	var commands []Command

	for _, pkg := range pkgs {
		for filename, f := range pkg.Files {
			if strings.HasSuffix(filename, "_test.go") || strings.HasSuffix(filename, "_example.go") {
				continue
			}

			ast.Inspect(f, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}
				if !strings.HasPrefix(fn.Name.Name, "New") {
					return true
				}

				cmdName := toSnakeCase(strings.TrimPrefix(fn.Name.Name, "New"))

				var args []string
				takesInput := false

				if fn.Type.Params != nil {
					for i, param := range fn.Type.Params.List {
						// type string
						typeName := ""
						isVariadic := false

						if ident, ok := param.Type.(*ast.Ident); ok {
							typeName = ident.Name
						} else if sel, ok := param.Type.(*ast.SelectorExpr); ok {
							// e.g. color.Color
							if x, ok := sel.X.(*ast.Ident); ok {
								typeName = x.Name + "." + sel.Sel.Name
							}
						} else if ell, ok := param.Type.(*ast.Ellipsis); ok {
							isVariadic = true
							if ident, ok := ell.Elt.(*ast.Ident); ok {
								typeName = "..." + ident.Name
							} else if _, ok := ell.Elt.(*ast.FuncType); ok {
								// ...func(any)
								typeName = "...func(any)"
							}
						}

						// Handle parameter names
						for range param.Names {
							if isVariadic {
								if typeName == "...func(any)" {
									// Skip this arg in CLI requirements
									continue
								}
								// Treat other variadics as unsupported for now, or strings
								args = append(args, typeName)
							} else if i == 0 && typeName == "image.Image" {
								takesInput = true
							} else {
								args = append(args, typeName)
							}
						}
					}
				}

				commands = append(commands, Command{
					Name:      cmdName,
					FuncName:  fn.Name.Name,
					Args:      args,
					TakesInput: takesInput,
				})
				return true
			})
		}
	}

	// Sort commands by name for stable output
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	// Generate file content
	var sb strings.Builder
	sb.WriteString("// Code generated by cmd/bootstrap/main.go; DO NOT EDIT.\n")
	sb.WriteString("package pattern_cli\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"image\"\n")
	sb.WriteString("\t\"strconv\"\n")
	sb.WriteString("\t\"github.com/arran4/go-pattern/dsl\"\n")
	sb.WriteString("\t\"github.com/arran4/go-pattern\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("func RegisterGeneratedCommands(fm dsl.FuncMap) {\n")

	for _, cmd := range commands {
		sb.WriteString(fmt.Sprintf("\tfm[\"%s\"] = func(args []string, input image.Image) (image.Image, error) {\n", cmd.Name))

		// Check arg count
		sb.WriteString(fmt.Sprintf("\t\tif len(args) < %d {\n", len(cmd.Args)))
		sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"%s requires %d arguments\")\n", cmd.Name, len(cmd.Args)))
		sb.WriteString("\t\t}\n")

		// Check support first
		supported := true
		for _, argType := range cmd.Args {
			switch argType {
			case "int", "float64", "bool", "color.Color", "string":
				// supported
			default:
				supported = false
			}
		}

		if supported {
			// Parse args
			callArgs := []string{}
			if cmd.TakesInput {
				sb.WriteString("\t\tif input == nil {\n")
				sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"%s requires an input image\")\n", cmd.Name))
				sb.WriteString("\t\t}\n")
				callArgs = append(callArgs, "input")
			}

			for i, argType := range cmd.Args {
				varName := fmt.Sprintf("arg%d", i)
				switch argType {
				case "int":
					sb.WriteString(fmt.Sprintf("\t\t%s, err := strconv.Atoi(args[%d])\n", varName, i))
					sb.WriteString("\t\tif err != nil {\n")
					sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"argument %d must be int: %%v\", err)\n", i))
					sb.WriteString("\t\t}\n")
					callArgs = append(callArgs, varName)
				case "float64":
					sb.WriteString(fmt.Sprintf("\t\t%s, err := strconv.ParseFloat(args[%d], 64)\n", varName, i))
					sb.WriteString("\t\tif err != nil {\n")
					sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"argument %d must be float: %%v\", err)\n", i))
					sb.WriteString("\t\t}\n")
					callArgs = append(callArgs, varName)
				case "bool":
					sb.WriteString(fmt.Sprintf("\t\t%s, err := strconv.ParseBool(args[%d])\n", varName, i))
					sb.WriteString("\t\tif err != nil {\n")
					sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"argument %d must be bool: %%v\", err)\n", i))
					sb.WriteString("\t\t}\n")
					callArgs = append(callArgs, varName)
				case "color.Color":
					sb.WriteString(fmt.Sprintf("\t\t%s, err := parseColor(args[%d])\n", varName, i))
					sb.WriteString("\t\tif err != nil {\n")
					sb.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"argument %d must be color: %%v\", err)\n", i))
					sb.WriteString("\t\t}\n")
					callArgs = append(callArgs, varName)
				case "string":
					callArgs = append(callArgs, fmt.Sprintf("args[%d]", i))
				}
			}

			sb.WriteString(fmt.Sprintf("\t\treturn pattern.%s(%s), nil\n", cmd.FuncName, strings.Join(callArgs, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("\t\treturn nil, fmt.Errorf(\"command %s has unsupported argument types\")\n", cmd.Name))
		}

		sb.WriteString("\t}\n")
	}

	sb.WriteString("}\n")

	return os.WriteFile(outfile, []byte(sb.String()), 0644)
}

func toSnakeCase(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
