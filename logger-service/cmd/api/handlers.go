package main

import (
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	//read the json into a var

	var requestPayload JSONPayload
	_ = app.ReadJSON(w, r, &requestPayload)

	//insert Data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Log entry inserted",
	}

	app.WriteJSON(w, http.StatusAccepted, resp)
}
