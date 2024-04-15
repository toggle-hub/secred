package schoolmodel

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/xsadia/secred/graph/model"
)

type SchoolModel struct {
	db *sql.DB
}

func (sm *SchoolModel) Create(input model.CreateSchoolInput) (*model.School, error) {
	var id uuid.UUID
	var createdAt, updatedAt time.Time
	err := sm.db.QueryRow(
		"INSERT INTO schools (name, address, phone_number) VALUES ($1, $2, $3) returning id, created_at, updated_at",
		input.Name, input.Address, input.PhoneNumber).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return nil, errors.New("school name already in use")
			}
		}

		return nil, errors.New("unexpected error")
	}
	return &model.School{
		ID:          id.String(),
		Name:        input.Name,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Orders:      nil,
		DeletedAt:   nil,
	}, nil
}

func New(db *sql.DB) *SchoolModel {
	return &SchoolModel{
		db,
	}
}
