package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &requestPayload)

	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.ErrorJSON(w, errors.New("invalid cred"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJSON(w, errors.New("invalid creds "), http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)
}
