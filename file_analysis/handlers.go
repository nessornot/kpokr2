package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func AnalyzeFile(c echo.Context) error {
	fileID := c.Param("id")

	var result AnalysisResult
	err := db.QueryRow(context.Background(), `
        SELECT file_id, paragraphs, words, characters 
        FROM analysis_results 
        WHERE file_id = $1
    `, fileID).Scan(&result.FileID, &result.Paragraphs, &result.Words, &result.Characters)

	if err == nil {
		return c.JSON(http.StatusOK, result)
	}

	if err != pgx.ErrNoRows {
		log.Printf("DB error when fetching analysis: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal database error",
		})
	}

	content, err := GetFileContentFromStoring(fileID)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Could not fetch file content from file_storing",
		})
	}

	p, w, ch := AnalyzeText(content)

	_, err = db.Exec(context.Background(), `
        INSERT INTO analysis_results (file_id, paragraphs, words, characters) 
        VALUES ($1, $2, $3, $4)
    `, fileID, p, w, ch)

	if err != nil {
		log.Printf("DB error when saving analysis: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save analysis result",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"file_id":    fileID,
		"paragraphs": p,
		"words":      w,
		"characters": ch,
	})
}
