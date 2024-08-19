package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	sa := option.WithCredentialsFile("C:/Users/user/Pictures/Go Programming/serviceAccountKey.json")

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

func StoreOriginalPathInFirestore(originalPath string) (string, error) {
	ctx := context.Background()
	docRef, _, err := firestoreClient.Collection("images").Add(ctx, map[string]interface{}{
		"originalPath": originalPath,
		"timestamp":    firestore.ServerTimestamp,
	})
	if err != nil {
		return "", fmt.Errorf("error storing path in Firestore: %w", err)
	}
	return docRef.ID, nil
}

func UpdateImagePathsInFirestore(docID string, paths map[string]string) error {
	ctx := context.Background()
	_, err := firestoreClient.Collection("images").Doc(docID).Set(ctx, paths, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error updating paths in Firestore: %w", err)
	}
	return nil
}

func GetImagePathByID(docID string) (string, error) {
	ctx := context.Background()
	doc, err := firestoreClient.Collection("images").Doc(docID).Get(ctx)
	if err != nil {
		return "", fmt.Errorf("error retrieving document from Firestore: %w", err)
	}

	originalPath, ok := doc.Data()["originalPath"].(string)
	if !ok {
		return "", fmt.Errorf("invalid document format")
	}

	return originalPath, nil
}

func UploadToFirebase(filePath string) (string, error) {
	ctx := context.Background()
	bucket, err := storageClient.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("error getting bucket: %w", err)
	}

	filename := filepath.Base(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	wc := bucket.Object(filename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("error uploading file: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	return filename, nil
}
