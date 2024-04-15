package usermodel

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/xsadia/secred/graph/model"
	"github.com/xsadia/secred/pkg/utils"
)

type UserModel struct {
	db *sql.DB
}

func (um *UserModel) FindById(id string) (*model.User, error) {
	var user model.User
	err := um.db.QueryRow(
		"SELECT id, name, email, created_at, updated_at, deleted_at FROM users where id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (um *UserModel) Create(input model.CreateUserInput) (*model.User, error) {
	var id uuid.UUID
	var createdAt, updatedAt time.Time

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	err = um.db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) returning id, created_at, updated_at",
		input.Name, input.Email, hash).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return nil, errors.New("email already in use")
			}
		}

		return nil, errors.New("unexpected error")
	}

	return &model.User{
		ID:        id.String(),
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: nil,
	}, nil
}

func New(db *sql.DB) *UserModel {
	return &UserModel{
		db,
	}
}
