package finder

import (
	"sync"

	sfinder "github.com/jpfaria/tests/example/internal/pkg/domain/service/product/finder"
	stdfinder "github.com/jpfaria/tests/example/internal/pkg/fx/module/provider/domain/service/product/finder"
	"go.uber.org/fx"
)

var once sync.Once

func Module() fx.Option {

	options := fx.Options()

	once.Do(func() {

		switch sfinder.Provider() {
		case "xpto":
			// options = xpto.Module()
		default:
			options = stdfinder.StandardModule()
		}

	})

	return options
}
