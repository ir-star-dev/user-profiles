package resp

import (
	"encoding/json"
	"net/http"
)

type JsonData struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

func Json(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	JsonData := &JsonData{
		Code:    statusCode,
		Message: data,
	}
	json.NewEncoder(w).Encode(JsonData)
}
