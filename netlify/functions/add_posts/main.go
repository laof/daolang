package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Blog struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	dbName := os.Getenv("TURSO_DATABASE_NAME")
	dbAuthToken := os.Getenv("TURSO_DATABASE_TOKEN")
	url := fmt.Sprintf("%s?authToken=%s", dbName, dbAuthToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	defer db.Close()
	blogPayload := Blog{}
	json.Unmarshal([]byte(request.Body), &blogPayload)
	title := blogPayload.Title
	description := blogPayload.Description
	content := blogPayload.Content

	_, err = db.Exec("INSERT INTO posts (title, description, content) VALUES (?, ?, ?)", title, description, content)
	if err != nil {
		fmt.Println("Error inserting row:", err)
		return &events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Added post",
	}, nil
}

func main() {
	lambda.Start(handler)
}
