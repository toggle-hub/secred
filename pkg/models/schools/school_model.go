package schoolmodel

import (
	"database/sql"
	"log"
	"time"

	"github.com/xsadia/secred/pkg/utils"
)

type SchoolModel struct {
	db *sql.DB
}

type School struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Address     *string `json:"address,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	// Orders      []*Order   `json:"orders"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (sm *SchoolModel) FindById(id string) (*School, error) {
	var school School
	err := sm.db.QueryRow(
		"SELECT * FROM schools WHERE id = $1 AND deleted_at IS NULL",
		id,
	).Scan(
		&school.ID,
		&school.Name,
		&school.Address,
		&school.PhoneNumber,
		&school.CreatedAt,
		&school.UpdatedAt,
		&school.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &school, nil
}

func (sm *SchoolModel) Create(name string, address *string, phoneNumber *string) (*School, error) {
	school := School{
		Name:        name,
		Address:     address,
		PhoneNumber: phoneNumber,
	}
	row := sm.db.QueryRow(
		"INSERT INTO schools (name, address, phone_number) VALUES ($1, $2, $3) returning id, created_at, updated_at",
		name, address, phoneNumber)

	err := utils.ParseDuplicateError(row.Scan(&school.ID, &school.CreatedAt, &school.UpdatedAt), "school name already in use")
	if err != nil {
		return nil, err
	}

	return &school, nil
}

func (sm *SchoolModel) LoadMany(page, limit int) ([]*School, bool) {
	var results []*School
	offset := (page - 1) * limit
	rows, err := sm.db.Query("SELECT * FROM schools WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", limit+1, offset)
	if err != nil {
		log.Println("Error fetching items: ", err.Error())
		return nil, false
	}

	for rows.Next() {
		var item School
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Address,
			&item.PhoneNumber,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return nil, false
		}

		results = append(results, &item)
	}

	if len(results) < 1 {
		return nil, false
	}

	if len(results) > limit {
		return results[:limit], true
	}

	return results, false
}

// func (sm *SchoolModel) CreateMany(input []model.CreateSchoolInput) error {
// 	query := "INSERT INTO schools (name, address, phone_number) VALUES "
// 	values := make([]interface{}, len(input)*3)

// 	for i, school := range input {
// 		query = query + fmt.Sprintf("($%d, $%d, $%d),", 3*i+1, 3*i+2, 3*i+3)

// 		values[3*i] = school.Name
// 		values[3*i+1] = school.Address
// 		values[3*i+2] = school.PhoneNumber
// 	}

// 	query = query[:len(query)-1]
// 	_, err := sm.db.Exec(query, values...)
// 	return err
// }

func New(db *sql.DB) *SchoolModel {
	return &SchoolModel{
		db,
	}
}
