package service

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

type PageText struct {
	PageNumber int32
	Text       string
}

func ExtractPDFText(content []byte) ([]PageText, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("open pdf reader: %w", err)
	}

	numPages := reader.NumPage()
	pages := make([]PageText, 0, numPages)

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return nil, fmt.Errorf("extract text from page %d: %w", i, err)
		}

		pages = append(pages, PageText{PageNumber: int32(i), Text: text})
	}

	return pages, nil
}
