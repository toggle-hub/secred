package utils

import (
	"context"
	"errors"
)

func GetUserFromContext(ctx context.Context) (string, error) {
	id := ctx.Value("user")
	if id == nil {
		return "", errors.New("authorization required")
	}

	val := id.(string)

	return val, nil
}
