//Filename: internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record Not Found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a wrapper for our data models
type Models struct {
	Courses CourseModel
}

// NewModels() will allow us to create a new model
func NewModels(db *sql.DB) Models {
	return Models{
		Courses: CourseModel{DB: db},
	}
}
