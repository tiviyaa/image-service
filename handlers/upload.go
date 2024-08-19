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

	// Store the original file path in Firestore and get the ID
	docID, err := utils.StoreOriginalPathInFirestore(originalFilePath)
	if err != nil {
		http.Error(w, "Failed to store path in Firestore", http.StatusInternalServerError)
		return
	}

	// Respond with the document ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      docID,
		"message": "Successfully uploaded the image",
	})
}
