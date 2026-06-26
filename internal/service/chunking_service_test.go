package service

import (
	"strings"
	"testing"
)

func TestChunkText(t *testing.T) {
	t.Run("empty text yields no chunks", func(t *testing.T) {
		if chunks := ChunkText("", 1000, 100); chunks != nil {
			t.Errorf("expected nil, got %v", chunks)
		}
	})

	t.Run("text shorter than chunk size yields a single chunk", func(t *testing.T) {
		text := strings.Repeat("a", 500)
		chunks := ChunkText(text, 1000, 100)
		if len(chunks) != 1 || chunks[0] != text {
			t.Fatalf("expected a single chunk equal to input, got %v", chunks)
		}
	})

	t.Run("respects chunk size and overlap for cyrillic text", func(t *testing.T) {
		text := strings.Repeat("абвгд", 500) // 2500 runes
		chunks := ChunkText(text, 1000, 100)

		if len(chunks) != 3 {
			t.Fatalf("expected 3 chunks for 2500 runes with size=1000/overlap=100, got %d", len(chunks))
		}
		for i, c := range chunks[:2] {
			if got := len([]rune(c)); got != 1000 {
				t.Errorf("chunk %d: expected length 1000, got %d", i, got)
			}
		}

		// verify the configured overlap actually repeats between adjacent chunks
		firstRunes := []rune(chunks[0])
		secondRunes := []rune(chunks[1])
		overlapFromFirst := string(firstRunes[len(firstRunes)-100:])
		overlapFromSecond := string(secondRunes[:100])
		if overlapFromFirst != overlapFromSecond {
			t.Errorf("expected 100-rune overlap between chunk 0 and 1, got %q vs %q", overlapFromFirst, overlapFromSecond)
		}
	})

	t.Run("zero or negative chunk size yields no chunks", func(t *testing.T) {
		if chunks := ChunkText("some text", 0, 100); chunks != nil {
			t.Errorf("expected nil for zero chunk size, got %v", chunks)
		}
	})
}

func TestBuildChunksFromPages(t *testing.T) {
	pages := []PageText{
		{PageNumber: 1, Text: strings.Repeat("a", 1500)},
		{PageNumber: 2, Text: strings.Repeat("b", 500)},
	}

	chunks := BuildChunksFromPages("doc-1", "file.pdf", pages, 1000, 100)

	wantPageNumbers := []int32{1, 1, 2}
	if len(chunks) != len(wantPageNumbers) {
		t.Fatalf("expected %d chunks, got %d", len(wantPageNumbers), len(chunks))
	}
	for i, want := range wantPageNumbers {
		if chunks[i].PageNumber != want {
			t.Errorf("chunk %d: page number = %d, want %d", i, chunks[i].PageNumber, want)
		}
		if chunks[i].DocumentID != "doc-1" || chunks[i].FileName != "file.pdf" {
			t.Errorf("chunk %d: unexpected document_id/file_name metadata: %+v", i, chunks[i])
		}
	}

	seenIDs := map[string]struct{}{}
	for _, c := range chunks {
		if _, dup := seenIDs[c.ChunkID]; dup {
			t.Errorf("duplicate chunk_id %q", c.ChunkID)
		}
		seenIDs[c.ChunkID] = struct{}{}
	}
}

func TestBuildChunksFromText(t *testing.T) {
	chunks := BuildChunksFromText("doc-2", "file.docx", strings.Repeat("x", 1500), 1000, 100)

	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	for _, c := range chunks {
		if c.PageNumber != 1 {
			t.Errorf("expected page_number=1 for docx (no real pagination), got %d", c.PageNumber)
		}
	}
}
