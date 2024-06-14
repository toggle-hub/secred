package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/database"
	schoolmodel "github.com/xsadia/secred/pkg/models/schools"
	authutils "github.com/xsadia/secred/pkg/utils/auth_utils"
)

type CreateSchoolRequestBody struct {
	Name        string  `json:"name" validate:"required"`
	Address     *string `json:"address,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty" validate:"omitempty,numeric"`
}

func CreateSchoolHandler(c echo.Context) error {
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

	body := new(CreateSchoolRequestBody)

	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, HttpError{
			Message: "bad request",
		})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, HttpError{
			Message: err.Error(),
		})
	}

	schoolModel := schoolmodel.New(db)
	school, err := schoolModel.Create(body.Name, body.Address, body.PhoneNumber)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	return c.JSON(http.StatusCreated, school)
}

type ListSchoolQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type ListSchoolsResponseBody struct {
	Data        []*schoolmodel.School `json:"data"`
	HasNextPage bool                  `json:"hasNextPage"`
}

func ListSchoolHandler(c echo.Context) error {
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

	queryParams := new(ListSchoolQueryParams)
	if err := c.Bind(queryParams); err != nil {
		return c.JSON(http.StatusBadRequest, HttpError{
			Message: "bad request",
		})
	}

	page := queryParams.Page
	limit := queryParams.Limit
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	schoolModel := schoolmodel.New(db)
	schools, hasNextPage := schoolModel.LoadMany(page, limit)
	return c.JSON(http.StatusOK, ListSchoolsResponseBody{
		Data:        schools,
		HasNextPage: hasNextPage,
	})
}
