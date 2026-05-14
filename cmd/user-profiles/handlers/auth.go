package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/configs"

	"user-profiles/internal/http/cookie"

	"github.com/go-chi/chi/v5"
)

type AuthHandlerDeps struct {
	Config      *configs.Config
	AuthService auth.AuthService
	JWTService  auth.JWTService
	TCache      panel.Templates
}

type AuthHandler struct {
	Config      *configs.Config
	AuthService auth.AuthService
	JWTService  auth.JWTService
	TCache      panel.Templates
}

func NewAuthHandler(router chi.Router, deps AuthHandlerDeps) *AuthHandler {
	return &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
		JWTService:  deps.JWTService,
		TCache:      deps.TCache,
	}
}

func (handler *AuthHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	data := panel.PageData{}
	handler.TCache.Render(w, r, http.StatusOK, "base", "login.tmpl", data)
}

func (handler *AuthHandler) SignupForm(w http.ResponseWriter, r *http.Request) {
	data := panel.PageData{}
	handler.TCache.Render(w, r, http.StatusOK, "base", "signup.tmpl", data)
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	data, err := handler.AuthService.Login(email, password)
	if err != nil {
		err = handler.TCache.RenderPartial(w, "login.tmpl", "form-submit-error", data)
		if err != nil {
			handler.TCache.ServerError(w, err)
		}
		return
	}
	cookie.Set(data.Access, "__up_access_token", 5*time.Minute, w)
	cookie.Set(data.Refresh, "__up_refresh_token", 7*24*time.Hour, w)

	w.Header().Set("HX-Redirect", "/panel/profile/"+strconv.Itoa(data.UserId))
}

func (handler *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	name := strings.TrimSpace(r.FormValue("name"))
	role := strings.TrimSpace(r.FormValue("role"))

	data, err := handler.AuthService.Register(email, password, name, role)
	if err != nil {
		err := handler.TCache.RenderPartial(w, "signup.tmpl", "form-submit-error", data)
		if err != nil {
			handler.TCache.ServerError(w, err)
		}
		return
	}
	w.Header().Set("HX-Redirect", "/auth/login")
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	token, err := cookie.Get("__up_refresh_token", r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	err = handler.AuthService.Logout(userId, token)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	cookie.Set("", "__up_access_token", -time.Minute, w)
	cookie.Set("", "__up_refresh_token", -time.Hour, w)

	w.Header().Set("HX-Redirect", "/auth/login")
}
