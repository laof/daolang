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
	ID          int    `json:"id"`
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
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var blogs []Blog

	for rows.Next() {
		var blog Blog

		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Description, &blog.Content); err != nil {
			fmt.Println("Error scanning row:", err)
		}

		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}

	blogsJSON, err := json.Marshal(blogs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blogs: %v", err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(blogsJSON),
	}, nil
}

func main() {
	lambda.Start(handler)
}
