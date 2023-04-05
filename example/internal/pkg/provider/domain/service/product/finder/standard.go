package finder

import (
	"context"

	"github.com/jpfaria/tests/example/internal/pkg/domain/model"
	"github.com/jpfaria/tests/example/internal/pkg/domain/repository/product"
	"github.com/jpfaria/tests/example/internal/pkg/domain/services/product/finder"
)

type standard struct {
	products product.Repository
}

func (s *standard) Finder(ctx context.Context) ([]model.Product, error) {
	// TODO: implement me
	return s.products.GetAll(ctx)
}

func NewStandard(products product.Repository) finder.Service {
	return &standard{products: products}
}
