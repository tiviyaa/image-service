package utils

import (
	"context"
	"fmt"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
	storageClient   *storage.Client
	bucketName      = "internshiptask-431606.appspot.com"
)

func init() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccountKey.json")

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

// Original Image - Store the Path in Firestore
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

// Update the Image Path (URL) in Firestore
func UpdateImagePathsInFirestore(docID string, paths map[string]string) error {
	ctx := context.Background()
	_, err := firestoreClient.Collection("images").Doc(docID).Set(ctx, paths, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error updating paths in Firestore: %w", err)
	}
	return nil
}

// Retrieve the Image Path by ID from the Cloud Firestore (images)
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

// Get data: Image from Postman
func GetImagePathsByID(docID string) (map[string]string, error) {
	ctx := context.Background()
	doc, err := firestoreClient.Collection("images").Doc(docID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving document from Firestore: %w", err)
	}

	paths := make(map[string]string)
	for key, value := range doc.Data() {
		if path, ok := value.(string); ok {
			paths[key] = path
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("no paths found for document ID %s", docID)
	}

	return paths, nil
}

// Upload the Image to the Firestore Storage
func UploadToFirebaseFromBytes(fileBytes []byte, filePath string) (string, error) {
	ctx := context.Background()
	bucket, err := storageClient.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("error getting bucket: %w", err)
	}

	filename := filepath.Base(filePath)
	wc := bucket.Object(filename).NewWriter(ctx)
	token := uuid.New().String()

	wc.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": token,
	}

	if _, err := wc.Write(fileBytes); err != nil {
		return "", fmt.Errorf("error writing file: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error closing writer: %w", err)
	}

	// Generate the full URL including the token
	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", bucketName, filename, token)

	return url, nil
}
