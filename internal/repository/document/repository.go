package document

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knowledge-search-system/document-processor/internal/apperrors"
	"github.com/knowledge-search-system/document-processor/internal/model"
	"github.com/knowledge-search-system/document-processor/internal/repository/sql"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) model.DocumentRepository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, doc model.Document) error {
	query, args, err := buildInsertQuery(doc)
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	executor := sql.ExecutorFromContext(ctx, r.pool)
	if _, err := executor.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert document: %w", err)
	}

	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, status model.DocumentStatus, errorMessage string) error {
	query, args, err := buildUpdateStatusQuery(id, status, errorMessage)
	if err != nil {
		return fmt.Errorf("build update status query: %w", err)
	}

	executor := sql.ExecutorFromContext(ctx, r.pool)
	if _, err := executor.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("update document status: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (model.Document, error) {
	query, args, err := buildGetQuery(id)
	if err != nil {
		return model.Document{}, fmt.Errorf("build get query: %w", err)
	}

	executor := sql.ExecutorFromContext(ctx, r.pool)
	row := executor.QueryRow(ctx, query, args...)

	doc, err := scanDocument(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Document{}, apperrors.ErrDocumentNotFound
		}
		return model.Document{}, fmt.Errorf("scan document: %w", err)
	}

	return doc, nil
}

func (r *Repository) List(ctx context.Context, page, pageSize int32) ([]model.Document, int32, error) {
	listQuery, listArgs, err := buildListQuery(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("build list query: %w", err)
	}

	executor := sql.ExecutorFromContext(ctx, r.pool)
	rows, err := executor.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()

	documents := make([]model.Document, 0, pageSize)
	for rows.Next() {
		doc, err := scanDocument(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan document row: %w", err)
		}
		documents = append(documents, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate document rows: %w", err)
	}

	countQuery, countArgs, err := buildCountQuery()
	if err != nil {
		return nil, 0, fmt.Errorf("build count query: %w", err)
	}

	var total int32
	if err := executor.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count documents: %w", err)
	}

	return documents, total, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanDocument(row scannable) (model.Document, error) {
	var doc model.Document
	var status string

	if err := row.Scan(&doc.ID, &doc.FileName, &status, &doc.UploadedAt, &doc.ErrorMessage); err != nil {
		return model.Document{}, err
	}

	doc.Status = model.DocumentStatus(status)

	return doc, nil
}
