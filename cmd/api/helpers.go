// FileName: cmd/api/helpers.go

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("Invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	// converting map into a json object
	js, err := json.Marshal(data)
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
	w.Write(js)
	return nil

}
