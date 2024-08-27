package handlers

import (
	"net/http"
	"strings"

	"example.com/go-programming/utils"
)

func GetImageHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/image/")
	if id == "" {
		http.Error(w, "Missing image ID", http.StatusBadRequest)
		return
	}

	// Fetch the image URL from Firestore
	imageURL, err := utils.GetImagePathByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to the image URL in Firebase Storage
	http.Redirect(w, r, imageURL, http.StatusFound)
}
