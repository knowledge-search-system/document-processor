-- +goose Up
CREATE TABLE documents (
    id            uuid PRIMARY KEY,
    file_name     text NOT NULL,
    status        text NOT NULL,
    uploaded_at   timestamptz NOT NULL DEFAULT now(),
    error_message text NOT NULL DEFAULT ''
);

CREATE INDEX idx_documents_uploaded_at ON documents (uploaded_at DESC);

-- +goose Down
DROP TABLE IF EXISTS documents;
