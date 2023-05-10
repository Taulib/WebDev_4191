// FileName: cmd/api/helpers.go

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	CourseID, err := strconv.ParseInt(params.ByName("CourseID"), 10, 64)
	if err != nil || CourseID < 1 {
		return 0, errors.New("invalid CourseID parameter")
	}
	return CourseID, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// converting map into a json object
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// adding a new line to make the viewing easier
	js = append(js, '\n')
	// add the headers
	for key, value := range headers {
		w.Header()[key] = value
	}

	// showing that we will serve our responses using json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// write the byte slice
	w.Write([]byte(js))
	return nil

}

func (app *application) readJOSN(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1mb
	maxBytes := 1_048_576
	// Decode the request body into the target destination
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := json.NewDecoder(r.Body).Decode(dst)
	// check for a bad request
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshaError *json.InvalidUnmarshalError

		//switch to check for the errors
		switch {
		// checking for syntax errors
		case errors.As(err, &syntaxError):
			return fmt.Errorf("bodt contains badly formed JSON(at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed json")
			// check for wrong types passed by the client
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrent JSON type (at character %d)", unmarshalTypeError.Offset)
			// Empty Body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
			// unmappable fields
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "josn unkown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

			// too large request
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
			//Passing a non-nil pointer error; if the promamer messes up
		case errors.As(err, &invalidUnmarshaError):
			panic(err)
			// default
		default:
			return err
		}
	}
	// Call Decode() again
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single value")

	}
	return nil
}
