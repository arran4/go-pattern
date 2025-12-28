package pattern

import (
	"image"
	"image/color"
)

var (
	RoadOutputFilename = "road.png"
	Road_curvedOutputFilename = "road_curved.png"
	Road_intersectionOutputFilename = "road_intersection.png"
	AsphaltOutputFilename = "asphalt.png"
)

// ExampleNewAsphalt demonstrates the Micro view (Asphalt concrete texture).
func ExampleNewAsphalt() image.Image {
	return NewAsphalt()
}

// ExampleNewRoad demonstrates the Macro view (Straight Road).
func ExampleNewRoad() image.Image {
	return NewRoad(
		SetLanes(2),
		SetWidth(120),
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
		SetWidth(100),
		SetDirection(0),
		SetShape(RoadCurved),
		SetCurveRadius(100),
		SetCurveAngle(180),
		SetAsphalt(NewAsphalt()),
		SetMarkingColor(color.RGBA{255, 255, 0, 255}), // Yellow
		SetCenter(127, 240), // Anchor at bottom to show arc
	)
}

// ExampleNewRoad_intersection demonstrates an intersection.
func ExampleNewRoad_intersection() image.Image {
	return NewRoad(
		SetLanes(4),
		SetWidth(80),
		SetDirection(45), // Rotated
		SetShape(RoadIntersection),
		SetAsphalt(NewAsphalt()),
		SetMarkingColor(color.White),
		SetCenter(127, 127),
	)
}

func GenerateRoad(rect image.Rectangle) image.Image {
	return ExampleNewRoad()
}

func GenerateRoad_curved(rect image.Rectangle) image.Image {
	return ExampleNewRoad_curved()
}

func GenerateRoad_intersection(rect image.Rectangle) image.Image {
	return ExampleNewRoad_intersection()
}

func GenerateAsphalt(rect image.Rectangle) image.Image {
	return ExampleNewAsphalt()
}

func init() {
	GlobalGenerators["Road"] = GenerateRoad
	GlobalGenerators["Road_curved"] = GenerateRoad_curved
	GlobalGenerators["Road_intersection"] = GenerateRoad_intersection
	GlobalGenerators["Asphalt"] = GenerateAsphalt
}
