package handlers

import (
	"net/http"
	"strconv"
	"time"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/configs"

	"errors"
	"user-profiles/internal/http/cookie"
	"user-profiles/internal/http/req"

	"github.com/go-chi/chi/v5"
)

type AuthHandlerDeps struct {
	Config      *configs.Config
	AuthService auth.AuthService
	JWTService  auth.JWTService
	Templates   panel.Templates
}

type AuthHandler struct {
	Config      *configs.Config
	AuthService auth.AuthService
	JWTService  auth.JWTService
	Templates   panel.Templates
}

func NewAuthHandler(router chi.Router, deps AuthHandlerDeps) *AuthHandler {
	return &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
		JWTService:  deps.JWTService,
		Templates:   deps.Templates,
	}
}

func (handler *AuthHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	data := panel.PageData{}
	err := handler.Templates.Render(w, "login", data,
		"././ui/pages/login.tmpl",
		"././ui/parts/forms/login.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *AuthHandler) SignupForm(w http.ResponseWriter, r *http.Request) {
	data := panel.PageData{}
	err := handler.Templates.Render(w, "signup", data,
		"././ui/pages/signup.tmpl",
		"././ui/parts/forms/signup.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	input := &auth.LoginInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	res, err := handler.AuthService.Login(input.Email, input.Password)
	if err != nil {
		err = handler.Templates.RenderPartial(
			w,
			"form-submit-error",
			auth.AuthViewError{
				Message: err.Error(),
			},
			"././ui/parts/validation/errors.tmpl",
			"./ui/parts/forms/error.tmpl",
		)

		if err != nil {
			handler.Templates.ServerError(w, err)
		}
		return
	}
	cookie.Set(res.Access, "__up_access_token", 5*time.Minute, w)
	cookie.Set(res.Refresh, "__up_refresh_token", 7*24*time.Hour, w)

	w.Header().Set("HX-Redirect", "/panel/profile/"+strconv.Itoa(res.UserId))
	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	input := &auth.RegisterInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Role:     r.FormValue("role"),
	}
	if err := req.IsValid(input); err != nil {
		var ve req.ValidationErrors
		if errors.As(err, &ve) {
			err := handler.Templates.RenderPartial(
				w,
				"validation-error",
				ve,
				"././ui/parts/validation/errors.tmpl",
			)
			if err != nil {
				handler.Templates.ServerError(w, err)
			}

			return
		}
	}

	err := handler.AuthService.Register(
		input.Email,
		input.Password,
		input.Name,
		input.Role,
	)
	if err != nil {
		err = handler.Templates.RenderPartial(
			w,
			"form-submit-error",
			err.Error(),

			"././ui/parts/forms/error.tmpl",
		)
		if err != nil {
			handler.Templates.ServerError(w, err)
		}
		return
	}
	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)
}
