package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

func main() {
	// Try to generate using Nano Banana (Gemini)
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey != "" {
		fmt.Println("GEMINI_API_KEY found. Attempting to generate assets using Nano Banana...")
		if err := generateWithGemini(apiKey); err != nil {
			fmt.Printf("Error generating with Gemini: %v. Falling back to programmatic generation.\n", err)
			generateProgrammatic()
		} else {
			fmt.Println("Successfully generated assets using Nano Banana.")
		}
	} else {
		fmt.Println("GEMINI_API_KEY not found. Skipping Nano Banana generation. Using programmatic generation.")
		generateProgrammatic()
	}
}

func generateWithGemini(apiKey string) error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Generate Gopher
	fmt.Println("Generating gopher.png...")
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-image", genai.Text("A pixel art style image of a blue Go Gopher, facing slightly to the right, on a transparent background."), nil)
	if err != nil {
		return fmt.Errorf("failed to generate gopher: %w", err)
	}
	if err := saveGeminiImage(resp, "assets/gopher.png"); err != nil {
		return fmt.Errorf("failed to save gopher: %w", err)
	}

	// Generate GO Text
	fmt.Println("Generating go.png...")
	resp, err = client.Models.GenerateContent(ctx, "gemini-2.5-flash-image", genai.Text("The text 'GO' in 3D blue block letters, isometric view, on a white background."), nil)
	if err != nil {
		return fmt.Errorf("failed to generate go text: %w", err)
	}
	if err := saveGeminiImage(resp, "assets/go.png"); err != nil {
		return fmt.Errorf("failed to save go text: %w", err)
	}

	return nil
}

func saveGeminiImage(resp *genai.GenerateContentResponse, filename string) error {
	if resp == nil || len(resp.Candidates) == 0 {
		return fmt.Errorf("no candidates returned")
	}

	// Iterate through parts to find the image
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			// Ensure directory exists
			if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
				return err
			}
			// Write bytes to file
			return os.WriteFile(filename, part.InlineData.Data, 0644)
		}
	}
	return fmt.Errorf("no inline image data found in response")
}

func generateProgrammatic() {
	generateGopher()
	generateGoText()
}

func generateGopher() {
	// Simple Gopher Pixel Art (Upscaled)
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

	writePngUpscaled("assets/gopher.png", data, 10) // Scale 10x
}

func generateGoText() {
	// "GO" text
	T := color.RGBA{0, 0, 0, 0}
	C := color.RGBA{0x00, 0xAD, 0xD8, 0xff} // Blue

	// 5x7 font roughly
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

	writePngUpscaled("assets/go.png", data, 20) // Scale 20x
}

func writePngUpscaled(filename string, data [][]color.Color, scale int) {
	height := len(data)
	width := len(data[0])

	img := image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))
	for y, row := range data {
		for x, col := range row {
			// Fill the scaled block
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.Set(x*scale+dx, y*scale+dy, col)
				}
			}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
