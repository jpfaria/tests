package echo

import (
	"sync"

	echofx "github.com/americanas-go/ignite/go.uber.org/fx.v1/module/labstack/echo.v4"
	e "github.com/americanas-go/ignite/labstack/echo.v4"
	"github.com/americanas-go/ignite/labstack/echo.v4/plugins/contrib/americanas-go/log.v1"
	status "github.com/americanas-go/ignite/labstack/echo.v4/plugins/contrib/americanas-go/rest-response.v1"
	json "github.com/americanas-go/ignite/labstack/echo.v4/plugins/contrib/bytedance/sonic.v1"
	datadog "github.com/americanas-go/ignite/labstack/echo.v4/plugins/contrib/datadog/dd-trace-go.v1"
	pprof "github.com/americanas-go/ignite/labstack/echo.v4/plugins/contrib/hiko1129/echo-pprof.v1"
	prometheus "github.com/americanas-go/ignite/labstack/echo.v4/plugins/contrib/prometheus/client_golang.v1"
	"github.com/americanas-go/ignite/labstack/echo.v4/plugins/extra/error_handler"
	"github.com/americanas-go/ignite/labstack/echo.v4/plugins/native/cors"
	"github.com/americanas-go/ignite/labstack/echo.v4/plugins/native/recover"
	"github.com/americanas-go/ignite/labstack/echo.v4/plugins/native/requestid"
	"go.uber.org/fx"
)

var once sync.Once

func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Module("echo",
			fx.Provide(
				func() []e.Plugin {
					return []e.Plugin{
						cors.Register,
						recover.Register,
						requestid.Register,
						log.Register,
						status.Register,
						prometheus.Register,
						pprof.Register,
						json.Register,
						error_handler.Register,
						datadog.Register,
					}
				},
			),
			echofx.Module(),
		)
	})

	return options
}
