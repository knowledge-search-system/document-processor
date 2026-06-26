package service

import (
	"strings"
	"testing"
)

func TestExtractPDFText(t *testing.T) {
	pages, err := ExtractPDFText(readFixture(t, "sample.pdf"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("expected at least one page")
	}
	if !strings.Contains(pages[0].Text, "Машинное обучение") {
		t.Errorf("expected extracted text to contain the source phrase, got %q", pages[0].Text)
	}
}

func TestExtractPDFText_InvalidContent(t *testing.T) {
	if _, err := ExtractPDFText([]byte("not a real pdf")); err == nil {
		t.Fatal("expected an error for invalid pdf content")
	}
}

func TestExtractDOCXText(t *testing.T) {
	text, err := ExtractDOCXText(readFixture(t, "sample.docx"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(text, "Нейронные сети") {
		t.Errorf("expected extracted text to contain the source phrase, got %q", text)
	}
}

func TestExtractDOCXText_InvalidContent(t *testing.T) {
	if _, err := ExtractDOCXText([]byte("not a real docx")); err == nil {
		t.Fatal("expected an error for invalid docx content")
	}
}
