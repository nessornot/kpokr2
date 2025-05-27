package main

import (
	"bytes"
	echoSwagger "github.com/swaggo/echo-swagger"
	"io"
	"log"
	"net/http"

	_ "kpokr2/api_gateway/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// UploadHandler
// @Summary Upload a text file
// @Description Uploads a text file and returns a unique file ID
// @Accept mpfd
// @Produce json
// @Param file formData file true "Text file to upload"
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /upload [post]
func UploadHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed to read request body: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read request body",
		})
	}

	// Создаем новый запрос
	req, _ := http.NewRequest("POST", "http://file_storing:8080/upload", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", c.Request().Header.Get("Content-Type")) // Передаем Content-Type

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("file_storing is down: %v\n", err)
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "File storing service is down",
		})
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	return c.HTMLBlob(resp.StatusCode, responseBody)
}

// AnalyzeHandler
// @Summary Analyze file by ID
// @Description Analyzes a previously uploaded text file and returns statistics
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /analyze/{id} [get]
func AnalyzeHandler(c echo.Context) error {
	id := c.Param("id")
	resp, err := http.Get("http://file_analysis:8081/analyze/" + id)
	if err != nil {
		log.Printf("file_analysis is down: %v\n", err)
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Analysis service is down",
		})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.JSONBlob(resp.StatusCode, body)
}

// DownloadHandler
// @Summary Get file content by ID
// @Description Retrieves the content of a previously uploaded text file
// @Produce text/plain
// @Param id path string true "File ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /files/{id} [get]
func DownloadHandler(c echo.Context) error {
	id := c.Param("id")
	resp, err := http.Get("http://file_storing:8080/files/" + id)
	if err != nil {
		log.Printf("file_storing is down: %v\n", err)
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "File storing service is down",
		})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.JSONBlob(resp.StatusCode, body)
}

// @title kpokr2 API

// @BasePath /
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// file upload
	e.POST("/upload", UploadHandler)

	// file analysis
	e.GET("/analyze/:id", AnalyzeHandler)

	// file download
	e.GET("/files/:id", DownloadHandler)

	// swagger ui
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	e.Logger.Fatal(e.Start(":8080"))
}
