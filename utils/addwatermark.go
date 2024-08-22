package utils

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

// AddWatermark applies a watermark to an image in a grid pattern
func AddWatermark(img image.Image, watermarkPath string, rows, cols int) (image.Image, error) {
	// Open the watermark file
	watermarkFile, err := os.Open(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("error opening watermark file: %w", err)
	}
	defer watermarkFile.Close()

	// Decode the watermark image
	watermarkImg, err := png.Decode(watermarkFile)
	if err != nil {
		return nil, fmt.Errorf("error decoding watermark image: %w", err)
	}

	bounds := img.Bounds()

	// Calculate the new size for the watermark based on the resized image dimensions
	newWatermarkWidth := bounds.Dx() / cols
	newWatermarkHeight := bounds.Dy() / rows

	// Resize the watermark image
	resizedWatermark := resize.Resize(uint(newWatermarkWidth), uint(newWatermarkHeight), watermarkImg, resize.Lanczos3)

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)

	// Calculate the spacing for the watermarks
	spacingX := bounds.Dx() / cols
	spacingY := bounds.Dy() / rows

	// Draw the watermark in a grid pattern on the resized image
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			pt := image.Point{
				X: x*spacingX + (spacingX-newWatermarkWidth)/2,
				Y: y*spacingY + (spacingY-newWatermarkHeight)/2,
			}
			draw.Draw(rgba, resizedWatermark.Bounds().Add(pt), resizedWatermark, image.Point{}, draw.Over)
		}
	}

	return rgba, nil
}
