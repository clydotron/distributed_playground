package main

import (
	common "common/json-utils"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
)

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

	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	log.Println("WriteLog:", event.Name, event.Data)
	if err := app.dataStore.LogRepo.Insert(event); err != nil {
		fmt.Println("ERROR! >>", err)
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
