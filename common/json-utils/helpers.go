package common

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	out, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	if len(headers) > 0 {
		for key, value := range headers[0] { //why index 0?
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return WriteJSON(w, statusCode, payload)
}
