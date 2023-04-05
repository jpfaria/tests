package main

import (
	"github.com/americanas-go/ignite"
	"github.com/americanas-go/ignite/go.uber.org/fx.v1/module/americanas-go/multiserver.v1"
	datadogfx "github.com/americanas-go/ignite/go.uber.org/fx.v1/module/datadog/dd-trace-go.v1"
	chifx "github.com/americanas-go/ignite/go.uber.org/fx.v1/module/go-chi/chi.v5"
	"github.com/americanas-go/log"
	app1 "github.com/jpfaria/tests/example/internal/app/app1/fx/module/server/echo"
	"github.com/jpfaria/tests/example/internal/pkg/fx/module/lib/echo"

	ignitelog "github.com/americanas-go/ignite/americanas-go/log.v1"
	ignitefx "github.com/americanas-go/ignite/go.uber.org/fx.v1"
	c "github.com/americanas-go/ignite/spf13/cobra.v1"
	"github.com/spf13/cobra"
)

func main() {

	ignite.Boot()
	ignitelog.New()

	cmd := &cobra.Command{
		Version: "v1",
		Run: func(cmd *cobra.Command, args []string) {
			ignitefx.NewApp(
				datadogfx.TracerModule(),
				datadogfx.ProfilerModule(),
				chifx.Module(),
				echo.Module(),
				app1.Module(),
				multiserver.Module(),
			).Run()
		},
	}

	if err := c.Run(cmd); err != nil {
		log.Panicf(err.Error())
	}
}
