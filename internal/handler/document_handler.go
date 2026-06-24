package handler

import (
	"context"

	documentprocessorv1 "github.com/knowledge-search-system/document-processor/proto/documentprocessor/v1"
	"github.com/knowledge-search-system/document-processor/internal/model"
	"github.com/knowledge-search-system/document-processor/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DocumentHandler struct {
	documentprocessorv1.UnimplementedDocumentServiceServer

	uploadService *service.UploadService
	queryService  *service.DocumentQueryService
}

func NewDocumentHandler(uploadService *service.UploadService, queryService *service.DocumentQueryService) *DocumentHandler {
	return &DocumentHandler{
		uploadService: uploadService,
		queryService:  queryService,
	}
}

func (h *DocumentHandler) ListDocuments(ctx context.Context, req *documentprocessorv1.ListDocumentsRequest) (*documentprocessorv1.ListDocumentsResponse, error) {
	documents, total, err := h.queryService.List(ctx, req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}

	resp := &documentprocessorv1.ListDocumentsResponse{
		Documents: make([]*documentprocessorv1.Document, 0, len(documents)),
		Total:     total,
	}
	for _, doc := range documents {
		resp.Documents = append(resp.Documents, toProtoDocument(doc))
	}

	return resp, nil
}

func (h *DocumentHandler) GetDocument(ctx context.Context, req *documentprocessorv1.GetDocumentRequest) (*documentprocessorv1.Document, error) {
	doc, err := h.queryService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return toProtoDocument(doc), nil
}

func (h *DocumentHandler) UploadDocument(ctx context.Context, req *documentprocessorv1.UploadDocumentRequest) (*documentprocessorv1.UploadDocumentResponse, error) {
	doc, err := h.uploadService.Upload(ctx, req.GetFileName(), req.GetContent())
	if err != nil {
		return nil, err
	}

	return &documentprocessorv1.UploadDocumentResponse{
		Id:     doc.ID,
		Status: string(doc.Status),
	}, nil
}

func toProtoDocument(doc model.Document) *documentprocessorv1.Document {
	return &documentprocessorv1.Document{
		Id:           doc.ID,
		FileName:     doc.FileName,
		Status:       string(doc.Status),
		UploadedAt:   timestamppb.New(doc.UploadedAt),
		ErrorMessage: doc.ErrorMessage,
	}
}
