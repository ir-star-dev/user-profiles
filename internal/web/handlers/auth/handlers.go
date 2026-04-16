package auth_handlers

import (
	"net/http"
	//"user-profiles/internal/middleware"
	auth_dto "user-profiles/internal/dto/auth"
	"user-profiles/internal/http/req"
	"errors"
)

func (handler *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := handler.Tmpl.ExecuteTemplate(w, "pages/login", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
// 	body, err := req.HandleBody[LoginInput](&w, r)
// 	if err != nil {
// 		return
// 	}
// 	tokens, err := handler.Service.Login(body.Email, body.Password)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	// Send refresh like HttpOnly cookie
// 	setCookie(tokens.Refresh, w)
// 	res := &LoginResponse{
// 		Token: tokens.Access,
// 	}
// 	resp.Json(w, res, http.StatusOK)
// }

func (handler *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	err := handler.Tmpl.ExecuteTemplate(w, "pages/register", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	input := &auth_dto.RegisterInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Role:     r.FormValue("role"),
	}
	if err := req.IsValid(input); err != nil {
		var ve req.ValidationErrors
		if errors.As(err, &ve) {
			if err := handler.Tmpl.ExecuteTemplate(w, "parts/validation/errors", ve); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	if err := handler.Service.Register(input.Email, input.Password, input.Name, input.Role); err != nil {
		handler.Tmpl.ExecuteTemplate(w, "parts/form-submit/error", err.Error())
		return
	}

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}

// func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
// 	userId, err := middleware.GetUserID(r.Context())
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	token, err := getRefreshTokenFromCookie(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err = handler.Service.Logout(userId, token)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}	
// 	setCookie("", w)
// 	resp.Json(w, "Logged out", http.StatusOK)
// }

// func (handler *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
// 	token, err := getRefreshTokenFromCookie(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	tokens, err := handler.Service.Refresh(token)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	setCookie(tokens.Refresh, w)
// 	res := &LoginResponse{
// 		Token: tokens.Access,
// 	}
// 	resp.Json(w, res, http.StatusOK)
// }