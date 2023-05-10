// FileName : internal/data/courses.go

package data

import (
	"database/sql"
	"errors"
	"time"
)

// this truct represents a new row of data in our table
type Courses struct {
	CourseID    int64     `json:"Course ID"`
	CourseName  string    `json:"Course Name"`
	CreditHours string    `josn:"Credit Hours"`
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
		RETURNING CourseID
	`
	// Collect the data fields into a slice
	args := []interface{}{
		Courses.CourseName,
		Courses.CreditHours,
	}
	return m.DB.QueryRow(query, args...).Scan(&Courses.CourseID)

}

// Get() will allow us to retrieve a specific Course
func (m CourseModel) Get(CourseID int64) (*Courses, error) {
	// Ensure that there is a valid id
	if CourseID < 1 {
		return nil, ErrRecordNotFound
	}
	// Create our query
	query := `
		SELECT *
		FROM courses 
		WHERE CourseID = $1
	`
	// dECLARE A COURSE VARIABLE to hold the returned data
	var Courses Courses
	// execute the query using QueryRow()
	err := m.DB.QueryRow(query, CourseID).Scan(
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
		AND CreditHours = $4
		RETURNING CreditHours
	`
	args := []interface{}{
		Courses.CourseName,
		Courses.CreditHours,
	}

	return m.DB.QueryRow(query, args...).Scan(&Courses.CreditHours)

}

// Delete will remove a specified course
func (m CourseModel) Delete(CourseID int64) error {
	// Ensure that there is a valid id
	if CourseID < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM courses
		WHERE CourseID = $1
		`
	// execute the query

	result, err := m.DB.Exec(query, CourseID)
	if err != nil {
		return err
	}

	// check how many rows where affected by the delete operation using method called rowsAffected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// check to see if there where zero rows affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
