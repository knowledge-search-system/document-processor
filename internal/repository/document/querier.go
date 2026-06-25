package document

import (
	"github.com/Masterminds/squirrel"
	"github.com/knowledge-search-system/document-processor/internal/model"
	"github.com/knowledge-search-system/document-processor/internal/repository/dbconsts"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func buildInsertQuery(doc model.Document) (string, []any, error) {
	return psql.Insert(dbconsts.DocumentsTable).
		Columns("id", "file_name", "status", "uploaded_at", "error_message").
		Values(doc.ID, doc.FileName, string(doc.Status), doc.UploadedAt, doc.ErrorMessage).
		ToSql()
}

func buildUpdateStatusQuery(id string, status model.DocumentStatus, errorMessage string) (string, []any, error) {
	return psql.Update(dbconsts.DocumentsTable).
		Set("status", string(status)).
		Set("error_message", errorMessage).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func buildGetQuery(id string) (string, []any, error) {
	return psql.Select("id", "file_name", "status", "uploaded_at", "error_message").
		From(dbconsts.DocumentsTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
}

func buildListQuery(page, pageSize int32) (string, []any, error) {
	offset := (page - 1) * pageSize
	return psql.Select("id", "file_name", "status", "uploaded_at", "error_message").
		From(dbconsts.DocumentsTable).
		OrderBy("uploaded_at DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
}

func buildCountQuery() (string, []any, error) {
	return psql.Select("COUNT(*)").
		From(dbconsts.DocumentsTable).
		ToSql()
}
