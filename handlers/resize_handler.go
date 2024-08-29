package handlers

import (
	"net/http"

	"example.com/go-programming/utils"
	"github.com/gin-gonic/gin"
)

func ResizeHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image ID"})
		return
	}

	resizedPaths, err := utils.ResizeImageByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image resized successfully",
		"id":      id,
		"paths":   resizedPaths,
	})
}
