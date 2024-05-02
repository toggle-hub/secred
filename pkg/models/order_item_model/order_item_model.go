package orderitemmodel

import (
	"database/sql"
	"fmt"
	"time"
)

type OrderItemModel struct {
	db *sql.DB
}

type OrderItem struct {
	ID        string
	ItemId    string
	OrderId   string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type OrderItemCreateInput struct {
	OrderId  string
	ItemId   string
	Quantity int
}

func (om *OrderItemModel) Create(input OrderItemCreateInput) (*OrderItem, error) {
	orderItem := OrderItem{
		OrderId:  input.OrderId,
		ItemId:   input.ItemId,
		Quantity: input.Quantity,
	}

	err := om.db.QueryRow(
		"INSERT INTO order_items (order_id, item_id, quantity) VALUES ($1, $2, $3) returning id, created_at, updated_at, deleted_at",
		orderItem.OrderId, orderItem.ItemId, orderItem.Quantity,
	).Scan(&orderItem.ID, &orderItem.CreatedAt, &orderItem.UpdatedAt, &orderItem.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}

func (om *OrderItemModel) CreateMany(input []OrderItemCreateInput) ([]*OrderItem, error) {
	query := "INSERT INTO order_items (order_id, item_id, quantity) VALUES"
	values := make([]interface{}, len(input)*3)
	for i, item := range input {
		query = query + fmt.Sprintf("($%d, $%d, $%d),", 3*i+1, 3*i+2, 3*i+3)

		values[3*i] = item.OrderId
		values[3*i+1] = item.ItemId
		values[3*i+2] = item.Quantity
	}

	query = query[:len(query)-1] + "returning id, order_id, item_id, quantity, created_at, updated_at, deleted_at"

	rows, err := om.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	var results []*OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderId,
			&item.ItemId,
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

func New(db *sql.DB) *OrderItemModel {
	return &OrderItemModel{
		db,
	}
}
