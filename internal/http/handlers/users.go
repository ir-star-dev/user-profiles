package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	// "user-profiles/internal/http/req"
	// "user-profiles/internal/http/resp"
	"user-profiles/internal/http/middleware"
	//"html/template"
	"time"
	"user-profiles/cmd/view"
	"user-profiles/configs"
	"user-profiles/internal/http/cookie"

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

	profile := &users.ProfileResponse{
		Id:        user.Id,
		Role:      user.Role,
		Name:      user.Name,
		Email:     user.Email,
		Banned:    user.Banned,
		CreatedAt: user.CreatedAt,
	}

	var tmpl *template.Template
	var tmplName string
	if user.Role == "admin" {
		tmpl, err = view.LoadTemplate(
			"././ui/templates/base.tmpl",
			"././ui/templates/pages/profile-admin.tmpl",
			"././ui/templates/parts/layout/user.tmpl",
		)
		tmplName = "profile-admin"
	} else {
		tmpl, err = view.LoadTemplate(
			"././ui/templates/base.tmpl",
			"././ui/templates/pages/profile.tmpl",
			"././ui/templates/parts/layout/user.tmpl",
		)
		tmplName = "profile"
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, tmplName, profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
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

func (handler *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uId, err := getUserId(r)
	if err != nil {
		return
	}
	err = handler.UsersService.Delete(uId)
	if err != nil {
		return
	}
	cookie.Set("", "__up_access_token", -time.Minute, w)
	cookie.Set("", "__up_refresh_token", -time.Hour, w)

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}

func getUserId(r *http.Request) (int, error) {
	uId, err := middleware.GetUserID(r.Context())
	if err != nil {
		return 0, err
	}
	return uId, nil
}

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
