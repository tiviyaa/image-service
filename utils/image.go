package utils

import (
	"fmt"
)

func GetOriginalImageURL(id string) string {
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/internshiptask-431606.appspot.com/o/uploaded-original-%s.jpg?alt=media", id)
}

func GetResizedImageURL(id string, size string) string {
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/internshiptask-431606.appspot.com/o/resized-%s-%s.jpg?alt=media", size, id)
}

func GetWatermarkedImageURL(id string, watermark string) string {
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/internshiptask-431606.appspot.com/o/watermarked-%s-%s.jpg?alt=media", watermark, id)
}
