package main

import (
	"github.com/knowledge-search-system/document-processor/config"
	"github.com/knowledge-search-system/document-processor/internal/apperrors"
	searchengineclient "github.com/knowledge-search-system/document-processor/internal/clients/searchengine"
	"github.com/knowledge-search-system/document-processor/internal/handler"
	"github.com/knowledge-search-system/document-processor/internal/logger"
	"github.com/knowledge-search-system/document-processor/internal/repository"
	"github.com/knowledge-search-system/document-processor/internal/repository/document"
	"github.com/knowledge-search-system/document-processor/internal/server"
	"github.com/knowledge-search-system/document-processor/internal/service"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		config.Module,
		logger.Module,
		apperrors.Module,
		repository.Module,
		document.Module,
		searchengineclient.Module,
		service.Module,
		handler.Module,
		server.Module,
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
