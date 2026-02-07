package dsl

import (
	"fmt"
	"image"
	"image/color"
	"testing"
)

// MockImage is a simple image.Image implementation for testing.
type MockImage struct {
	Name string
}

func (m *MockImage) ColorModel() color.Model { return color.RGBAModel }
func (m *MockImage) Bounds() image.Rectangle { return image.Rect(0, 0, 1, 1) }
func (m *MockImage) At(x, y int) color.Color { return color.RGBA{0, 0, 0, 255} }

func TestExecute(t *testing.T) {
	// Success case
	t.Run("Success", func(t *testing.T) {
		fm := FuncMap{
			"cmd1": func(args []string, input image.Image) (image.Image, error) {
				return &MockImage{Name: "img1"}, nil
			},
			"cmd2": func(args []string, input image.Image) (image.Image, error) {
				in := input.(*MockImage)
				return &MockImage{Name: in.Name + "->img2"}, nil
			},
		}

		pipeline := Pipeline{
			{Name: "cmd1"},
			{Name: "cmd2"},
		}

		initial := &MockImage{Name: "initial"}
		result, err := pipeline.Execute(fm, initial)

		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		resImg, ok := result.(*MockImage)
		if !ok {
			t.Fatalf("Expected MockImage, got %T", result)
		}

		expectedName := "img1->img2"
		if resImg.Name != expectedName {
			t.Errorf("Expected image name %q, got %q", expectedName, resImg.Name)
		}
	})

	// Unknown command case
	t.Run("UnknownCommand", func(t *testing.T) {
		fm := FuncMap{}
		pipeline := Pipeline{
			{Name: "unknown"},
		}

		_, err := pipeline.Execute(fm, nil)
		if err == nil {
			t.Error("Expected error for unknown command, got nil")
		}
	})

	// Command failure case
	t.Run("CommandFailure", func(t *testing.T) {
		fm := FuncMap{
			"fail": func(args []string, input image.Image) (image.Image, error) {
				return nil, fmt.Errorf("intentional failure")
			},
		}
		pipeline := Pipeline{
			{Name: "fail"},
		}

		_, err := pipeline.Execute(fm, nil)
		if err == nil {
			t.Error("Expected error for command failure, got nil")
		}
	})
}
