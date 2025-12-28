package pattern

import (
	"image"
	"image/color"
)

var (
	RoadOutputFilename = "road.png"
	Road_curvedOutputFilename = "road_curved.png"
	Road_intersectionOutputFilename = "road_intersection.png"
)

// ExampleNewRoad demonstrates a straight road with markings.
func ExampleNewRoad() image.Image {
	return NewRoad(
		SetLanes(2),
		SetWidth(100),
		SetDirection(0), // Horizontal
		SetShape(RoadStraight),
		SetAsphalt(NewAsphalt()),
		SetMarkingColor(color.White),
		SetCenter(127, 127),
	)
}

// ExampleNewRoad_curved demonstrates a curved road.
func ExampleNewRoad_curved() image.Image {
	return NewRoad(
		SetLanes(2),
		SetWidth(80),
		SetDirection(0),
		SetShape(RoadCurved),
		SetCurveRadius(80),
		SetCurveAngle(180), // Semi-circle
		SetAsphalt(NewAsphalt()),
		SetMarkingColor(color.RGBA{255, 255, 0, 255}), // Yellow lines
		SetCenter(127, 200), // Move center down so arc is visible
	)
}

// ExampleNewRoad_intersection demonstrates an intersection.
func ExampleNewRoad_intersection() image.Image {
	return NewRoad(
		SetLanes(4),
		SetWidth(120),
		SetDirection(45), // Rotated intersection
		SetShape(RoadIntersection),
		SetAsphalt(NewAsphalt()),
		SetMarkingColor(color.White),
		SetCenter(127, 127),
	)
}

// ExampleNewAsphalt is already defined in asphalt.go?
// No, I put it there but `asphalt.go` is not an example file (doesn't end in _example.go).
// Wait, `asphalt.go` is a source file.
// But I added `ExampleNewAsphalt` inside `asphalt.go`.
// The bootstrap tool looks for `ExampleNew*` in `_example.go` files usually?
// "The bootstrap tool identifies documentation targets by looking for functions named `ExampleNew<Pattern>` in `_example.go` files"
// So `ExampleNewAsphalt` in `asphalt.go` (if I put it there) might be ignored for doc generation if it's not in an `_example.go` file.
// I should move `ExampleNewAsphalt` to here or `asphalt_example.go`.
// Since I haven't created `asphalt_example.go` and user said "replace road examples", I will put it here.
// I will remove it from `asphalt.go` later or just ignore it there.

// Redefining here would cause conflict if it's exported in `asphalt.go`.
// Let's check `asphalt.go` content. I defined `ExampleNewAsphalt` there.
// I should probably rename `asphalt.go` to `asphalt_example.go` OR move the Example function here.
// But `NewAsphalt` (the generator) should be in `asphalt.go`.
// So I will fix `asphalt.go` to NOT have the Example function, and put it here.
// Or just create `asphalt_example.go`.

// Let's assume I will fix `asphalt.go`.

func GenerateRoad(rect image.Rectangle) image.Image {
	return ExampleNewRoad()
}

func GenerateRoad_curved(rect image.Rectangle) image.Image {
	return ExampleNewRoad_curved()
}

func GenerateRoad_intersection(rect image.Rectangle) image.Image {
	return ExampleNewRoad_intersection()
}

func init() {
	GlobalGenerators["Road"] = GenerateRoad
	GlobalGenerators["Road_curved"] = GenerateRoad_curved
	GlobalGenerators["Road_intersection"] = GenerateRoad_intersection
	// Asphalt generator will be registered in asphalt_example.go if I make it.
}
