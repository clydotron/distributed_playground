package main

import (
	"bytes"
	common "common/json-utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// TODO use oneof
type requestPayload struct {
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

const (
	authServiceURL = "http://auth-service/authenticate"
	logServiceURL  = "http://log-service/log"
)

func (app *App) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hit the broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *App) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var payload requestPayload

	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.ErrorJSON(w, err)
		return
	}
	switch payload.Action {
	case "auth":
		app.authenticate(w, payload.Auth)
	case "log":
		app.logItem(w, payload.Log)
	default:
		common.ErrorJSON(w, errors.New("unknown action"))
	}

}

func (app *App) authenticate(w http.ResponseWriter, auth AuthPayload) {
	// send request to the authentication service:
	jsonData, _ := json.MarshalIndent(auth, "", "\t") //TODO remove before prod
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		common.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		common.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	switch response.StatusCode {
	case http.StatusUnauthorized:
		common.ErrorJSON(w, errors.New("invalid credentials"))
		return
	case http.StatusAccepted:
	default:
		common.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// process the response from the auth service:
	var jsonFromService jsonResponse

	// decode the json from the auth service
	if err = json.NewDecoder(response.Body).Decode(&jsonFromService); err != nil {
		common.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		common.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// send message to the log service
	err = app.logRequest("authentication", fmt.Sprintf("%s successfully logged in", auth.Email))
	if err != nil {
		fmt.Printf("Error logging auth status:%v\n", err)
	}
	sendResponse(w, false, "Authenticated", jsonFromService.Data)
}

func (app *App) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		common.ErrorJSON(w, err) //wrap it?
		return
	}

	log.Printf("log item: %v", entry)

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		common.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		common.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Println("Failed:", response.StatusCode)
		common.ErrorJSON(w, fmt.Errorf("failed to contact log service: %d", response.StatusCode))
		return
	}

	sendResponse(w, false, "logged", nil)
}

func sendResponse(w http.ResponseWriter, err bool, msg string, data any) {
	payload := jsonResponse{
		Error:   err,
		Message: msg,
		Data:    data,
	}
	common.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *App) logRequest(name, data string) error {
	payload := LogPayload{Name: name, Data: data}
	log.Printf("logRequest: %s %s", name, data)

	// uses http for now: replace with
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

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
