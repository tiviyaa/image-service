package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func ResizeAndAddWatermark(srcPath, dstResizedPath, dstWatermarkedPath, watermarkPath, size string) error {
	// Open the original file
	file, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer file.Close()

	// Decode the original image
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	//Image Size (Small, Medium, Large) with WaterMark quantity by column and row
	var newWidth uint
	var rows, cols int
	switch size {
	case "small":
		newWidth = 256
		cols = 2
		rows = 2
	case "medium":
		newWidth = 512
		cols = 12
		rows = 10
	case "large":
		newWidth = 1024
		cols = 18
		rows = 15
	default:
		return fmt.Errorf("invalid size specified")
	}

	// Resize the image
	resizedImage := resize.Resize(newWidth, 0, img, resize.Lanczos3)

	// Add watermark to the resized image
	watermarkedImage, err := AddWatermark(resizedImage, watermarkPath, rows, cols)
	if err != nil {
		return fmt.Errorf("error adding watermark: %w", err)
	}

	// Save the resized image
	out, err := os.Create(dstResizedPath)
	if err != nil {
		return fmt.Errorf("error creating resized file: %w", err)
	}
	defer out.Close()

	if err := jpeg.Encode(out, resizedImage, nil); err != nil {
		return fmt.Errorf("error encoding resized image: %w", err)
	}

	// Save the watermarked image
	outWatermarked, err := os.Create(dstWatermarkedPath)
	if err != nil {
		return fmt.Errorf("error creating watermarked file: %w", err)
	}
	defer outWatermarked.Close()

	if err := jpeg.Encode(outWatermarked, watermarkedImage, nil); err != nil {
		return fmt.Errorf("error encoding watermarked image: %w", err)
	}

	return nil
}
