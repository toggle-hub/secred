package usermodel

import (
	"database/sql"

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
	user := model.User{
		Name:  input.Name,
		Email: input.Email,
	}
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	row := um.db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) returning id, created_at, updated_at",
		input.Name, input.Email, hash)

	err = utils.ParseDuplicateError(row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt), "email already in use")
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func New(db *sql.DB) *UserModel {
	return &UserModel{
		db,
	}
}
