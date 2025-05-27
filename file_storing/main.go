package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
)

func runMigrations() {
	_, err := db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS files (
            id UUID PRIMARY KEY,
            name TEXT NOT NULL,
            location TEXT NOT NULL
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
	e.POST("/upload", UploadFile)
	e.GET("/files/:id", GetFile)

	e.Start(":8080")
}
