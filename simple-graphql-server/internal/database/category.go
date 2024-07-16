package database

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) (*Category, error) {
	err := createCategoryTable(db)
	if err != nil {
		return nil, err
	}

	return &Category{
		db: db,
	}, nil
}

func createCategoryTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS category (id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT)")
	if err != nil {
		return err
	}

	return nil
}

func (c *Category) Create(name string, description string) (Category, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Category{}, err
	}

	_, err = c.db.Exec("INSERT INTO category (id, name, description) VALUES ($1, $2, $3)", id, name, description)
	if err != nil {
		return Category{}, err
	}

	return Category{
		ID:          id.String(),
		Name:        name,
		Description: description,
	}, nil
}

func (c *Category) FindById(id string) (Category, error) {
	err := uuid.Validate(id)
	if err != nil {
		return Category{}, errors.New("id should be a valid UUID")
	}

	rows, err := c.db.Query("SELECT id, name, description FROM category WHERE id = ?", id)
	if err != nil {
		return Category{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var categoryCategory Category
		err = rows.Scan(
			&categoryCategory.ID,
			&categoryCategory.Name,
			&categoryCategory.Description,
		)

		if err != nil {
			return Category{}, errors.New("error trying to find category")
		}

		return categoryCategory, nil
	}

	return Category{}, errors.New("no category found")
}

func (c *Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query(`
		SELECT id, name, description from category;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)

	for rows.Next() {
		var category Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)
		if err != nil {
			continue
		}

		categories = append(categories, category)
	}

	return categories, nil
}
