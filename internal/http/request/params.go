package request

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func GetIdFromReq(r *http.Request) (int, error) {
	idStr := strings.TrimSpace(r.PathValue("id"))
	if idStr == "" {
		return 0, errors.New("Missing param")
	}
	uId, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uId, nil
}

func GetUserId(r *http.Request) (int, error) {
	uId, err := GetUserID(r.Context())
	if err != nil {
		return 0, err
	}
	return uId, nil
}

func GetPageFromReq(r *http.Request) int {
	page := r.URL.Query().Get("page")
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		return 1
	}
	return p
}

func GetFilterValue(r *http.Request, key string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return ""
	}
	return value
}
