package main

import (
	"bytes"
	"encoding/json"
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
	fmt.Println("Json request is : ", requestPayload)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		fmt.Println("invalid creds")
		app.ErrorJSON(w, errors.New("invalid cred"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJSON(w, errors.New("invalid creds "), http.StatusBadRequest)
		return
	}

	//log authentication

	err = app.LogRequest("authentication", fmt.Sprintf("%s  logged in", user.Email))

	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	fmt.Println("response payload ", payload)
	app.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) LogRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
