package server

import (
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	documentprocessorv1 "github.com/knowledge-search-system/document-processor/proto/documentprocessor/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxUploadMemoryBytes = 32 << 20

func newUploadHandler(client documentprocessorv1.DocumentServiceClient, logger *zap.Logger) http.Handler {
	errorMux := runtime.NewServeMux()
	marshaler := &runtime.JSONPb{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			renderError(errorMux, marshaler, w, r, status.Error(codes.InvalidArgument, "method not allowed, use POST"))
			return
		}

		if err := r.ParseMultipartForm(maxUploadMemoryBytes); err != nil {
			renderError(errorMux, marshaler, w, r, status.Error(codes.InvalidArgument, "invalid multipart/form-data request"))
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			renderError(errorMux, marshaler, w, r, status.Error(codes.InvalidArgument, "missing \"file\" field in form data"))
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			renderError(errorMux, marshaler, w, r, status.Error(codes.Internal, "failed to read uploaded file"))
			return
		}

		resp, err := client.UploadDocument(r.Context(), &documentprocessorv1.UploadDocumentRequest{
			FileName: header.Filename,
			Content:  content,
		})
		if err != nil {
			renderError(errorMux, marshaler, w, r, err)
			return
		}

		w.Header().Set("Content-Type", marshaler.ContentType(resp))
		w.WriteHeader(http.StatusOK)
		if err := marshaler.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to encode upload response", zap.Error(err))
		}
	})
}

func renderError(mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	runtime.DefaultHTTPErrorHandler(r.Context(), mux, marshaler, w, r, err)
}
