package finder

import (
	findersvc "github.com/jpfaria/tests/example/internal/pkg/fx/module/domain/service/product/finder"
	"github.com/jpfaria/tests/example/internal/pkg/server/rest/handler/echo/product/finder"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		findersvc.Module(),
		fx.Provide(
			finder.New,
		),
		fx.Invoke(
			finder.Router,
		),
	)
}
