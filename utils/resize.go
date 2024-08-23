package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/nfnt/resize"
)

func ResizeImageByID(id string) (map[string]string, error) {
	// Retrieve the original path (URL) from Firestore
	originalURL, err := GetImagePathByID(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving original image path: %w", err)
	}

	// Download the image from Firebase Storage using the URL
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

	// Create paths for resized images and store them in Firestore
	resizedPaths := make(map[string]string)
	sizes := map[string]uint{
		"small":  256,
		"medium": 512,
		"large":  1024,
	}

	for size, newWidth := range sizes {
		// Resize the image in-memory
		resizedImage := resize.Resize(newWidth, 0, img, resize.Lanczos3)

		// Encode the resized image into a byte buffer
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, resizedImage, nil); err != nil {
			return nil, fmt.Errorf("error encoding resized image: %w", err)
		}

		// Upload the resized image to Firebase Storage
		filePath := fmt.Sprintf("resized-%s-%s.jpg", size, id)
		uploadPath, err := UploadToFirebaseFromBytes(buf.Bytes(), filePath)
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
