package apperrors

import "testing"

func TestTranslatorTranslate(t *testing.T) {
	translator := NewTranslator()

	tests := []struct {
		name       string
		messageKey string
		lang       string
		want       string
	}{
		{name: "known key ru", messageKey: "upload.invalid_file_format", lang: "ru", want: "Недопустимый формат файла. Поддерживаются только PDF и DOCX"},
		{name: "known key en", messageKey: "upload.invalid_file_format", lang: "en", want: "Unsupported file format. Only PDF and DOCX are allowed"},
		{name: "unknown lang falls back to ru", messageKey: "documents.not_found", lang: "fr", want: "Документ не найден"},
		{name: "unknown key falls back to internal error", messageKey: "does.not.exist", lang: "ru", want: "Внутренняя ошибка сервера"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translator.Translate(tt.messageKey, tt.lang)
			if got != tt.want {
				t.Errorf("Translate(%q, %q) = %q, want %q", tt.messageKey, tt.lang, got, tt.want)
			}
		})
	}
}

func TestSentinelErrorCodeAndMessageKey(t *testing.T) {
	err := ErrDocumentNotFound

	if err.Code() != int(CodeNotFound) {
		t.Errorf("Code() = %d, want %d", err.Code(), CodeNotFound)
	}
	if err.MessageKey() != "documents.not_found" {
		t.Errorf("MessageKey() = %q, want %q", err.MessageKey(), "documents.not_found")
	}

	wrapped := err.WithErr(errBoom)
	if wrapped.Unwrap() != errBoom {
		t.Error("Unwrap() did not return the wrapped error")
	}
}

var errBoom = &boomErr{"boom"}

type boomErr struct{ msg string }

func (e *boomErr) Error() string { return e.msg }
