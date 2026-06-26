package service

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type docxText struct {
	Value string `xml:",chardata"`
}

type docxRun struct {
	Text []docxText `xml:"t"`
}

type docxParagraph struct {
	Runs []docxRun `xml:"r"`
}

type docxBody struct {
	Paragraphs []docxParagraph `xml:"p"`
}

type docxDocument struct {
	Body docxBody `xml:"body"`
}

const docxDocumentXMLPath = "word/document.xml"

func ExtractDOCXText(content []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("open docx as zip: %w", err)
	}

	data, err := readZipFile(zr, docxDocumentXMLPath)
	if err != nil {
		return "", err
	}

	var doc docxDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("unmarshal %s: %w", docxDocumentXMLPath, err)
	}

	var sb strings.Builder
	for _, p := range doc.Body.Paragraphs {
		for _, r := range p.Runs {
			for _, t := range r.Text {
				sb.WriteString(t.Value)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func readZipFile(zr *zip.Reader, name string) ([]byte, error) {
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", name, err)
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}

		return data, nil
	}

	return nil, fmt.Errorf("%s not found in docx archive", name)
}
