package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func ResizeImageByID(id string) (map[string]string, error) {
	// Retrieve the original path from Firestore
	originalPath, err := GetImagePathByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving original image path: %w", err)
	}

	file, err := os.Open(originalPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	// Create paths for resized images and store them in Firestore
	resizedPaths := make(map[string]string)
	sizes := map[string]uint{
		"small":  256,
		"medium": 512,
		"large":  1024,
	}

	for size, newWidth := range sizes {
		dstResizedPath := fmt.Sprintf("uploads/%s-%s.jpg", size, id)
		resizedImage := resize.Resize(newWidth, 0, img, resize.Lanczos3)

		out, err := os.Create(dstResizedPath)
		if err != nil {
			return nil, fmt.Errorf("error creating resized file: %w", err)
		}
		defer out.Close()

		// Encode the resized image as JPEG
		if err := jpeg.Encode(out, resizedImage, nil); err != nil {
			return nil, fmt.Errorf("error encoding resized image: %w", err)
		}

		// Upload the resized image to Firebase Storage
		uploadPath, err := UploadToFirebase(dstResizedPath)
		if err != nil {
			return nil, fmt.Errorf("error uploading resized image to Firebase Storage: %w", err)
		}

		// Store the path in the map
		resizedPaths[size] = uploadPath
	}

	// Update Firestore with the resized image paths
	if err := UpdateImagePathsInFirestore(id, resizedPaths); err != nil {
		return nil, fmt.Errorf("error updating Firestore with resized paths: %w", err)
	}

	return resizedPaths, nil
}
