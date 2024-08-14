package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"example.com/go-programming/utils"
)

type UploadRequest struct {
	ImageBase64 string `json:"imageBase64"`
	Filename    string `json:"filename"`
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var uploadRequest UploadRequest
	err := json.NewDecoder(r.Body).Decode(&uploadRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode base64 image
	decodedImage, err := base64.StdEncoding.DecodeString(uploadRequest.ImageBase64)
	if err != nil {
		http.Error(w, "Invalid base64 data", http.StatusBadRequest)
		return
	}

	// Save the original file
	originalFilePath := fmt.Sprintf("uploads/original-%d-%s", time.Now().Unix(), uploadRequest.Filename)
	err = ioutil.WriteFile(originalFilePath, decodedImage, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Path to the watermark image
	watermarkPath := "resources/watermark.png"

	// Initialize paths map
	paths := map[string]string{
		"originalPath": originalFilePath,
	}

	// Process resized images and add watermark
	sizes := map[string]string{
		"small":  "small",
		"medium": "medium",
		"large":  "large",
	}

	for size, sizeName := range sizes {
		resizedFilePath := fmt.Sprintf("uploads/%s-%d-%s", sizeName, time.Now().Unix(), uploadRequest.Filename)
		watermarkedFilePath := fmt.Sprintf("uploads/%s-watermarked-%d-%s", sizeName, time.Now().Unix(), uploadRequest.Filename)

		err = utils.ResizeAndAddWatermark(originalFilePath, resizedFilePath, watermarkedFilePath, watermarkPath, sizeName)
		if err != nil {
			http.Error(w, fmt.Sprintf("error processing size %s: %v", sizeName, err), http.StatusInternalServerError)
			return
		}

		// Upload resized and watermarked files to Firebase
		utils.UploadToFirebase(resizedFilePath)
		utils.UploadToFirebase(watermarkedFilePath)

		// Store the paths of resized and watermarked images
		paths[size+"Path"] = resizedFilePath
		paths[size+"WatermarkPath"] = watermarkedFilePath
	}

	// Store all paths in Firestore
	utils.StorePathsInFirestore(paths)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully Uploaded Files",
	})
}
