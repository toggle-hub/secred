package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/xsadia/secred/pkg/database"
	usermodel "github.com/xsadia/secred/pkg/models/user_model"
	"github.com/xsadia/secred/pkg/utils"
)

type AuthResponseBody struct {
	User  *usermodel.User `json:"user"`
	Token string          `json:"token"`
}

type RegisterRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,gte=8"`
}

func RegisterHandler(c echo.Context) error {
	body := new(RegisterRequestBody)
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

	storage, err := database.GetInstance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	userModel := usermodel.New(storage.DB())

	_, err = userModel.FindByEmail(body.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, HttpError{
			Message: "email in already in use",
		})
	}

	user, err := userModel.Create(body.Email, body.Name, body.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	token, err := utils.CreateJWT(user.ID, utils.JWTExpireTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	return c.JSON(http.StatusCreated, AuthResponseBody{
		User:  user,
		Token: token,
	})
}

type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func LoginHandler(c echo.Context) error {
	body := new(LoginRequestBody)
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, HttpError{
			Message: err.Error(),
		})
	}

	storage, err := database.GetInstance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	userModel := usermodel.New(storage.DB())

	user, err := userModel.FindByEmail(body.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, HttpError{
			Message: "wrong password or email",
		})
	}

	if err := utils.ComparePassword(body.Password, user.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, HttpError{
			Message: "wrong password or email",
		})
	}

	token, err := utils.CreateJWT(user.ID, utils.JWTExpireTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HttpError{
			Message: "something went wrong",
		})
	}

	return c.JSON(http.StatusOK, AuthResponseBody{
		User:  user,
		Token: token,
	})
}
