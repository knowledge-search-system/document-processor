package searchengine

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("searchengine_client",
	fx.Provide(NewClient),
	fx.Invoke(registerLifecycle),
)

func registerLifecycle(lc fx.Lifecycle, client *Client) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})
}
