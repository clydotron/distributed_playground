package main

import (
	common "common/json-utils"
	"errors"
	"fmt"
	"net/http"
)

func (a *App) Authenticate(w http.ResponseWriter, r *http.Request) {
	// TODO put these someplace common (but not in common)
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := common.ReadJSON(w, r, &requestPayload); err != nil {
		common.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := a.UserRepo.GetByEmail(requestPayload.Email)
	if err != nil {
		common.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		common.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := common.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
	}
	common.WriteJSON(w, http.StatusAccepted, payload)
}
