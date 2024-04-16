package schoolmodel

import (
	"database/sql"

	"github.com/xsadia/secred/graph/model"
	"github.com/xsadia/secred/pkg/utils"
)

type SchoolModel struct {
	db *sql.DB
}

func (sm *SchoolModel) Create(input model.CreateSchoolInput) (*model.School, error) {
	school := model.School{
		Name:        input.Name,
		Address:     input.Address,
		PhoneNumber: input.PhoneNumber,
	}
	row := sm.db.QueryRow(
		"INSERT INTO schools (name, address, phone_number) VALUES ($1, $2, $3) returning id, created_at, updated_at",
		input.Name, input.Address, input.PhoneNumber)

	err := utils.ParseDuplicateError(row.Scan(&school.ID, &school.CreatedAt, &school.UpdatedAt), "school name already in use")
	if err != nil {
		return nil, err
	}

	return &school, nil
}

func New(db *sql.DB) *SchoolModel {
	return &SchoolModel{
		db,
	}
}
