// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/jpfaria/tests/example/internal/pkg/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: ctx
func (_m *Repository) GetAll(ctx context.Context) ([]model.Product, error) {
	ret := _m.Called(ctx)

	var r0 []model.Product
	if rf, ok := ret.Get(0).(func(context.Context) []model.Product); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Product)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
