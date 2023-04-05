package finder

import (
	"github.com/jpfaria/tests/example/internal/pkg/fx/module/domain/repository/product"
	"github.com/jpfaria/tests/example/internal/pkg/provider/domain/service/product/finder"
	"go.uber.org/fx"
)

func StandardModule() fx.Option {
	return fx.Module("provider_domain_service_product_finder_standard",
		product.Module(),
		fx.Provide(finder.NewStandard),
	)
}
