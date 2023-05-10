//fileName: cmd/api/errors.go

package main

import (
	"fmt"
	"net/http"
)

// logging an error to the console
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// To send JSON-Formatted error message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	// Create a JSON Response

	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Server Error Response
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	// Prepare a message with the error
	message := "the server encountered a problem and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The Not Found response
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	// Create our Message
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// A method not allowed Response
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	// Create our Message
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// User Provides a bad request
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	// Create our Message
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
