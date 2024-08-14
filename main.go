package main

import (
	"net/http"

	"example.com/go-programming/handlers"
)

func main() {
	// Serve static files from the uploads directory
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Handle image uploads
	http.HandleFunc("/upload", handlers.UploadHandler)

	http.ListenAndServe(":8080", nil)
}
