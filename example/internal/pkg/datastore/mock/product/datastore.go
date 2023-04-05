package product

import (
	"context"

	"github.com/jpfaria/tests/example/internal/pkg/domain/model"
	"github.com/jpfaria/tests/example/internal/pkg/domain/repository/product"
)

type dataStore struct {
}

func (d dataStore) GetAll(ctx context.Context) ([]model.Product, error) {
	return []model.Product{
		{
			ID:    "ABC",
			Name:  "Ovo de Pascoa",
			Price: 34.0,
		},
	}, nil
}

func New() product.Repository {
	return &dataStore{}
}
