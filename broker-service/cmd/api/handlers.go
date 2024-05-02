package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	Payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = app.WriteJSON(w, http.StatusOK, Payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var RequestPayload RequestPayload

	err := app.ReadJSON(w, r, &RequestPayload)

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, RequestPayload.Auth)
	case "log":
		app.LogItem(w, RequestPayload.Log)

	default:
		app.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) LogItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New(response.Status))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged successfully"

	app.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the service

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	//make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJSON(w, errors.New("invalid creds "))
		return
	} else if response.StatusCode != http.StatusAccepted {
		fmt.Println("this is status : ", response.StatusCode)
		app.ErrorJSON(w, errors.New("error calling auth service "))
		return
	}

	// create a variable we will read response body into
	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.ErrorJSON(w, err, http.StatusUnauthorized)
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.WriteJSON(w, http.StatusAccepted, payload)

}
