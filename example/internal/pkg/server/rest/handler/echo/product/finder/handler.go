package finder

import (
	"github.com/americanas-go/errors"
	"github.com/americanas-go/log"
	"github.com/jpfaria/tests/example/internal/pkg/domain/service/product/finder"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	productFinder finder.Service
}

func New(productFinder finder.Service) *Handler {
	return &Handler{
		productFinder,
	}
}

func (h *Handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	logger := log.FromContext(ctx)
	logger.Infof("recebendo um request")

	_, err := h.productFinder.FindAll(ctx)
	if err != nil {
		return err
	}

	return errors.Internalf("teste de erro")

	// return c.JSON(http.StatusOK, products)
}
