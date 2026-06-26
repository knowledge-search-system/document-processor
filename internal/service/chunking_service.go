package service

import (
	"fmt"

	"github.com/knowledge-search-system/document-processor/internal/model"
)

func ChunkText(text string, chunkSize, overlap int) []string {
	if chunkSize <= 0 {
		return nil
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	step := chunkSize - overlap
	if step <= 0 {
		step = chunkSize
	}

	var chunks []string
	for start := 0; start < len(runes); start += step {
		end := min(start+chunkSize, len(runes))

		chunks = append(chunks, string(runes[start:end]))

		if end == len(runes) {
			break
		}
	}

	return chunks
}

func BuildChunksFromPages(documentID, fileName string, pages []PageText, chunkSize, overlap int) []model.Chunk {
	var chunks []model.Chunk
	chunkIndex := 0

	for _, page := range pages {
		for _, text := range ChunkText(page.Text, chunkSize, overlap) {
			chunks = append(chunks, model.Chunk{
				ChunkID:    fmt.Sprintf("%s-%d", documentID, chunkIndex),
				DocumentID: documentID,
				FileName:   fileName,
				PageNumber: page.PageNumber,
				Text:       text,
			})
			chunkIndex++
		}
	}

	return chunks
}

func BuildChunksFromText(documentID, fileName, text string, chunkSize, overlap int) []model.Chunk {
	var chunks []model.Chunk

	for i, t := range ChunkText(text, chunkSize, overlap) {
		chunks = append(chunks, model.Chunk{
			ChunkID:    fmt.Sprintf("%s-%d", documentID, i),
			DocumentID: documentID,
			FileName:   fileName,
			PageNumber: 1,
			Text:       t,
		})
	}

	return chunks
}
