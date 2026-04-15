package profile

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"user-profiles/internal/http/req"
	"user-profiles/internal/http/resp"
)

func (handler *ProfileHandler) ViewByAdmin(w http.ResponseWriter, r *http.Request) {
	uId, err := getIdFromReq(w, r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusForbidden)
		return
	}
	user, err := handler.Service.View(uId)
	if err != nil {
		resp.Json(w, UserNotFound, http.StatusNotFound)
		return
	}
	profile := &FullProfileResponseForAdmin{
		Id:        user.Id,
		Name:      user.Name,
		Role:      user.Role,
		Email:     user.Email,
		Banned:    user.Banned,
		CreatedAt: user.CreatedAt,
	}
	resp.Json(w, profile, http.StatusOK)
}

func (handler *ProfileHandler) ViewAllByAdmin(w http.ResponseWriter, r *http.Request) {
	users, err := handler.Service.ViewAll()
	if err != nil {
		resp.Json(w, EmptyUsers, http.StatusNotFound)
		return
	}
	var profiles []FullProfileResponseForAdmin
	for _, user := range users {
		profiles = append(profiles, FullProfileResponseForAdmin{
			Id:        user.Id,
			Name:      user.Name,
			Role:      user.Role,
			Email:     user.Email,
			Banned:    user.Banned,
			CreatedAt: user.CreatedAt,
		})
	}
	resp.Json(w, profiles, http.StatusOK)
}

func (handler *ProfileHandler) UpdateByAdmin(w http.ResponseWriter, r *http.Request) {
	uId, err := getIdFromReq(w, r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusForbidden)
		return
	}
	res, err := req.HandleBody[UpdateNameRequest](&w, r)
	if err != nil {
		resp.Json(w, UserNotFound, http.StatusNotFound)
		return
	}
	user, err := handler.Service.ChangeName(uId, res.Name)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile := &ShortProfileResponseForAdmin{
		Id:        user.Id,
		Name:      user.Name,
	}
	resp.Json(w, profile, http.StatusOK)
}

func (handler *ProfileHandler) DeleteByAdmin(w http.ResponseWriter, r *http.Request) {
	uId, err := selfDeletionDetected(w, r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.Service.Delete(*uId)
	if err != nil {
		resp.Json(w, DeleteFailed, http.StatusInternalServerError)
		return
	}
	resp.Json(w, "Profile deleted", http.StatusOK)
}

func (handler *ProfileHandler) BanByAdmin(w http.ResponseWriter, r *http.Request) {
	uId, err := selfDeletionDetected(w, r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.Service.Ban(*uId)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Json(w, "User banned", http.StatusOK)
}

func (handler *ProfileHandler) UnbanByAdmin(w http.ResponseWriter, r *http.Request) {
	uId, err := selfDeletionDetected(w, r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.Service.Unban(*uId)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Json(w, "User unbanned", http.StatusOK)
}

func getIdFromReq(w http.ResponseWriter, r *http.Request) (int, error) {
	idStr := strings.TrimSpace(r.PathValue("id"))
	if idStr == "" {
		resp.Json(w, EmptyParams, http.StatusForbidden)
		return 0, errors.New(EmptyParams)
	}
	uId, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New(WrongParam)
	}
	return uId, nil
}

func selfDeletionDetected(w http.ResponseWriter, r *http.Request) (*int, error) {
	uId, err := getIdFromReq(w, r)
	if err != nil {
		return nil, err
	}
	adminId, err := getUserId(r)
	if err != nil {
		return nil, err
	}
	if uId == adminId {
		return nil, errors.New(DeleteYourself)
	}
	return &uId, nil
}