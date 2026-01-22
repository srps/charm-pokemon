package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
)

type ConversionConfig struct {
	Width int
}

func main() {
	spriteDirs := []string{
		"assets/sprites/standard",
		"assets/sprites/shiny",
	}
	outputDir := "assets/art"

	os.MkdirAll(outputDir, 0755)

	config := ConversionConfig{
		Width: 40,
	}

	for _, spriteDir := range spriteDirs {
		files, err := filepath.Glob(filepath.Join(spriteDir, "*.png"))
		if err != nil {
			panic(err)
		}

		isShiny := strings.Contains(spriteDir, "shiny")
		suffix := ""
		if isShiny {
			suffix = "_shiny"
		}

		for _, file := range files {
			id := strings.TrimSuffix(filepath.Base(file), ".png")
			img, err := loadImg(file)
			if err != nil {
				fmt.Printf("Error loading %s: %v\n", id, err)
				continue
			}

			// Generate Half-block ASCII
			ascii := imageToHalfBlocks(img, config)
			asciiPath := filepath.Join(outputDir, id+suffix+".ascii")
			os.WriteFile(asciiPath, []byte(ascii), 0644)

			// Generate Sixel
			sixelArt := imageToSixel(img, config)
			sixelPath := filepath.Join(outputDir, id+suffix+".sixel")
			os.WriteFile(sixelPath, []byte(sixelArt), 0644)

			fmt.Printf("Converted %s%s\n", id, suffix)
		}
	}
}

func loadImg(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func imageToHalfBlocks(img image.Image, config ConversionConfig) string {
	bounds := img.Bounds()
	ratio := float64(bounds.Dy()) / float64(bounds.Dx())
	// Divide by 2 to account for terminal character aspect ratio (chars are ~2x tall as wide)
	height := int(float64(config.Width) * ratio / 2.0)

	// Resize image
	newRect := image.Rect(0, 0, config.Width, height*2)
	resized := image.NewRGBA(newRect)
	draw.ApproxBiLinear.Scale(resized, newRect, img, bounds, draw.Over, nil)
	newBounds := resized.Bounds()

	var buf bytes.Buffer
	for y := 0; y < newBounds.Dy(); y += 2 {
		for x := 0; x < newBounds.Dx(); x++ {
			c1 := resized.At(x, y)
			c2 := resized.At(x, y+1)

			_, _, _, a1 := c1.RGBA()
			_, _, _, a2 := c2.RGBA()

			// Check transparency (alpha < 50%)
			isT1 := a1 < 32768
			isT2 := a2 < 32768

			fg, _ := colorful.MakeColor(c1)
			bg, _ := colorful.MakeColor(c2)

			if isT1 && isT2 {
				buf.WriteString(" ")
			} else if isT1 {
				// Only bottom pixel exists
				buf.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm▄\x1b[0m", int(bg.R*255), int(bg.G*255), int(bg.B*255)))
			} else if isT2 {
				// Only top pixel exists
				buf.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm▀\x1b[0m", int(fg.R*255), int(fg.G*255), int(fg.B*255)))
			} else {
				// Both exist
				buf.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m",
					int(fg.R*255), int(fg.G*255), int(fg.B*255),
					int(bg.R*255), int(bg.G*255), int(bg.B*255)))
			}
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func imageToSixel(img image.Image, config ConversionConfig) string {
	bounds := img.Bounds()
	ratio := float64(bounds.Dy()) / float64(bounds.Dx())
	// For Sixel we can use a slightly higher resolution as it's pixel-based
	targetWidth := config.Width * 8
	targetHeight := int(float64(targetWidth) * ratio)

	newRect := image.Rect(0, 0, targetWidth, targetHeight)
	resized := image.NewRGBA(newRect)
	draw.ApproxBiLinear.Scale(resized, newRect, img, bounds, draw.Over, nil)

	var buf bytes.Buffer
	enc := sixel.NewEncoder(&buf)
	err := enc.Encode(resized)
	if err != nil {
		return ""
	}

	return buf.String()
}
