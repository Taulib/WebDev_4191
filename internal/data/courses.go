// FileName : internal/data/courses.go

package data

import (
	"database/sql"
	"time"
)

type Courses struct {
	ID          int64     `json:"id"`
	CourseName  string    `json:"name,omitempty"`
	CreditHours int64     `josn:"hours"`
	CreatedAt   time.Time `json:"-"`
}

// Defining the course model to wrap a sql.DB connection pool
type CourseModel struct {
	DB *sql.DB
}

// Inserting data into our new courses
func (m CourseModel) Insert(Courses *Courses) error {
	query := `
		INSERT INTO courses (courseName, creditHours)
		VALUES ($1, $2)
		RETURNING id, createdAT, version
	`
	// Collect the data fields into a slice
	args := []interface{}{
		Courses.CourseName,
		Courses.CreditHours,
	}
	return m.DB.QueryRow(query, args...).Scan(&Courses.ID, &Courses.CreditHours, &Courses.CreatedAt)

}

// Get() will allow us to retrieve a specific Course
func (m CourseModel) Get(id int64) (*Courses, error) {
	return nil, nil
}

// Update() allows us to update or edit a specfic course
func (m CourseModel) Update(course *Courses) error {
	return nil
}

// Delete will remove a specified course
func (m CourseModel) Delete(id int64) error {
	return nil
}
