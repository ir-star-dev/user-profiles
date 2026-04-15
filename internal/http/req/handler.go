package req

import (
	"errors"
	"net/http"
	"user-profiles/internal/http/resp"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		resp.Json(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	err = IsValid(body)
	if err != nil {
		var ve ValidationErrors
		if errors.As(err, &ve) {
			resp.Json(*w, ve, http.StatusBadRequest)
			return nil, err
		}

		resp.Json(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}
