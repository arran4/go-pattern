package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	generateGopher()
	generateGoText()
}

func generateGopher() {
	// Simple Gopher Pixel Art
	// Colors
	T := color.RGBA{0, 0, 0, 0}             // Transparent
	B := color.RGBA{0x00, 0xAD, 0xD8, 0xff} // Gopher Blue
	Y := color.RGBA{0xFC, 0xE8, 0x83, 0xff} // Yellow/Beige (Snout/Hands/Feet)
	W := color.White                        // Eye White
	K := color.Black                        // Eye Pupil / Nose / Tooth

	// 15x18 Grid
	data := [][]color.Color{
		{T, T, T, T, T, B, B, B, B, B, T, T, T, T, T},
		{T, T, T, T, B, B, B, B, B, B, B, T, T, T, T},
		{T, T, T, B, B, B, B, B, B, B, B, B, T, T, T},
		{T, T, B, B, B, B, B, B, B, B, B, B, B, T, T},
		{T, T, B, B, B, W, W, B, B, W, W, B, B, T, T},
		{T, T, B, B, B, W, K, B, B, W, K, B, B, T, T},
		{T, T, B, B, B, B, B, B, B, B, B, B, B, T, T},
		{T, T, T, B, B, Y, Y, K, Y, Y, B, B, T, T, T},
		{T, T, T, B, B, Y, Y, Y, Y, Y, B, B, T, T, T},
		{T, T, B, B, Y, W, Y, K, Y, W, Y, B, B, T, T}, // With teeth
		{T, B, B, B, Y, Y, Y, Y, Y, Y, Y, B, B, B, T},
		{T, B, B, B, B, B, B, B, B, B, B, B, B, B, T},
		{B, B, B, B, B, B, B, B, B, B, B, B, B, B, B},
		{B, B, B, B, B, B, B, B, B, B, B, B, B, B, B},
		{B, B, B, B, B, B, B, B, B, B, B, B, B, B, B},
		{B, B, T, Y, Y, T, T, T, T, T, Y, Y, T, B, B},
		{B, B, T, Y, Y, T, T, T, T, T, Y, Y, T, B, B},
		{B, B, T, T, T, T, T, T, T, T, T, T, T, B, B},
	}

	writePng("assets/gopher.png", data)
}

func generateGoText() {
	// "GO" text
	T := color.RGBA{0, 0, 0, 0}
	C := color.RGBA{0x00, 0xAD, 0xD8, 0xff} // Blue

	// 5x7 font roughly
	// G
	//  XXX
	// X
	// X  XX
	// X   X
	//  XXX

	// O
	//  XXX
	// X   X
	// X   X
	// X   X
	//  XXX

	data := [][]color.Color{
		{T, T, T, T, T, T, T, T, T, T, T, T, T, T, T},
		{T, T, C, C, C, T, T, T, T, C, C, C, T, T, T},
		{T, C, T, T, T, T, T, T, C, T, T, T, C, T, T},
		{T, C, T, T, T, T, T, T, C, T, T, T, C, T, T},
		{T, C, T, C, C, T, T, T, C, T, T, T, C, T, T},
		{T, C, T, T, C, T, T, T, C, T, T, T, C, T, T},
		{T, T, C, C, C, T, T, T, T, C, C, C, T, T, T},
		{T, T, T, T, T, T, T, T, T, T, T, T, T, T, T},
	}

	writePng("assets/go.png", data)
}

func writePng(filename string, data [][]color.Color) {
	height := len(data)
	width := len(data[0])

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y, row := range data {
		for x, col := range row {
			img.Set(x, y, col)
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
