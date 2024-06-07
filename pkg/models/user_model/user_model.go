package usermodel

import (
	"database/sql"
	"time"

	"github.com/xsadia/secred/pkg/utils"
)

const UserType = "User"

type UserModel struct {
	db *sql.DB
}

type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (um *UserModel) FindByEmail(email string) (*User, error) {
	var user User
	err := um.db.QueryRow(
		"SELECT id, name, email, password, created_at, updated_at, deleted_at FROM users where email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (um *UserModel) FindById(id string) (*User, error) {
	var user User
	err := um.db.QueryRow(
		"SELECT id, name, email, created_at, updated_at, deleted_at FROM users where id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (um *UserModel) Create(email, name, password string) (*User, error) {
	user := User{
		Name:  name,
		Email: email,
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	row := um.db.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) returning id, created_at, updated_at",
		name, email, hash)

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
