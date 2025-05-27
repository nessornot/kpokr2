package main

import "github.com/google/uuid"

type AnalysisResult struct {
	FileID     uuid.UUID `json:"file_id"`
	Paragraphs int       `json:"paragraphs"`
	Words      int       `json:"words"`
	Characters int       `json:"characters"`
}
