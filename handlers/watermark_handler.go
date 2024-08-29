package handlers

import (
	"net/http"

	"example.com/go-programming/utils"
	"github.com/gin-gonic/gin"
)

func WatermarkHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing image ID"})
		return
	}

	watermarkPaths, err := utils.WatermarkImageByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image watermarked successfully",
		"id":      id,
		"paths":   watermarkPaths,
	})
}
