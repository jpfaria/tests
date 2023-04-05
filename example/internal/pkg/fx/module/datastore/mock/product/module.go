package product

import (
	"github.com/jpfaria/tests/example/internal/pkg/datastore/mock/product"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("datastore_mock_product",
		fx.Provide(product.New),
	)
}
