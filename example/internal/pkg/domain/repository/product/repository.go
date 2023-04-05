//go:generate mockery --name Repository --case underscore
package product

import (
	"context"

	"github.com/jpfaria/tests/example/internal/pkg/domain/model"
)

type Repository interface {
	GetAll(ctx context.Context) ([]model.Product, error)
}
