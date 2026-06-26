package service

import (
	"errors"
	"os"
	"testing"

	"github.com/knowledge-search-system/document-processor/config"
	"github.com/knowledge-search-system/document-processor/internal/apperrors"
)

func newTestValidationService(maxSize int64) *ValidationService {
	return NewValidationService(&config.Config{Upload: config.UploadConfig{MaxFileSizeBytes: maxSize}})
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("../../tests/fixtures/" + name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return content
}

func TestValidationService_Validate(t *testing.T) {
	s := newTestValidationService(20 * 1024 * 1024)

	t.Run("valid pdf passes", func(t *testing.T) {
		if err := s.Validate("sample.pdf", readFixture(t, "sample.pdf")); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid docx passes", func(t *testing.T) {
		if err := s.Validate("sample.docx", readFixture(t, "sample.docx")); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty file rejected", func(t *testing.T) {
		err := s.Validate("empty.pdf", readFixture(t, "empty.pdf"))
		if !errors.Is(err, apperrors.ErrEmptyFile) {
			t.Errorf("expected ErrEmptyFile, got %v", err)
		}
	})

	t.Run("corrupted pdf signature rejected", func(t *testing.T) {
		err := s.Validate("corrupted.pdf", readFixture(t, "corrupted.pdf"))
		if !errors.Is(err, apperrors.ErrInvalidFileFormat) {
			t.Errorf("expected ErrInvalidFileFormat, got %v", err)
		}
	})

	t.Run("unsupported extension rejected", func(t *testing.T) {
		err := s.Validate("notes.txt", []byte("hello world"))
		if !errors.Is(err, apperrors.ErrInvalidFileFormat) {
			t.Errorf("expected ErrInvalidFileFormat, got %v", err)
		}
	})

	t.Run("file over size limit rejected", func(t *testing.T) {
		tiny := newTestValidationService(10)
		err := tiny.Validate("sample.pdf", readFixture(t, "sample.pdf"))
		if !errors.Is(err, apperrors.ErrFileTooLarge) {
			t.Errorf("expected ErrFileTooLarge, got %v", err)
		}
	})
}
