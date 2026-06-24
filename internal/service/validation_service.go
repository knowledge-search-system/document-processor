package service

import (
	"bytes"
	"strings"

	"github.com/knowledge-search-system/document-processor/config"
	"github.com/knowledge-search-system/document-processor/internal/apperrors"
)

type ValidationService struct {
	maxFileSizeBytes int64
}

func NewValidationService(cfg *config.Config) *ValidationService {
	return &ValidationService{maxFileSizeBytes: cfg.Upload.MaxFileSizeBytes}
}

func (s *ValidationService) Validate(fileName string, content []byte) error {
	if len(content) == 0 {
		return apperrors.ErrEmptyFile
	}

	if int64(len(content)) > s.maxFileSizeBytes {
		return apperrors.ErrFileTooLarge
	}

	if !hasValidSignature(fileName, content) {
		return apperrors.ErrInvalidFileFormat
	}

	return nil
}

func hasValidSignature(fileName string, content []byte) bool {
	lower := strings.ToLower(fileName)

	switch {
	case strings.HasSuffix(lower, ".pdf"):
		return bytes.HasPrefix(content, []byte("%PDF"))
	case strings.HasSuffix(lower, ".docx"):
		return bytes.HasPrefix(content, []byte("PK\x03\x04"))
	default:
		return false
	}
}
