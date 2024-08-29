package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"example.com/go-programming/utils"
	"github.com/gin-gonic/gin"
)

type UploadRequest struct {
	ImageBase64 string `json:"imageBase64"`
	Filename    string `json:"filename"`
}

func UploadHandler(c *gin.Context) {
	var uploadRequest UploadRequest
	if err := c.ShouldBindJSON(&uploadRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Decode the base64 image
	decodedImage, err := base64.StdEncoding.DecodeString(uploadRequest.ImageBase64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 data"})
		return
	}

	// Store the path in Firestore to get the document ID
	docID, err := utils.StoreOriginalPathInFirestore("temporary")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store path in Firestore"})
		return
	}

	// Use the document ID as part of the file name
	originalFileName := uploadRequest.Filename
	filePath := fmt.Sprintf("uploaded-%s-%s", docID, originalFileName)

	// Upload the image to Firebase Storage with the docID in the filename
	url, err := utils.UploadToFirebaseFromBytes(decodedImage, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Firebase Storage"})
		return
	}

	// Update the Firestore record with the actual URL
	err = utils.UpdateImagePathsInFirestore(docID, map[string]string{"originalPath": url})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Firestore with image URL"})
		return
	}

	// Respond with the document ID and URL
	c.JSON(http.StatusOK, gin.H{
		"id":      docID,
		"message": "Successfully uploaded the image",
		"path":    url,
	})
}
