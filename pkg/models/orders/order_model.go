package ordermodel

import (
	"database/sql"
	"time"
)

type OrderModel struct {
	db *sql.DB
}

type Order struct {
	ID          string
	InvoiceUrl  *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeliveredAt *time.Time
	DeletedAt   *time.Time
}

func (om *OrderModel) Create(invoiceUrl *string) (*Order, error) {
	var order Order
	err := om.db.QueryRow(
		"INSERT INTO orders (invoice_url) VALUES ($1) returning id, created_at, updated_at",
		invoiceUrl,
	).Scan(
		&order.ID, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func New(db *sql.DB) *OrderModel {
	return &OrderModel{
		db,
	}
}
