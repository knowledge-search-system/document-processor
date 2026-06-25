package document

import (
	"testing"
	"time"

	"github.com/knowledge-search-system/document-processor/internal/model"
)

func TestBuildInsertQuery(t *testing.T) {
	doc := model.Document{
		ID:         "11111111-1111-1111-1111-111111111111",
		FileName:   "lecture.pdf",
		Status:     model.DocumentStatusUploading,
		UploadedAt: time.Date(2026, 6, 25, 10, 0, 0, 0, time.UTC),
	}

	query, args, err := buildInsertQuery(doc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantQuery := `INSERT INTO documents (id,file_name,status,uploaded_at,error_message) VALUES ($1,$2,$3,$4,$5)`
	if query != wantQuery {
		t.Errorf("query = %q, want %q", query, wantQuery)
	}
	if len(args) != 5 {
		t.Fatalf("expected 5 args, got %d", len(args))
	}
	if args[0] != doc.ID || args[1] != doc.FileName || args[2] != string(doc.Status) {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestBuildUpdateStatusQuery(t *testing.T) {
	query, args, err := buildUpdateStatusQuery("doc-1", model.DocumentStatusFailed, "boom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantQuery := `UPDATE documents SET status = $1, error_message = $2 WHERE id = $3`
	if query != wantQuery {
		t.Errorf("query = %q, want %q", query, wantQuery)
	}
	if len(args) != 3 || args[0] != "failed" || args[1] != "boom" || args[2] != "doc-1" {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestBuildListQuery(t *testing.T) {
	query, args, err := buildListQuery(2, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantQuery := `SELECT id, file_name, status, uploaded_at, error_message FROM documents ORDER BY uploaded_at DESC LIMIT 10 OFFSET 10`
	if query != wantQuery {
		t.Errorf("query = %q, want %q", query, wantQuery)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}
