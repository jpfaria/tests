package finder

import (
	"context"

	"github.com/americanas-go/errors"

	"github.com/jpfaria/tests/example/internal/pkg/domain/model"
	"github.com/jpfaria/tests/example/internal/pkg/domain/repository/product"
	"github.com/jpfaria/tests/example/internal/pkg/domain/service/product/finder"
)

type standard struct {
	products product.Repository
}

func (s *standard) FindAll(ctx context.Context) ([]model.Product, error) {
	products, err := s.products.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, errors.NotFoundf("no records found")
	}
	return products, nil
}

func NewStandard(products product.Repository) finder.Service {
	return &standard{products: products}
}
