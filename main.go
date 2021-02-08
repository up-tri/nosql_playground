package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

func loadDotEnv() {
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		log.Fatal(err)
	}
}

func makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error")
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

func createClient(ctx context.Context) *firestore.Client {
	projectID := os.Getenv("GCP_PROJECT_ID")

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func addUser(ctx context.Context, client *firestore.Client, userName string) error {
	_, _, err := client.Collection("Users").Add(ctx, map[string]interface{}{
		"name": userName,
	})
	if err != nil {
		return err
	}

	return nil
}

func printAllUser(ctx context.Context, client *firestore.Client) error {
	iter := client.Collection("Users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(doc.Data()["name"])
	}

	return nil
}

func main() {
	loadDotEnv()

	ctx := context.Background()
	client := createClient(ctx)

	randName, err := makeRandomStr(10)
	if err != nil {
		log.Fatal(err)
	}

	addUser(ctx, client, randName)
	printAllUser(ctx, client)
}
