package pattern

import (
	"image"
	"reflect"
	"testing"
)

func TestAvailableGenerators(t *testing.T) {
	// Save the original map to restore it later
	originalGenerators := GlobalGenerators
	defer func() {
		GlobalGenerators = originalGenerators
	}()

	// Clear the map for testing
	GlobalGenerators = make(map[string]func(image.Rectangle) image.Image)

	// Register some dummy generators out of order
	RegisterGenerator("Zebra", nil)
	RegisterGenerator("Alpha", nil)
	RegisterGenerator("Gamma", nil)
	RegisterGenerator("Beta", nil)

	expected := []string{"Alpha", "Beta", "Gamma", "Zebra"}
	actual := AvailableGenerators()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("AvailableGenerators() = %v, want %v", actual, expected)
	}
}

func TestAvailableReferences(t *testing.T) {
	// Save the original map to restore it later
	originalReferences := GlobalReferences
	defer func() {
		GlobalReferences = originalReferences
	}()

	// Clear the map for testing
	GlobalReferences = make(map[string]func() (map[string]func(image.Rectangle) image.Image, []string))

	// Register some dummy references out of order
	RegisterReferences("Xray", nil)
	RegisterReferences("Delta", nil)
	RegisterReferences("Echo", nil)
	RegisterReferences("Charlie", nil)

	expected := []string{"Charlie", "Delta", "Echo", "Xray"}
	actual := AvailableReferences()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("AvailableReferences() = %v, want %v", actual, expected)
	}
}
