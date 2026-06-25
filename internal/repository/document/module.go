package document

import "go.uber.org/fx"

var Module = fx.Module("document_repository",
	fx.Provide(NewRepository),
)
