package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// AnalyzeText analyses text
func AnalyzeText(text string) (int, int, int) {
	paragraphs := len(strings.Split(text, "\n\n"))
	words := len(strings.Fields(text))
	characters := len(text)
	return paragraphs, words, characters
}

// GetFileContentFromStoring gets file content from storage
func GetFileContentFromStoring(fileID string) (string, error) {
	resp, err := http.Get("http://file_storing:8080/files/" + fileID)
	if err != nil {
		log.Printf("file_storing is unreachable: %v\n", err)
		return "", fmt.Errorf("file_storing service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("file_storing returned non-200 status: %d\n", resp.StatusCode)
		return "", fmt.Errorf("file_storing returned error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return "", fmt.Errorf("failed to read response from file_storing")
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Failed to parse JSON response: %v\n", err)
		return "", fmt.Errorf("invalid JSON response")
	}

	content, ok := result["content"].(string)
	if !ok {
		log.Printf("Content field missing or not a string")
		return "", fmt.Errorf("invalid content format")
	}

	return content, nil
}
