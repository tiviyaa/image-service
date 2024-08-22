package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func WatermarkImageByID(id string) (map[string]string, error) {
	// Retrieve the original path from Firestore
	originalPath, err := GetImagePathByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving original image path: %w", err)
	}

	// Open the original image file
	file, err := os.Open(originalPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Decode the original image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	// Watermark settings based on image size
	watermarkPaths := make(map[string]string)
	sizes := map[string]struct {
		width uint
		rows  int
		cols  int
	}{
		"smallWatermark":  {width: 256, rows: 2, cols: 2},
		"mediumWatermark": {width: 512, rows: 10, cols: 12},
		"largeWatermark":  {width: 1024, rows: 15, cols: 18},
	}

	// Relative path to the watermark image
	watermarkPath := filepath.Join("resources", "watermark.png")

	for size, settings := range sizes {
		// Resize the original image to the desired width while maintaining aspect ratio
		resizedImg := resize.Resize(settings.width, 0, img, resize.Lanczos3)

		// Apply watermark to the resized image
		watermarkedImage, err := AddWatermark(resizedImg, watermarkPath, settings.rows, settings.cols)
		if err != nil {
			return nil, fmt.Errorf("error applying watermark: %w", err)
		}

		// Save the watermarked image
		dstWatermarkedPath := filepath.Join("uploads", fmt.Sprintf("watermarked-%s-%s.jpg", size, id))
		out, err := os.Create(dstWatermarkedPath)
		if err != nil {
			return nil, fmt.Errorf("error creating watermarked file: %w", err)
		}
		defer out.Close()

		if err := jpeg.Encode(out, watermarkedImage, nil); err != nil {
			return nil, fmt.Errorf("error encoding watermarked image: %w", err)
		}

		// Upload the watermarked image to Firebase Storage and store the path
		uploadPath, err := UploadToFirebase(dstWatermarkedPath)
		if err != nil {
			return nil, fmt.Errorf("error uploading watermarked image to Firebase Storage: %w", err)
		}

		// Store the path in the map
		watermarkPaths[size] = uploadPath

		// Update Firestore with the watermarked image path
		if err := UpdateImagePathsInFirestore(id, map[string]string{size: uploadPath}); err != nil {
			return nil, fmt.Errorf("error updating Firestore with watermarked paths: %w", err)
		}
	}

	return watermarkPaths, nil
}
