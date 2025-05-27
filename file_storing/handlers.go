package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"os"
)

func UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No file provided",
		})
	}

	id := uuid.New()
	dst := "./storage/" + id.String() + ".txt"

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	outFile, err := os.Create(dst)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create local file",
		})
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save file content",
		})
	}

	_, err = db.Exec(context.Background(), `
        INSERT INTO files (id, name, location) VALUES ($1, $2, $3)
    `, id, file.Filename, dst)

	if err != nil {
		log.Printf("Database error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to store file metadata",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func GetFile(c echo.Context) error {
	id := c.Param("id")
	var location string

	err := db.QueryRow(context.Background(), `
        SELECT location FROM files WHERE id = $1
    `, id).Scan(&location)

	if err == pgx.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "File not found in DB",
		})
	}

	if err != nil {
		log.Printf("DB query error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}

	content, err := os.ReadFile(location)
	if err != nil {
		log.Printf("Failed to read file: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read file from disk",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"content": string(content),
	})
}
