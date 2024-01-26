package main

import (
	common "common/json-utils"
	"log"
	"log-service/models"
	"net/http"
)

// TODO figure out better place to put these
type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *App) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload jsonPayload
	err := common.ReadJSON(w, r, &payload)
	if err != nil {
		log.Println("error reading the json:", err)
		return
	}

	event := models.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	if err := app.DataStore.LogRepo.Insert(event); err != nil {
		common.ErrorJSON(w, err)
		return
	}

	log.Println("successfully wrote logs to DB")
	resp := common.JsonResponse{
		Error:   false,
		Message: "logged",
	}
	common.WriteJSON(w, http.StatusAccepted, resp)
}
