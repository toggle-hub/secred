package authutils

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"
	usermodel "github.com/xsadia/secred/pkg/models/user_model"
)

var ErrAuthorizationRequired = errors.New("authorization required")

func AuthenticateUser(c echo.Context, db *sql.DB) (*usermodel.User, error) {

	id := c.Get("user")
	if id == nil {
		return nil, ErrAuthorizationRequired
	}

	val := id.(string)
	userModel := usermodel.New(db)
	user, err := userModel.FindById(val)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuthorizationRequired
		}

		return nil, errors.New("unexpected error")
	}

	return user, nil
}
