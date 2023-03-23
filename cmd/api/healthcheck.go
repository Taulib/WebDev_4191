// FileName: cmd/api/healthcheck.go

package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	//create a map to hold healthcheck data
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	// converting map into a json object
	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The Server has ran into a error and could not meet your expectations", http.StatusInternalServerError)
		return
	}

	// adding a new line to make the viewing easier
	js = append(js, '\n')
	// showing that we will serve our responses using json
	w.Header().Set("Content-Type", "application/json")
	// write the byte slice
	w.Write(js)

}
