package resp

import (
	"encoding/json"
	"net/http"
)

type JsonData struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

func Json(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	JsonData := &JsonData{
		Code: statusCode,
		Data: data,
	}
	json.NewEncoder(w).Encode(JsonData)
}
