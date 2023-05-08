// FileName: cmd/api/courses.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"quiz1.taulib.net/internal/data"
)

// createCourseHandler for the "POST /v1/courses"
func (app *application) createCoursesHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "create a course ")
}

// showsCoursesHandler for "GET /v1/courses/:id"
func (app *application) showCoursesHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return

	}
	// create a new instance of the course struct containing the ID we extracted from our URL and some sample data
	courses := data.Courses{
		ID:          id,
		CourseName:  "course name test",
		CreditHours: 3,
		CreatedAt:   time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, courses, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not retrun your request", http.StatusInternalServerError)
	}

}
