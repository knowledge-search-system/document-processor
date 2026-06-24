package apperrors

var (
	ErrInvalidFileFormat = New(CodeInvalidArgument, "upload.invalid_file_format")
	ErrFileTooLarge      = New(CodeInvalidArgument, "upload.file_too_large")
	ErrEmptyFile         = New(CodeInvalidArgument, "upload.empty_file")
	ErrMissingFile       = New(CodeInvalidArgument, "upload.missing_file")
	ErrInvalidPagination = New(CodeInvalidArgument, "documents.invalid_pagination")
	ErrDocumentNotFound  = New(CodeNotFound, "documents.not_found")
	ErrExtractionFailed  = New(CodeInternal, "upload.extraction_failed")
	ErrIndexingFailed    = New(CodeUnavailable, "upload.indexing_failed")
	ErrInternal          = New(CodeInternal, "common.internal_error")
)
