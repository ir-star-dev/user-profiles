package auth

import (
	"user-profiles/pkg/req"
	"user-profiles/pkg/resp"
	"user-profiles/internal/middleware"
	"net/http"
)

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[LoginInput](&w, r)
	if err != nil {
		return
	}
	tokens, err := handler.Service.Login(body.Email, body.Password)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Send refresh like HttpOnly cookie
	setCookie(tokens.Refresh, w)
	res := &LoginResponse{
		Token: tokens.Access,
	}
	resp.Json(w, res, http.StatusOK)
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[RegisterInput](&w, r)
	if err != nil {
		return
	}
	res, err := handler.Service.Register(body.Email, body.Password, body.Name)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Json(w, res, http.StatusCreated)
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetUserID(r.Context())
	if err != nil {
		resp.Json(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := getRefreshTokenFromCookie(r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.Service.Logout(userId, token)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}	
	setCookie("", w)
	resp.Json(w, "Logged out", http.StatusOK)
}

func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	token, err := getRefreshTokenFromCookie(r)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := handler.Service.Refresh(token)
	if err != nil {
		resp.Json(w, err.Error(), http.StatusBadRequest)
		return
	}
	setCookie(tokens.Refresh, w)
	res := &LoginResponse{
		Token: tokens.Access,
	}
	resp.Json(w, res, http.StatusOK)
}