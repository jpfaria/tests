package echo

import (
	"github.com/jpfaria/tests/example/internal/pkg/fx/module/server/rest/handler/echo/product/finder"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		finder.Module(),
	)
}
