package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	Payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = app.WriteJSON(w, http.StatusOK, Payload)

}
