package database

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

func NewCourse(db *sql.DB) (*Course, error) {
	err := createCourseTable(db)
	if err != nil {
		return nil, err
	}

	return &Course{
		db: db,
	}, nil
}

func createCourseTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS 
    	course (
    	    id TEXT PRIMARY KEY, 
    	    name TEXT NOT NULL, 
    	    description TEXT, 
    	    categoryId TEXT,
    	    FOREIGN KEY(categoryId) REFERENCES category(id) ON DELETE CASCADE
    	);
    `)

	return err
}

func (c *Course) Create(categoryId, name, description string) (Course, error) {
	err := uuid.Validate(categoryId)
	if err != nil {
		return Course{}, errors.New("categoryId must be a valid UUID")
	}

	id, err := uuid.NewUUID()

	_, err = c.db.Exec("INSERT INTO course (id, name, description, categoryId) VALUES (?, ?, ?, ?)", id.String(), name, description, categoryId)
	if err != nil {
		return Course{}, err
	}

	return Course{nil, id.String(), name, description, categoryId}, nil
}

func (c *Course) FindAll() ([]Course, error) {
	courses := make([]Course, 0)

	rows, err := c.db.Query("SELECT id, name, description, categoryId FROM course")
	if err != nil {
		return courses, err
	}
	defer rows.Close()

	for rows.Next() {
		var course Course
		err = rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.CategoryID,
		)
		if err != nil {
			continue
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (c *Course) FindAllFromCategory(categoryId string) ([]Course, error) {
	courses := make([]Course, 0)

	err := uuid.Validate(categoryId)
	if err != nil {
		return courses, errors.New("categoryId must be a valid UUID")
	}

	rows, err := c.db.Query("SELECT id, name, description, categoryId FROM course where categoryId = ?", categoryId)
	if err != nil {
		return courses, err
	}
	defer rows.Close()

	for rows.Next() {
		var course Course
		err = rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.CategoryID,
		)
		if err != nil {
			continue
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (c *Course) FindCategory(courseId string) (Category, error) {
	err := uuid.Validate(courseId)
	if err != nil {
		return Category{}, errors.New("courseId must be a valid UUID")
	}

	rows, err := c.db.Query("SELECT ca.id, ca.name, ca.description FROM category ca inner join course co on ca.id = co.categoryId where co.id = ?", courseId)
	if err != nil {
		return Category{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)
		if err != nil {
			return Category{}, err
		}

		return category, nil
	}

	return Category{}, errors.New("no category found")
}
