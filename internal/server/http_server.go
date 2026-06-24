package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/knowledge-search-system/document-processor/config"
	pkggrpc "github.com/knowledge-search-system/document-processor/pkg/grpc"
	documentprocessorv1 "github.com/knowledge-search-system/document-processor/proto/documentprocessor/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newHTTPHandler(cfg *config.Config, logger *zap.Logger) (http.Handler, error) {
	endpoint := fmt.Sprintf("localhost:%d", cfg.GRPC.Port)

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial grpc endpoint %q: %w", endpoint, err)
	}

	gwMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(header string) (string, bool) {
			if strings.EqualFold(header, pkggrpc.LangMetadataKey) {
				return pkggrpc.LangMetadataKey, true
			}
			return runtime.DefaultHeaderMatcher(header)
		}),
	)
	if err := documentprocessorv1.RegisterDocumentServiceHandler(context.Background(), gwMux, conn); err != nil {
		return nil, fmt.Errorf("register grpc-gateway handler: %w", err)
	}

	documentClient := documentprocessorv1.NewDocumentServiceClient(conn)

	rootMux := http.NewServeMux()
	rootMux.Handle("/docs", newSwaggerHandler())
	rootMux.Handle("/docs/", newSwaggerHandler())
	rootMux.Handle("/api/v1/documents/upload", newUploadHandler(documentClient, logger))
	rootMux.Handle("/", gwMux)

	return rootMux, nil
}

func registerHTTPLifecycle(lc fx.Lifecycle, httpHandler http.Handler, cfg *config.Config, logger *zap.Logger) {
	srv := &http.Server{Handler: httpHandler}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("listen http on %s: %w", addr, err)
			}

			go func() {
				logger.Info("http server started", zap.String("addr", addr))
				if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
					logger.Error("http server stopped", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
