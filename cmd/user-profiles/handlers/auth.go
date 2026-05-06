package handlers

import (
	"net/http"
	"strconv"
	"time"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/view"
	"user-profiles/configs"

	"errors"
	"user-profiles/internal/http/cookie"
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/http/req"

	"github.com/go-chi/chi/v5"
)

type AuthHandlerDeps struct {
	Config      *configs.Config
	AuthService auth.AuthService
	JWTService  auth.JWTService
}

type AuthHandler struct {
	Config      *configs.Config
	AuthService auth.AuthService
	JWTService  auth.JWTService
}

func NewAuthHandler(router chi.Router, deps AuthHandlerDeps) *AuthHandler {
	return &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
		JWTService:  deps.JWTService,
	}
}

func (handler *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/parts/layout/nav.tmpl",
		"././ui/templates/pages/login.tmpl",
		"././ui/templates/parts/auth/form-login.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := view.PageData{}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/parts/layout/nav.tmpl",
		"././ui/templates/pages/register.tmpl",
		"././ui/templates/parts/auth/form-register.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := view.PageData{}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	input := &auth.LoginInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/parts/layout/nav.tmpl",
		"././ui/templates/pages/register.tmpl",
		"././ui/templates/parts/auth/form-register.tmpl",
		"././ui/templates/parts/validation/errors.tmpl",
		"././ui/templates/parts/form-submit/error.tmpl",
	)
	res, err := handler.AuthService.Login(input.Email, input.Password)
	if err != nil {
		tmpl.ExecuteTemplate(w, "form-submit-error", err.Error())
		return
	}
	cookie.Set(res.Access, "__up_access_token", 5*time.Minute, w)
	cookie.Set(res.Refresh, "__up_refresh_token", 7*24*time.Hour, w)

	w.Header().Set("HX-Redirect", "/panel/profile/"+ strconv.Itoa(res.UserId))
	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	input := &auth.RegisterInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Role:     r.FormValue("role"),
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/parts/layout/nav.tmpl",
		"././ui/templates/pages/register.tmpl",
		"././ui/templates/parts/auth/form-register.tmpl",
		"././ui/templates/parts/validation/errors.tmpl",
		"././ui/templates/parts/form-submit/error.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := req.IsValid(input); err != nil {
		var ve req.ValidationErrors
		if errors.As(err, &ve) {
			if err := tmpl.ExecuteTemplate(w, "validation-error", ve); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	if err := handler.AuthService.Register(input.Email, input.Password, input.Name, input.Role); err != nil {
		tmpl.ExecuteTemplate(w, "form-submit-error", err.Error())
		return
	}

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetUserID(r.Context())
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
	w.WriteHeader(http.StatusOK)
}
