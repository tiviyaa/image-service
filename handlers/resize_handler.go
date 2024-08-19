package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/go-programming/utils"
)

func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/resize/")
	if id == "" {
		http.Error(w, "Missing image ID", http.StatusBadRequest)
		return
	}

	resizedPaths, err := utils.ResizeImageByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image resized successfully",
		"id":      id,
		"paths":   resizedPaths,
	})
}
