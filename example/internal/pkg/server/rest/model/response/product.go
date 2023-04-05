package model

import "github.com/jpfaria/tests/example/internal/pkg/domain/model"

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (m *Product) FromDomainModel(p model.Product) Product {
	return Product{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}
}
