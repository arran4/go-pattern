package pattern_cli

import (
	"bufio"
	"fmt"
	"github.com/arran4/go-pattern/dsl"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/arran4/go-pattern"
	"golang.org/x/image/colornames"
)

// Repl is a subcommand `pattern-cli repl`
func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	funcMap := make(dsl.FuncMap)
	registerCommands(funcMap)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "exit" || input == "quit" {
			break
		}
		if err := process(input, funcMap); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Print("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}
}

// Run is a subcommand `pattern-cli run`
func Run(pipeline string) {
	funcMap := make(dsl.FuncMap)
	registerCommands(funcMap)
	if err := process(pipeline, funcMap); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func process(input string, fm dsl.FuncMap) error {
	p, err := dsl.Parse(input)
	if err != nil {
		return err
	}
	_, err = p.Execute(fm, nil)
	return err
}

func registerCommands(fm dsl.FuncMap) {
	fm["checkers"] = func(args []string, input image.Image) (image.Image, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("checkers requires 2 color arguments")
		}
		c1, err := parseColor(args[0])
		if err != nil {
			return nil, err
		}
		c2, err := parseColor(args[1])
		if err != nil {
			return nil, err
		}
		return pattern.NewChecker(c1, c2), nil
	}

	fm["zoom"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("zoom requires an input image")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("zoom requires a factor argument")
		}
		factor, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, fmt.Errorf("invalid zoom factor: %v", err)
		}
		return pattern.NewSimpleZoom(input, factor), nil
	}

	fm["transposed"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("transposed requires an input image")
		}
		x, y := 0, 0
		var err error
		if len(args) >= 1 {
			x, err = strconv.Atoi(args[0])
			if err != nil {
				return nil, fmt.Errorf("invalid x offset: %v", err)
			}
		}
		if len(args) >= 2 {
			y, err = strconv.Atoi(args[1])
			if err != nil {
				return nil, fmt.Errorf("invalid y offset: %v", err)
			}
		}
		return pattern.NewTransposed(input, x, y), nil
	}

	fm["null"] = func(args []string, input image.Image) (image.Image, error) {
		return pattern.NewNull(), nil
	}

	fm["save"] = func(args []string, input image.Image) (image.Image, error) {
		if input == nil {
			return nil, fmt.Errorf("save requires an input image")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("save requires a filename argument")
		}
		filename := args[0]
		f, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if strings.HasSuffix(filename, ".png") {
			if err := png.Encode(f, input); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsupported file format: %s", filename)
		}
		fmt.Printf("Saved to %s\n", filename)
		return input, nil
	}
}

func parseColor(s string) (color.Color, error) {
	if c, ok := colornames.Map[s]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("unknown color: %s", s)
}
