package handlers

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"html/template"
	//"time"
	"user-profiles/cmd/view"
	"user-profiles/configs"

	//"user-profiles/internal/http/cookie"
	"user-profiles/internal/http/cookie"
	"user-profiles/internal/http/middleware"

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

type UsersPageData struct {
	RequestedUserId int
	Cards           []UserCardData
}
type UserCardData struct {
	CurrentUserId   int
	CurrentUserRole string
	Profile         users.UserWithRole
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
	uId, err := getIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	profile, err := handler.UsersService.View(uId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/profile.tmpl",
		"././ui/templates/parts/layout/user.tmpl",
		"././ui/templates/parts/layout/user-card.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	currUserId, err := getUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(currUserId)

	var cards []UserCardData
	cards = append(cards, UserCardData{
		Profile:         *profile,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
	})
	data := UsersPageData{
		RequestedUserId: uId,
		Cards:           cards,
	}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *UserHandler) UsersProfilePage(w http.ResponseWriter, r *http.Request) {
	uId, err := getUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	var tmpl *template.Template
	tmpl, err = view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/users-profile.tmpl",
		"././ui/templates/parts/layout/users.tmpl",
		"././ui/templates/parts/layout/user-card.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profiles, err := handler.UsersService.ViewAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cards []UserCardData

	for _, profile := range profiles {
		cards = append(cards, UserCardData{
			Profile:       profile,
			CurrentUserId: uId,
		})
	}

	data := UsersPageData{
		Cards: cards,
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(buf.Bytes())
}

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
	reqId, currId, err := selfDeletionDetected(r)
	if err != nil {
		w.Header().Set("HX-Retarget", "#flash-message")
		w.Write([]byte(`
			<div class="alert alert-danger" role="alert">
				You can't delete yourself!
			</div>
		`))
		return
	}
	role, err := handler.UsersService.Role(currId)
	if err != nil {
		return
	}
	err = handler.UsersService.Delete(reqId)
	if err != nil {
		return
	}
	if role == "admin" {
		w.Header().Set("HX-Trigger", "userDeleted")
		w.WriteHeader(http.StatusOK)
		return
	}
	cookie.Set("", "__up_access_token", -time.Minute, w)
	cookie.Set("", "__up_refresh_token", -time.Hour, w)

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}

// func (handler *UserHandler) Ban(w http.ResponseWriter, r *http.Request) {
// 	reqId, currId, err := selfDeletionDetected(r)
// 	if err != nil {
// 		return
// 	}
// 	tmpl, err := view.LoadTemplate(
// 		"./ui/templates/parts/layout/user-card.tmpl",
// 	)
// 	if err != nil {
// 		return
// 	}
// 	err = handler.UsersService.Ban(reqId)
// 	if err != nil {
// 		return
// 	}
// 	user, err := handler.UsersService.View(reqId)
// 	if err != nil {
// 		return
// 	}
// 	data := UsersPageData{
// 		Profile:         user,
// 		CurrentUserId:   currId,
// 		RequestedUserId: reqId,
// 	}
// 	tmpl.ExecuteTemplate(w, "user-card", data)
// }

// func (handler *UserHandler) Unban(w http.ResponseWriter, r *http.Request) {
// 	reqId, currId, err := selfDeletionDetected(r)
// 	if err != nil {
// 		return
// 	}
// 	tmpl, err := view.LoadTemplate(
// 		"./ui/templates/parts/layout/user-card.tmpl",
// 	)
// 	if err != nil {
// 		return
// 	}
// 	err = handler.UsersService.Unban(reqId)
// 	if err != nil {
// 		return
// 	}
// 	user, err := handler.UsersService.View(reqId)
// 	if err != nil {
// 		return
// 	}
// 	data := UsersPageData{
// 		Profile:         user,
// 		CurrentUserId:   currId,
// 		RequestedUserId: reqId,
// 	}
// 	tmpl.ExecuteTemplate(w, "user-card", data)
// }

func getUserId(r *http.Request) (int, error) {
	uId, err := middleware.GetUserID(r.Context())
	if err != nil {
		return 0, err
	}
	return uId, nil
}

func getIdFromReq(r *http.Request) (int, error) {
	idStr := strings.TrimSpace(r.PathValue("id"))
	if idStr == "" {
		return 0, errors.New("Missing param")
	}
	uId, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uId, nil
}

func selfDeletionDetected(r *http.Request) (int, int, error) {
	reqId, err := getIdFromReq(r)
	if err != nil {
		return 0, 0, err
	}
	currId, err := getUserId(r)
	if err != nil {
		return 0, 0, err
	}
	if reqId == currId {
		return 0, 0, errors.New("You can't delete yourself.")
	}
	return reqId, currId, nil
}
