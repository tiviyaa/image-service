package handlers

import (
	"net/http"

	"example.com/go-programming/utils"
	"github.com/gin-gonic/gin"
)

func GetImageHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image ID"})
		return
	}

	// Query Params in Postman
	imageType := c.DefaultQuery("type", "original") // Default to "original"
	imageSize := c.DefaultQuery("size", "")

	// Retrieve image paths from Firestore
	paths, err := utils.GetImagePathsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var imageURL string
	switch imageType {
	case "resized":
		imageURL = paths[imageSize]
	case "watermarked":
		imageURL = paths[imageSize+"Watermark"]
	case "original":
		imageURL = paths["originalPath"]
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type"})
		return
	}

	if imageURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.Redirect(http.StatusFound, imageURL)
}
