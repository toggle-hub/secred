package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/database"
	itemmodel "github.com/xsadia/secred/pkg/models/items"
	authutils "github.com/xsadia/secred/pkg/utils/auth_utils"
)

type CreateItemRequestBody struct {
	Name     string `json:"name" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,gte=1"`
}

func CreateItemHandler(c echo.Context) error {
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

	body := new(CreateItemRequestBody)
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

	itemModel := itemmodel.New(db)
	foundItem, err := itemModel.FindByName(body.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			item, err := itemModel.Create(itemmodel.CreateItemInput{
				Name:     body.Name,
				Quantity: body.Quantity,
			})

			if err != nil {
				return c.JSON(http.StatusInternalServerError, HttpError{
					Message: "something went wrong",
				})
			}

			return c.JSON(http.StatusCreated, item)
		}

		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	updatedItem, err := itemModel.Update(foundItem.ID, foundItem.Quantity+body.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	return c.JSON(http.StatusCreated, updatedItem)
}

type ListItemsQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type ListItemsResponseBody struct {
	Data        []*itemmodel.Item `json:"data"`
	HasNextPage bool              `json:"hasNextPage"`
}

func ListItemsHandler(c echo.Context) error {
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

	queryParams := new(ListItemsQueryParams)
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

	itemModel := itemmodel.New(db)
	items, hasNextPage := itemModel.FindMany(page, limit)
	return c.JSON(http.StatusOK, ListItemsResponseBody{
		Data:        items,
		HasNextPage: hasNextPage,
	})
}
