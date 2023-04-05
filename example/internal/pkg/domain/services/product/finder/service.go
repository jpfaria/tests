//go:generate mockery --name Service --case underscore
package finder

import (
	"context"

	"github.com/jpfaria/tests/example/internal/pkg/domain/model"
)

type Service interface {
	Finder(ctx context.Context) ([]model.Product, error)
}