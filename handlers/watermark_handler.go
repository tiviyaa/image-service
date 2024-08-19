package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/go-programming/utils"
)

func WatermarkHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/watermark/")
	if id == "" {
		http.Error(w, "Missing image ID", http.StatusBadRequest)
		return
	}

	watermarkPaths, err := utils.WatermarkImageByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image watermarked successfully",
		"id":      id,
		"paths":   watermarkPaths,
	})
}
