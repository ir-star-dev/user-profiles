package auth

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-profiles/configs"
	"user-profiles/internal/http/request"
	"user-profiles/internal/models"
	"user-profiles/internal/templates"
	"user-profiles/internal/utils"

	"github.com/go-chi/chi/v5"
)

type AHandlerDeps struct {
	Config     *configs.Config
	AService   AuthService
	JWTService JWTService
	*templates.BaseHandler
}

type AHandler struct {
	Config     *configs.Config
	AService   AuthService
	JWTService JWTService
	*templates.BaseHandler
}

func NewAuthHandler(router chi.Router, deps AHandlerDeps) *AHandler {
	return &AHandler{
		Config:     deps.Config,
		AService:   deps.AService,
		JWTService: deps.JWTService,
		BaseHandler: deps.BaseHandler,
	}
}

func (h *AHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	data := h.NewBasePageData(r)
	h.BaseHandler.TCache.Render(w, r, http.StatusOK, "base", "login.tmpl", data)
}

func (h *AHandler) SignupForm(w http.ResponseWriter, r *http.Request) {
	data := h.NewBasePageData(r)
	h.BaseHandler.TCache.Render(w, r, http.StatusOK, "base", "signup.tmpl", data)
}

func (h *AHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	base := h.NewBasePageData(r)
	data, formValidErr, err := h.AService.Login(email, password)
	if err != nil {
		data := models.AuthData{
			FormValidationErr: formValidErr,
			BasePageData: base,
		}
		err = h.BaseHandler.TCache.RenderPartial(w, "login.tmpl", "form-submit-error", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	}
	request.SetCookie(data.Access, "__up_access_token", 5*time.Minute, w)
	request.SetCookie(data.Refresh, "__up_refresh_token", 7*24*time.Hour, w)

	agent := r.UserAgent()
	ip := r.RemoteAddr
	device := utils.DetectDevice(agent)

	h.AService.CreateLoginLog(data.UserId, agent, ip, device)

	w.Header().Set("HX-Redirect", "/panel/profile/"+strconv.Itoa(data.UserId))
}

func (h *AHandler) Signup(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	name := strings.TrimSpace(r.FormValue("name"))
	role := "user"
	base := h.NewBasePageData(r)
	formValidErr, err := h.AService.Register(email, password, name, role)
	if err != nil {
		data := models.AuthData{
			FormValidationErr: formValidErr,
			BasePageData: base,
		}
		err := h.BaseHandler.TCache.RenderPartial(w, "signup.tmpl", "form-submit-error", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	} else {
		w.Header().Set("HX-Redirect", "/auth/login")
	}
}

func (h *AHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := request.GetUserID(r.Context())
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	token, err := request.GetCookie("__up_refresh_token", r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	err = h.AService.Logout(userId, token)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	request.SetCookie("", "__up_access_token", -time.Minute, w)
	request.SetCookie("", "__up_refresh_token", -time.Hour, w)

	w.Header().Set("HX-Redirect", "/auth/login")
}
