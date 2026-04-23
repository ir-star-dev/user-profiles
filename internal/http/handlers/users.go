package handlers

import (
	"net/http"
	// "user-profiles/internal/http/req"
	// "user-profiles/internal/http/resp"
	// "user-profiles/internal/middleware"
	//"html/template"
	"user-profiles/cmd/view"
	"user-profiles/configs"

	"user-profiles/internal/auth"
	"user-profiles/internal/users"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Config       *configs.Config
	AuthService  auth.AuthService
	UsersService users.UsersService
	JWTService   auth.JWTService
}

type UserHandlerDeps struct {
	Config       *configs.Config
	AuthService  auth.AuthService
	UsersService users.UsersService
	JWTService   auth.JWTService
}

func NewUserHandler(router chi.Router, deps UserHandlerDeps) *UserHandler {
	return &UserHandler{
		Config:       deps.Config,
		AuthService:  deps.AuthService,
		UsersService: deps.UsersService,
		JWTService:   deps.JWTService,
	}
}

func (handler *UserHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/profile.tmpl",
		"././ui/templates/parts/layout/user-data.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uId, err := getUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	user, err := handler.UsersService.View(uId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	profile := &users.UsersProfileResponse{
		Name:   user.Name,
		Email:  user.Email,
		Banned: user.Banned,
	}
	err = tmpl.ExecuteTemplate(w, "profile", profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func (handler *Handler) ViewById(w http.ResponseWriter, r *http.Request) {

// 	user, err := handler.AuthService.View(uId)
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
// 	_, err = handler.AuthService.ChangeName(uId, res.Name)
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
// 	err = handler.AuthService.Delete(uId)
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
