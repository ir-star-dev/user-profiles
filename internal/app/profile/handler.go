package profile

import (
	"net/http"
	"strconv"
	"user-profiles/internal/http/req"
	"user-profiles/internal/http/resp"
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
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	resp.Json(w, profile, http.StatusOK)
}

func (handler *ProfileHandler) ViewAll(w http.ResponseWriter, r *http.Request) {
	users, err := handler.Service.ViewAll()
	if err != nil {
		resp.Json(w, EmptyUsers, http.StatusNotFound)
		return
	}
	var profiles []ProfileResponseForAdmin
	for _, user := range users {
		profiles = append(profiles, ProfileResponseForAdmin{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		})
	}
	resp.Json(w, profiles, http.StatusOK)
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
	_, err = handler.Service.ChangeName(uId, res.Name)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Json(w, "Updated", http.StatusOK)
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
	if idStr := r.PathValue("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			return id, nil
		}
	}

	uId, err := middleware.GetUserID(r.Context())
	if err != nil {
		return 0, err
	}
	return uId, nil
}
