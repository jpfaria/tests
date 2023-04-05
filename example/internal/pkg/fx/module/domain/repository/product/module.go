package product

import (
	"sync"

	rproduct "github.com/jpfaria/tests/example/internal/pkg/domain/repository/product"
	mproduct "github.com/jpfaria/tests/example/internal/pkg/fx/module/datastore/mock/product"
	"go.uber.org/fx"
)

var once sync.Once

func Module() fx.Option {

	options := fx.Options()

	once.Do(func() {

		switch rproduct.DataStore() {
		case "mongodb":
			// options = mproduct.Module()
		default:
			options = mproduct.Module()
		}

	})

	return options
}
