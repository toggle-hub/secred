package itemmodel

import (
	"database/sql"
	"strings"

	"github.com/xsadia/secred/graph/model"
	"github.com/xsadia/secred/pkg/utils"
)

type ItemModel struct {
	db *sql.DB
}

func (im *ItemModel) Create(input model.CreateItemInput) (*model.Item, error) {
	rawName := strings.ToLower(input.Name)
	name := utils.RemoveDiacritics(rawName)
	item := model.Item{
		Name:     name,
		RawName:  rawName,
		Quantity: input.Quantity,
	}
	row := im.db.QueryRow(
		"INSERT INTO items (name, raw_name, quantity) VALUES ($1, $2, $3) returning id, created_at, updated_at",
		name, rawName, input.Quantity)

	err := utils.ParseDuplicateError(row.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt), "item with given name already registered")
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func New(db *sql.DB) *ItemModel {
	return &ItemModel{
		db,
	}
}
