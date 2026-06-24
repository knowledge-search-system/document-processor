package model

import (
	"context"
	"time"
)

type DocumentStatus string

const (
	DocumentStatusUploading DocumentStatus = "uploading"
	DocumentStatusIndexing  DocumentStatus = "indexing"
	DocumentStatusReady     DocumentStatus = "ready"
	DocumentStatusFailed    DocumentStatus = "failed"
)

type Document struct {
	ID           string
	FileName     string
	Status       DocumentStatus
	UploadedAt   time.Time
	ErrorMessage string
}

type Chunk struct {
	ChunkID    string
	DocumentID string
	FileName   string
	PageNumber int32
	Text       string
}

type DocumentRepository interface {
	Create(ctx context.Context, doc Document) error
	UpdateStatus(ctx context.Context, id string, status DocumentStatus, errorMessage string) error
	Get(ctx context.Context, id string) (Document, error)
	List(ctx context.Context, page, pageSize int32) ([]Document, int32, error)
}
