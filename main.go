package main

import (
	"example.com/go-programming/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) //Release Mode
	r := gin.Default()

	r.POST("/upload", handlers.UploadHandler)
	r.POST("/resize/:id", handlers.ResizeHandler)
	r.POST("/watermark/:id", handlers.WatermarkHandler)
	r.GET("/image/:id", handlers.GetImageHandler)

	r.Run(":8080")
}
