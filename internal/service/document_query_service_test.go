package service

import (
	"errors"
	"testing"

	"github.com/knowledge-search-system/document-processor/internal/apperrors"
)

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name         string
		page         int32
		pageSize     int32
		wantPage     int32
		wantPageSize int32
		wantErr      error
	}{
		{name: "defaults when zero", page: 0, pageSize: 0, wantPage: defaultPage, wantPageSize: defaultPageSize},
		{name: "explicit values kept", page: 3, pageSize: 20, wantPage: 3, wantPageSize: 20},
		{name: "negative page rejected", page: -1, pageSize: 10, wantErr: apperrors.ErrInvalidPagination},
		{name: "negative page size rejected", page: 1, pageSize: -10, wantErr: apperrors.ErrInvalidPagination},
		{name: "page size over max rejected", page: 1, pageSize: maxPageSize + 1, wantErr: apperrors.ErrInvalidPagination},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageSize, err := normalizePagination(tt.page, tt.pageSize)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotPage != tt.wantPage || gotPageSize != tt.wantPageSize {
				t.Errorf("got (%d, %d), want (%d, %d)", gotPage, gotPageSize, tt.wantPage, tt.wantPageSize)
			}
		})
	}
}
