package server

import (
	"context"
	"fmt"
	"net"

	"github.com/knowledge-search-system/document-processor/config"
	"github.com/knowledge-search-system/document-processor/internal/apperrors"
	"github.com/knowledge-search-system/document-processor/internal/handler"
	pkggrpc "github.com/knowledge-search-system/document-processor/pkg/grpc"
	documentprocessorv1 "github.com/knowledge-search-system/document-processor/proto/documentprocessor/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newGRPCServer(documentHandler *handler.DocumentHandler, logger *zap.Logger, translator *apperrors.Translator) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			pkggrpc.RecoveryInterceptor(logger),
			pkggrpc.LoggingInterceptor(logger),
			pkggrpc.ErrorTranslationInterceptor(translator),
		),
	)

	documentprocessorv1.RegisterDocumentServiceServer(srv, documentHandler)
	reflection.Register(srv)

	return srv
}

func registerGRPCLifecycle(lc fx.Lifecycle, srv *grpc.Server, cfg *config.Config, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.GRPC.Port)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("listen grpc on %s: %w", addr, err)
			}

			go func() {
				logger.Info("grpc server started", zap.String("addr", addr))
				if err := srv.Serve(lis); err != nil {
					logger.Error("grpc server stopped", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})
}
