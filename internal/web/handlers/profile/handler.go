package profile_handlers

import (
	"net/http"
	// "user-profiles/internal/http/req"
	// "user-profiles/internal/http/resp"
	// "user-profiles/internal/middleware"
)

func (handler *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	err := handler.Tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


// func (handler *Handler) View(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getUserId(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	user, err := handler.Service.View(uId)
// 	if err != nil {
// 		resp.Json(w, UserNotFound, http.StatusNotFound)
// 		return
// 	}
// 	profile := &ProfileResponse{
// 		Name:  user.Name,
// 		Email: user.Email,
// 		Banned: user.Banned,
// 	}
// 	resp.Json(w, profile, http.StatusOK)
// }

// func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getUserId(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	res, err := req.HandleBody[UpdateNameRequest](&w, r)
// 	if err != nil {
// 		resp.Json(w, UserNotFound, http.StatusNotFound)
// 		return
// 	}
// 	_, err = handler.Service.ChangeName(uId, res.Name)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "Profile updated", http.StatusOK)
// }

// func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getUserId(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	err = handler.Service.Delete(uId)
// 	if err != nil {
// 		resp.Json(w, DeleteFailed, http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "Profile deleted", http.StatusOK)
// }

// func getUserId(r *http.Request) (int, error) {
// 	uId, err := middleware.GetUserID(r.Context())
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uId, nil
// }
