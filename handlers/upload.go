package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

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

	// Use the filename provided in the request
	originalFileName := uploadRequest.Filename

	// Save the image directly to Firebase Storage
	filePath := fmt.Sprintf("uploaded-%s", originalFileName)
	url, err := utils.UploadToFirebaseFromBytes(decodedImage, filePath)
	if err != nil {
		http.Error(w, "Failed to upload image to Firebase Storage", http.StatusInternalServerError)
		return
	}

	// Store the file path in Firestore and get the document ID
	docID, err := utils.StoreOriginalPathInFirestore(url)
	if err != nil {
		http.Error(w, "Failed to store path in Firestore", http.StatusInternalServerError)
		return
	}

	// Respond with the document ID and Firebase Storage URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      docID,
		"message": "Successfully uploaded the image",
		"path":    url,

	})
}
