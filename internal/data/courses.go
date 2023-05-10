// FileName : internal/data/courses.go

package data

import (
	"database/sql"
	"errors"
	"time"
)

type Courses struct {
	CourseID    int64     `json:"id"`
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
		INSERT INTO courses (CourseName, CreditHours)
		VALUES ($1, $2)
		RETURNING CourseID, CreatedAT
	`
	// Collect the data fields into a slice
	args := []interface{}{
		Courses.CourseName,
		Courses.CreditHours,
	}
	return m.DB.QueryRow(query, args...).Scan(&Courses.CourseID, &Courses.CreditHours, &Courses.CreatedAt)

}

// Get() will allow us to retrieve a specific Course
func (m CourseModel) Get(id int64) (*Courses, error) {
	// Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create our query
	query := `
		SELECT CourseID, CourseName, CreditHours
		FROM courses 
		WHERE CourseID = $1
	`
	// dECLARE A COURSE VARIABLE to hold the returned data
	var Courses Courses
	// execute the query using QueryRow()
	err := m.DB.QueryRow(query, id).Scan(
		&Courses.CourseID,
		&Courses.CourseName,
		&Courses.CreditHours,
		&Courses.CreatedAt,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &Courses, nil
}

// Update() allows us to update or edit a specfic course
func (m CourseModel) Update(Courses *Courses) error {
	// Create a query
	query := `
		UPDATE courses
		SET CourseName = $1, CreditHours = $2
		WHERE CourseID = $3
		AND CreditHours = $3
		RETURNING CreditHours
	`
	args := []interface{}{
		Courses.CourseName,
		Courses.CreditHours,
	}

	return m.DB.QueryRow(query, args...).Scan(&Courses.CreditHours)

}

// Delete will remove a specified course
func (m CourseModel) Delete(id int64) error {
	return nil
}
