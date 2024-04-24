package itemmodel

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/xsadia/secred/graph/model"
	"github.com/xsadia/secred/pkg/utils"
)

type ItemModel struct {
	db *sql.DB
}

var EmptyItemsList = []*model.Item{}

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

func (im *ItemModel) LoadAll(page, limit int) ([]*model.Item, bool) {
	var results []*model.Item
	offset := (page - 1) * limit
	rows, err := im.db.Query("SELECT * FROM items WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", limit+1, offset)
	if err != nil {
		log.Println("Error fetching items: ", err.Error())
		return EmptyItemsList, false
	}

	for rows.Next() {
		var item model.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.RawName,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return EmptyItemsList, false
		}

		results = append(results, &item)
	}

	if len(results) < 1 {
		return EmptyItemsList, false
	}

	if len(results) > limit {
		return results[:limit], true
	}

	return results, false
}

func (im *ItemModel) CreateMany(input []model.CreateItemInput) error {
	query := "INSERT INTO items (name, raw_name, quantity) VALUES "
	values := make([]interface{}, len(input)*3)

	for i, item := range input {
		query = query + fmt.Sprintf("($%d, $%d, $%d),", 3*i+1, 3*i+2, 3*i+3)
		rawName := strings.ToLower(item.Name)
		name := utils.RemoveDiacritics(rawName)

		values[3*i] = name
		values[3*i+1] = rawName
		values[3*i+2] = item.Quantity
	}

	query = query[:len(query)-1]
	_, err := im.db.Exec(query, values...)
	return err
}

func New(db *sql.DB) *ItemModel {
	return &ItemModel{
		db,
	}
}
