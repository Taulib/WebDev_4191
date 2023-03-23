// FileName: cmd/api/courses.go

package main

import (
	"fmt"
	"net/http"
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
	// Display the course ID
	fmt.Fprintf(w, "show the details for courses %d\n", id)

}
