package schoolorderitemmodel

import (
	"database/sql"
	"fmt"
	"time"
)

type SchoolOrderItemModel struct {
	db *sql.DB
}

type SchoolOrderItem struct {
	ID        string
	ItemID    string
	OrderID   string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type SchoolOrderItemCreateInput struct {
	ItemID   string
	OrderID  string
	Quantity int
}

func (sm *SchoolOrderItemModel) CreateMany(input []SchoolOrderItemCreateInput) ([]*SchoolOrderItem, error) {
	values := make([]interface{}, len(input)*3)
	query := "INSERT INTO school_order_items (item_id, order_id, quantity) VALUES "

	for i, entry := range input {
		query = query + fmt.Sprintf("($%d, $%d, $%d),", 3*i+1, 3*i+2, 3*i+3)
		values[3*i] = entry.ItemID
		values[3*i+1] = entry.OrderID
		values[3*i+2] = entry.Quantity
	}

	query = query[:len(query)-1] + "RETURNING id, order_id, item_id, quantity, created_at, updated_at, deleted_at"

	rows, err := sm.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	var results []*SchoolOrderItem
	for rows.Next() {
		var item SchoolOrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ItemID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, &item)
	}

	return results, nil
}

func New(db *sql.DB) *SchoolOrderItemModel {
	return &SchoolOrderItemModel{
		db,
	}
}
