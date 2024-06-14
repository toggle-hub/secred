package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/database"
	schoolordermodel "github.com/xsadia/secred/pkg/models/school_orders"
	authutils "github.com/xsadia/secred/pkg/utils/auth_utils"
)

type SchoolPathParam struct {
	ID string `param:"id"`
}

func CreateSchoolOrderHandler(c echo.Context) error {
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

	params := new(SchoolPathParam)
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, HttpError{
			Message: "bad request",
		})
	}

	schoolOrderModel := schoolordermodel.New(db)
	schoolOrder, err := schoolOrderModel.Create(params.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	return c.JSON(http.StatusCreated, schoolOrder)
}
