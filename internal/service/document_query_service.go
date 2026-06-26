package service

import (
	"context"

	"github.com/knowledge-search-system/document-processor/internal/apperrors"
	"github.com/knowledge-search-system/document-processor/internal/model"
)

const (
	defaultPage     = int32(1)
	defaultPageSize = int32(10)
	maxPageSize     = int32(100)
)

type DocumentQueryService struct {
	repo model.DocumentRepository
}

func NewDocumentQueryService(repo model.DocumentRepository) *DocumentQueryService {
	return &DocumentQueryService{repo: repo}
}

func (s *DocumentQueryService) List(ctx context.Context, page, pageSize int32) ([]model.Document, int32, error) {
	page, pageSize, err := normalizePagination(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.List(ctx, page, pageSize)
}

func (s *DocumentQueryService) Get(ctx context.Context, id string) (model.Document, error) {
	return s.repo.Get(ctx, id)
}

func normalizePagination(page, pageSize int32) (int32, int32, error) {
	if page < 0 || pageSize < 0 {
		return 0, 0, apperrors.ErrInvalidPagination
	}

	if page == 0 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		return 0, 0, apperrors.ErrInvalidPagination
	}

	return page, pageSize, nil
}
