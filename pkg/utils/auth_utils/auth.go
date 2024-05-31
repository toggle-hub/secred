package authutils

import (
	"context"
	"database/sql"
	"errors"

	usermodel "github.com/xsadia/secred/pkg/models/user_model"
)

var ErrAuthorizationRequired = errors.New("authorization required")

func AuthenticateUser(ctx context.Context, db *sql.DB) (*usermodel.User, error) {
	id := ctx.Value("user")
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
