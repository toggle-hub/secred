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

func (im *ItemModel) LoadManyByName(input []string) ([]*model.Item, error) {
	query := "SELECT * FROM items WHERE name in ("
	values := make([]interface{}, len(input))
	for i, entry := range input {
		query = query + fmt.Sprintf("$%d,", i+1)
		values[i] = utils.RemoveDiacritics(entry)
	}
	query = query[:len(query)-1] + ")"

	rows, err := im.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	var results []*model.Item
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
			return nil, err
		}

		results = append(results, &item)
	}

	return results, nil
}

func (im *ItemModel) LoadMany(page, limit int) ([]*model.Item, bool) {
	var results []*model.Item
	offset := (page - 1) * limit
	rows, err := im.db.Query("SELECT * FROM items WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", limit+1, offset)
	if err != nil {
		log.Println("Error fetching items: ", err.Error())
		return nil, false
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

func (im *ItemModel) CreateMany(input []model.CreateItemInput) ([]*model.Item, error) {
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

	query = query[:len(query)-1] + " returning id, name, raw_name, quantity, created_at, updated_at"
	rows, err := im.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	var results []*model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.RawName,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, &item)
	}
	return results, nil
}

type UpdateItemInput struct {
	ID       string
	Name     string
	Quantity int
}

func (im *ItemModel) UpdateMany(input []UpdateItemInput) ([]*model.Item, error) {
	query := `UPDATE items 
	SET name = data.name,
	raw_name = data.raw_name,
	quantity = data.quantity
	FROM (
		VALUES 
	`
	values := make([]interface{}, len(input)*4)
	for i, entry := range input {
		rawName := strings.ToLower(entry.Name)
		name := utils.RemoveDiacritics(rawName)
		query = query + fmt.Sprintf("($%d::uuid, $%d, $%d, $%d::integer),", 4*i+1, 4*i+2, 4*i+3, 4*i+4)

		values[4*i] = entry.ID
		values[4*i+1] = name
		values[4*i+2] = rawName
		values[4*i+3] = entry.Quantity
	}
	query = query[:len(query)-1] + `) 
	AS data(id, name, raw_name, quantity)
	WHERE items.id = data.id RETURNING items.id, items.name, items.raw_name, items.quantity, items.created_at, items.updated_at`
	rows, err := im.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	var results []*model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.RawName,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, &item)
	}

	return results, nil
}

type LoadCreateItem struct {
	Inserted bool
	*model.Item
}

func (im *ItemModel) LoadOrCreateManyByName(input []model.CreateItemInput) ([]LoadCreateItem, error) {
	query := "SELECT id, name, raw_name, quantity, created_at, updated_at FROM items WHERE name IN ("
	values := make([]interface{}, len(input))
	for i, item := range input {
		query = query + fmt.Sprintf("$%d,", i+1)
		name := utils.RemoveDiacritics(strings.ToLower(item.Name))
		values[i] = name
	}
	query = query[:len(query)-1] + ")"

	rows, err := im.db.Query(query, values...)
	if err != nil {
		return nil, err
	}

	var results []LoadCreateItem
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.RawName,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		results = append(results, LoadCreateItem{
			Inserted: false,
			Item:     &item,
		})
	}

	if len(input) != len(results) {
		m := make(map[string]string, len(results))
		for _, result := range results {
			m[result.Name] = result.ID
		}

		var notFound []model.CreateItemInput

		for _, item := range input {
			if _, ok := m[utils.RemoveDiacritics(strings.ToLower(item.Name))]; !ok {
				notFound = append(notFound, model.CreateItemInput{
					Name:     item.Name,
					Quantity: item.Quantity,
				})
			}
		}

		insertedItems, err := im.CreateMany(notFound)
		if err != nil {
			return nil, err
		}

		inserted := make([]LoadCreateItem, len(insertedItems))
		for i, item := range insertedItems {
			inserted[i] = LoadCreateItem{
				Inserted: true,
				Item:     item,
			}
		}

		results = append(results, inserted...)
	}

	return results, nil
}

func New(db *sql.DB) *ItemModel {
	return &ItemModel{
		db,
	}
}
