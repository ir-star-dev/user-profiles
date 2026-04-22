package handlers

import (
	"net/http"
	// "user-profiles/internal/http/req"
	// "user-profiles/internal/http/resp"
	// "user-profiles/internal/middleware"
	"html/template"
	"user-profiles/configs"

	"user-profiles/internal/auth"
	"user-profiles/internal/users"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	Config     *configs.Config
	Service    users.AdminService
	JWTService auth.JWTService
	Tmpl       *template.Template
}

type AdminHandlerDeps struct {
	Config     *configs.Config
	Service    users.AdminService
	JWTService auth.JWTService
	Tmpl       *template.Template
}

func NewAdminHandler(router chi.Router, deps AdminHandlerDeps) *AdminHandler {
	return &AdminHandler{
		Config:     deps.Config,
		Service:    deps.Service,
		JWTService: deps.JWTService,
		Tmpl:       deps.Tmpl,
	}
}

func (handler *AdminHandler) AdminProfilePage(w http.ResponseWriter, r *http.Request) {
	err := handler.Tmpl.ExecuteTemplate(w, "profile", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *AdminHandler) UserProfilePage(w http.ResponseWriter, r *http.Request) {
	err := handler.Tmpl.ExecuteTemplate(w, "profile", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func (handler *Handler) ViewById(w http.ResponseWriter, r *http.Request) {
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

// func (handler *Handler) UpdateById(w http.ResponseWriter, r *http.Request) {
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

// func (handler *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
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


// func (handler *Handler) ViewByAdmin(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getIdFromReq(w, r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusForbidden)
// 		return
// 	}
// 	user, err := handler.Service.View(uId)
// 	if err != nil {
// 		resp.Json(w, UserNotFound, http.StatusNotFound)
// 		return
// 	}
// 	profile := &FullProfileResponseForAdmin{
// 		Id:        user.Id,
// 		Name:      user.Name,
// 		Role:      user.Role,
// 		Email:     user.Email,
// 		Banned:    user.Banned,
// 		CreatedAt: user.CreatedAt,
// 	}
// 	resp.Json(w, profile, http.StatusOK)
// }

// func (handler *Handler) ViewAllByAdmin(w http.ResponseWriter, r *http.Request) {
// 	users, err := handler.Service.ViewAll()
// 	if err != nil {
// 		resp.Json(w, EmptyUsers, http.StatusNotFound)
// 		return
// 	}
// 	var profiles []FullProfileResponseForAdmin
// 	for _, user := range users {
// 		profiles = append(profiles, FullProfileResponseForAdmin{
// 			Id:        user.Id,
// 			Name:      user.Name,
// 			Role:      user.Role,
// 			Email:     user.Email,
// 			Banned:    user.Banned,
// 			CreatedAt: user.CreatedAt,
// 		})
// 	}
// 	resp.Json(w, profiles, http.StatusOK)
// }

// func (handler *Handler) UpdateByAdmin(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getIdFromReq(w, r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusForbidden)
// 		return
// 	}
// 	res, err := req.HandleBody[UpdateNameRequest](&w, r)
// 	if err != nil {
// 		resp.Json(w, UserNotFound, http.StatusNotFound)
// 		return
// 	}
// 	user, err := handler.Service.ChangeName(uId, res.Name)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	profile := &ShortProfileResponseForAdmin{
// 		Id:        user.Id,
// 		Name:      user.Name,
// 	}
// 	resp.Json(w, profile, http.StatusOK)
// }

// func (handler *Handler) DeleteByAdmin(w http.ResponseWriter, r *http.Request) {
// 	uId, err := selfDeletionDetected(w, r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err = handler.Service.Delete(*uId)
// 	if err != nil {
// 		resp.Json(w, DeleteFailed, http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "Profile deleted", http.StatusOK)
// }

// func (handler *Handler) BanByAdmin(w http.ResponseWriter, r *http.Request) {
// 	uId, err := selfDeletionDetected(w, r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err = handler.Service.Ban(*uId)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "User banned", http.StatusOK)
// }

// func (handler *Handler) UnbanByAdmin(w http.ResponseWriter, r *http.Request) {
// 	uId, err := selfDeletionDetected(w, r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err = handler.Service.Unban(*uId)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "User unbanned", http.StatusOK)
// }

// func getIdFromReq(w http.ResponseWriter, r *http.Request) (int, error) {
// 	idStr := strings.TrimSpace(r.PathValue("id"))
// 	if idStr == "" {
// 		resp.Json(w, EmptyParams, http.StatusForbidden)
// 		return 0, errors.New(EmptyParams)
// 	}
// 	uId, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return 0, errors.New(WrongParam)
// 	}
// 	return uId, nil
// }

// func selfDeletionDetected(w http.ResponseWriter, r *http.Request) (*int, error) {
// 	uId, err := getIdFromReq(w, r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	adminId, err := getUserId(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if uId == adminId {
// 		return nil, errors.New(DeleteYourself)
// 	}
// 	return &uId, nil
// }