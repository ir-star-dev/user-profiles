package profile

import (
	"net/http"
	"user-profiles/pkg/req"
	"user-profiles/pkg/resp"
	"user-profiles/internal/middleware"
)

func (handler *ProfileHandler) View(w http.ResponseWriter, r *http.Request) {
	uId, err := getUserId(r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := handler.Service.View(uId)
	if err != nil {
		resp.Json(w, UserNotFound, http.StatusNotFound)
		return
	}
	profile := &ProfileResponse{
		Name: user.Name,
		Email: user.Email,
	}
	resp.Json(w, profile, http.StatusOK)
}

func (handler *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	uId, err := getUserId(r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusUnauthorized)
		return
	}
	res, err := req.HandleBody[UpdateNameRequest](&w, r)
	if err != nil {
		resp.Json(w, UserNotFound, http.StatusNotFound)
		return
	}
	u, err := handler.Service.UpdateName(uId, res.Name)
	if err != nil {
		resp.Json(w, UpdateFailed, http.StatusInternalServerError)
		return
	}
	profile := &ProfileResponse{
		Name: u.Name,
		Email: u.Email,
	}
	resp.Json(w, profile, http.StatusOK)
}

func (handler *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uId, err := getUserId(r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusUnauthorized)
		return
	}
	err = handler.Service.Delete(uId)
	if err != nil {
		resp.Json(w, DeleteFailed, http.StatusInternalServerError)
		return
	}
	resp.Json(w, "Profile deleted", http.StatusOK)
}

func getUserId(r *http.Request) (int, error) {
	uId, err := middleware.GetUserID(r.Context())
	if err != nil {
		return 0, err
	}
	return uId, nil
}