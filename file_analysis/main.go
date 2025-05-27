package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
)

func runMigrations() {
	_, err := db.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS analysis_results (
        file_id UUID PRIMARY KEY,
        paragraphs INT NOT NULL,
        words INT NOT NULL,
        characters INT NOT NULL
    );
`)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}
}

func main() {
	initDB()
	runMigrations()

	e := echo.New()
	e.GET("/analyze/:id", AnalyzeFile)
	err := e.Start(":8081")
	if err != nil {
		log.Fatal(err)
	}
}
