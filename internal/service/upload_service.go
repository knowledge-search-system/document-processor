package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/knowledge-search-system/document-processor/config"
	"github.com/knowledge-search-system/document-processor/internal/apperrors"
	"github.com/knowledge-search-system/document-processor/internal/clients/searchengine"
	"github.com/knowledge-search-system/document-processor/internal/model"
	"go.uber.org/zap"
)

type UploadService struct {
	validation   *ValidationService
	repo         model.DocumentRepository
	searchClient *searchengine.Client
	chunkSize    int
	chunkOverlap int
	logger       *zap.Logger
}

func NewUploadService(
	validation *ValidationService,
	repo model.DocumentRepository,
	searchClient *searchengine.Client,
	cfg *config.Config,
	logger *zap.Logger,
) *UploadService {
	return &UploadService{
		validation:   validation,
		repo:         repo,
		searchClient: searchClient,
		chunkSize:    cfg.Upload.ChunkSize,
		chunkOverlap: cfg.Upload.ChunkOverlap,
		logger:       logger,
	}
}

func (s *UploadService) Upload(ctx context.Context, fileName string, content []byte) (model.Document, error) {
	if err := s.validation.Validate(fileName, content); err != nil {
		return model.Document{}, err
	}

	doc := model.Document{
		ID:         uuid.New().String(),
		FileName:   fileName,
		Status:     model.DocumentStatusUploading,
		UploadedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return model.Document{}, apperrors.ErrInternal.WithErr(fmt.Errorf("create document record: %w", err))
	}

	chunks, err := s.extractAndChunk(doc.ID, fileName, content)
	if err != nil {
		s.markFailed(ctx, doc.ID, err)
		return model.Document{}, err
	}

	if err := s.repo.UpdateStatus(ctx, doc.ID, model.DocumentStatusIndexing, ""); err != nil {
		s.logger.Error("failed to update document status to indexing", zap.Error(err), zap.String("document_id", doc.ID))
	}

	if _, err := s.searchClient.IndexChunks(ctx, chunks); err != nil {
		wrapped := apperrors.ErrIndexingFailed.WithErr(err)
		s.markFailed(ctx, doc.ID, wrapped)
		return model.Document{}, wrapped
	}

	if err := s.repo.UpdateStatus(ctx, doc.ID, model.DocumentStatusReady, ""); err != nil {
		return model.Document{}, apperrors.ErrInternal.WithErr(fmt.Errorf("update document status to ready: %w", err))
	}

	doc.Status = model.DocumentStatusReady

	return doc, nil
}

func (s *UploadService) extractAndChunk(documentID, fileName string, content []byte) ([]model.Chunk, error) {
	switch strings.ToLower(filepath.Ext(fileName)) {
	case ".pdf":
		pages, err := ExtractPDFText(content)
		if err != nil {
			return nil, apperrors.ErrExtractionFailed.WithErr(err)
		}
		return BuildChunksFromPages(documentID, fileName, pages, s.chunkSize, s.chunkOverlap), nil
	case ".docx":
		text, err := ExtractDOCXText(content)
		if err != nil {
			return nil, apperrors.ErrExtractionFailed.WithErr(err)
		}
		return BuildChunksFromText(documentID, fileName, text, s.chunkSize, s.chunkOverlap), nil
	default:
		return nil, apperrors.ErrInvalidFileFormat
	}
}

func (s *UploadService) markFailed(ctx context.Context, documentID string, err error) {
	if updateErr := s.repo.UpdateStatus(ctx, documentID, model.DocumentStatusFailed, err.Error()); updateErr != nil {
		s.logger.Error("failed to mark document as failed", zap.Error(updateErr), zap.String("document_id", documentID))
	}
}
