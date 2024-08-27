package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/nfnt/resize"
)

func WatermarkImageByID(id string) (map[string]string, error) {
	// Retrieve the original path from Firestore
	originalURL, err := GetImagePathByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving original image path: %w", err)
	}

	// Fetch the image from the originalPath (Firebase URL) and decode it
	resp, err := http.Get(originalURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading image from Firebase Storage: %w", err)
	}
	defer resp.Body.Close()

	// Decode the image from the response body
	img, _, err := image.Decode(resp.Body)
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
	watermarkPath := "resources/watermark.png"

	for size, settings := range sizes {
		// Resize the original image to the desired width while maintaining aspect ratio
		resizedImg := resize.Resize(settings.width, 0, img, resize.Lanczos3)

		// Apply watermark to the resized image
		watermarkedImage, err := AddWatermark(resizedImg, watermarkPath, settings.rows, settings.cols)
		if err != nil {
			return nil, fmt.Errorf("error applying watermark: %w", err)
		}

		// Encode the watermarked image to bytes
		var imageBytes []byte
		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, watermarkedImage, nil); err != nil {
			return nil, fmt.Errorf("error encoding watermarked image: %w", err)
		}
		imageBytes = buffer.Bytes()

		// Upload the watermarked image directly to Firebase Storage
		filePath := fmt.Sprintf("watermarked-%s-%s.jpg", size, id)
		uploadPath, err := UploadToFirebaseFromBytes(imageBytes, filePath)
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