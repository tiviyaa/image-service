package utils

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
	storageClient   *storage.Client
	bucketName      = "internshiptask-431606.appspot.com"
)

func init() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("C:/Users/user/Documents/Go Programming/serviceAccountKey.json")

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: "internshiptask-431606"}, sa)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		return
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		fmt.Printf("error initializing firestore: %v\n", err)
		return
	}

	storageClient, err = app.Storage(ctx)
	if err != nil {
		fmt.Printf("error initializing storage: %v\n", err)
		return
	}
}

func UploadToFirebase(filePath string) {
	ctx := context.Background()
	bucket, err := storageClient.Bucket(bucketName)
	if err != nil {
		fmt.Printf("error getting bucket: %v\n", err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer file.Close()

	wc := bucket.Object(filePath).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		fmt.Printf("error uploading file: %v\n", err)
		return
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("error closing writer: %v\n", err)
		return
	}
}

func StorePathsInFirestore(paths map[string]string) {
	ctx := context.Background()
	_, _, err := firestoreClient.Collection("images").Add(ctx, map[string]interface{}{
		"originalPath":        paths["originalPath"],
		"smallPath":           paths["smallPath"],
		"mediumPath":          paths["mediumPath"],
		"largePath":           paths["largePath"],
		"smallWatermarkPath":  paths["smallWatermarkPath"],
		"mediumWatermarkPath": paths["mediumWatermarkPath"],
		"largeWatermarkPath":  paths["largeWatermarkPath"],
		"timestamp":           firestore.ServerTimestamp,
	})
	if err != nil {
		fmt.Printf("error storing paths in firestore: %v\n", err)
	}
}
