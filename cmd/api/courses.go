// FileName: cmd/api/courses.go

package main

import (
	"errors"
	"fmt"
	"net/http"

	"quiz1.taulib.net/internal/data"
)

// createCourseHandler for the "POST /v1/courses"
func (app *application) createCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// our taget decode destination
	var input struct {
		CourseName  string `json:"name,omitempty"`
		CreditHours int64  `josn:"hours"`
	}
	Courses := &data.Courses{
		CourseName:  input.CourseName,
		CreditHours: input.CreditHours,
	}
	// Initialize a new json.Decoder instance
	err := app.readJOSN(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// display the request
	fmt.Fprintf(w, "%+v\n", input)

	// create a school
	err = app.Models.Courses.Insert(Courses)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// create a location header for new created resource/course
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("v1,courses/%d", Courses.CourseID))
	// write the JSON response with 201 - Created status code wiht the body being the schooldata and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"course": Courses}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showsCoursesHandler for "GET /v1/courses/:id"
func (app *application) showCoursesHandler(w http.ResponseWriter, r *http.Request) {

	CourseID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return

	}
	// fetch the specific Course
	Courses, err := app.Models.Courses.Get(CourseID)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Courses": Courses}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateCoursesHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a partial replacement
	// Get the id for the school that needs updating
	CourseID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the orginal record from the database
	Courses, err := app.Models.Courses.Get(CourseID)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Create an input struct to hold data read in from the client
	// We update input struct to use pointers because pointers have a
	// default value of nil
	// If a field remains nil then we know that the client did not update it
	var input struct {
		CourseName  string `json:"name,omitempty"`
		CreditHours int64  `josn:"hours"`
	}

	/*// Initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}*/

	// update the fields /values in the school variable using the fields in the input struct
	Courses.CourseName = input.CourseName
	Courses.CreditHours = input.CreditHours

	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"Courses": Courses}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCoursesHandler(w http.ResponseWriter, r *http.Request) {

	CourseID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return

	}

	// deleting the Course form the database, send a 404 not found status code to the client if there is no matching record
	err = app.Models.Courses.Delete(CourseID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Return the 200 Status Ok to the Client with a sucess message
	err = app.writeJSON(w, http.StatusOK, envelope{"Message": "Course Sucessfully Deleted!"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
