package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/database"
	ordermodel "github.com/xsadia/secred/pkg/models/orders"
	authutils "github.com/xsadia/secred/pkg/utils/auth_utils"
)

type CreateOrderRequestBody struct {
	InvoiceUrl *string `json:"invoiceUrl,omitempty" validate:"omitempty,url"`
}

func CreateOrderHandler(c echo.Context) error {
	storage, err := database.GetInstance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	db := storage.DB()
	_, err = authutils.AuthenticateUser(c, db)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, HttpError{
			Message: err.Error(),
		})
	}

	body := new(CreateOrderRequestBody)
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, HttpError{
			Message: "bad request",
		})
	}

	orderModel := ordermodel.New(db)
	order, err := orderModel.Create(body.InvoiceUrl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	return c.JSON(http.StatusCreated, order)
}
