package utils

import (
	"errors"

	"github.com/lib/pq"
)

func ParseDuplicateError(err error, message string) error {
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return errors.New(message)
			}
		}

		return errors.New("unexpected error")
	}

	return nil
}
