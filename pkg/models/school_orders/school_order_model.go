package schoolordermodel

import (
	"database/sql"
	"time"

	"github.com/xsadia/secred/graph/model"
)

type SchoolOrderModel struct {
	db *sql.DB
}

type SchoolOrder struct {
	ID          string
	SchoolID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeliveredAt *time.Time
	DeletedAt   *time.Time
}

func (sm *SchoolOrderModel) Create(schoolId string) (*SchoolOrder, error) {
	var order SchoolOrder
	err := sm.db.QueryRow(
		"INSERT INTO school_orders (school_id) values ($1) RETURNING id, school_id, created_at, updated_at, deleted_at",
		schoolId,
	).Scan(
		&order.ID,
		&order.SchoolID,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (sm *SchoolOrderModel) LoadMany(schoolId string) ([]*model.SchoolOrder, error) {
	rows, err := sm.db.Query("SELECT id, school_id, created_at, updated_at, deleted_at FROM school_orders WHERE school_id = $1", schoolId)
	if err != nil {
		return nil, err
	}

	var results []*model.SchoolOrder
	for rows.Next() {
		var order model.SchoolOrder
		if err := rows.Scan(
			&order.ID,
			&order.SchoolID,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, &order)
	}

	return results, nil
}

func New(db *sql.DB) *SchoolOrderModel {
	return &SchoolOrderModel{
		db,
	}
}
