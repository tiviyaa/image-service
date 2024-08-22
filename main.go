package main

import (
	"net/http"

	"example.com/go-programming/handlers"
)

func main() {

	//http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/resize/", handlers.ResizeHandler)
	http.HandleFunc("/watermark/", handlers.WatermarkHandler)
	http.HandleFunc("/image/", handlers.GetImageHandler)
	http.ListenAndServe(":8080", nil)
}
