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

type PageData struct {
	RequestedUserId int
	CurrentUserId   int
	CurrentUserRole string
	Cards           []UserCardData
	Pagination      PaginationData
}

type UserCardData struct {
	Profiles        users.UserWithRole
	CanDelete       bool
	CanBan			bool
	CurrentUserId   int
	CurrentUserRole string
}

type PaginationData struct {
	Page       int
	TotalPages int
	Pages      []int
	HasPrev    bool
	HasNext    bool
	PrevPage   int
	NextPage   int
	ShowDots   bool
	LastPage   int
}

func NewUserHandler(router chi.Router, deps UserHandlerDeps) *UserHandler {
	return &UserHandler{
		Config:       deps.Config,
		AuthService:  deps.AuthService,
		UsersService: deps.UsersService,
		JWTService:   deps.JWTService,
	}
}

func (handler *UserHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
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

	data := PageData{
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
	}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *UserHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	uIdFromReq, err := getIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	profile, err := handler.UsersService.View(uIdFromReq)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/profile.tmpl",
		"././ui/templates/parts/layout/user.tmpl",
		"././ui/templates/parts/layout/user-card.tmpl",
		"././ui/templates/parts/layout/delete-modal.tmpl",
		
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
		Profiles:        *profile,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		CanDelete: canDelete(currUserRole, currUserId, profile.Id),
		CanBan: canBan(currUserRole, currUserId, profile.Id),
	})
	data := PageData{
		RequestedUserId: uIdFromReq,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
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

func (handler *UserHandler) UserListPage(w http.ResponseWriter, r *http.Request) {
	page, err := getPageFromReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if page <= 0 {
		http.Error(w, "Page does not exist", http.StatusInternalServerError)
		return
	}

	uId, err := getUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, err := handler.UsersService.Role(uId)
	if err != nil {
		return
	}
	var tmpl *template.Template
	tmpl, err = view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/users-profile.tmpl",
		"././ui/templates/parts/layout/users.tmpl",
		"././ui/templates/parts/layout/pagination.tmpl",
		"././ui/templates/parts/layout/user-row.tmpl",
		"././ui/templates/parts/layout/delete-modal.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profiles, res, err := handler.UsersService.ViewAll(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cards []UserCardData
	for _, profile := range profiles {
		cards = append(cards, UserCardData{
			Profiles: profile,
			CanDelete: canDelete(currRole, uId, profile.Id),
			CanBan: canBan(currRole, uId, profile.Id),
		})
	}

	totalPages := (res + 10 - 1) / 10

	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := BuildPagination(page, totalPages)
	data := PageData{
		Cards:         cards,
		Pagination:    pagination,
		CurrentUserId: uId,
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

func (handler *UserHandler) DeleteConfirm(w http.ResponseWriter, r *http.Request) {
	uId, err := getIdFromReq(r)
	if err != nil {
		return
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/parts/layout/delete-modal.tmpl",
	)
	if err != nil {
		return
	}
	profile, err := handler.UsersService.View(uId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	var cards []UserCardData
	cards = append(cards, UserCardData{
		Profiles:        *profile,
	})
	data := PageData{
		Cards: cards,
	}
	tmpl.ExecuteTemplate(w, "delete-modal", data)
}

func (handler *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	currId, err := getUserId(r)
	if err != nil {
		w.Header().Set("HX-Redirect", "/auth/login")
		w.WriteHeader(http.StatusOK)
		return
	} 
	role, err := handler.UsersService.Role(currId)
	if err != nil {
		return
	}
	if role == "admin" {
		w.Header().Set("HX-Trigger", "userDeleted")
		w.WriteHeader(http.StatusOK)
		return
	} else {
		cookie.Set("", "__up_access_token", -time.Minute, w)
		cookie.Set("", "__up_refresh_token", -time.Hour, w)

		w.Header().Set("HX-Redirect", "/auth/login")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func BuildPagination(currentPage, totalPages int) PaginationData {
	p := PaginationData{
		Page:       currentPage,
		TotalPages: totalPages,
		HasPrev:    currentPage > 1,
		HasNext:    currentPage < totalPages,
		PrevPage:   currentPage - 1,
		NextPage:   currentPage + 1,
		LastPage:   totalPages,
	}

	start := currentPage - 1
	if start < 1 {
		start = 1
	}

	end := start + 2
	if end > totalPages {
		end = totalPages
		start = end - 2
		if start < 1 {
			start = 1
		}
	}

	for i := start; i <= end; i++ {
		p.Pages = append(p.Pages, i)
	}

	if end < totalPages {
		p.ShowDots = true
	}

	return p
}

func (handler *UserHandler) Ban(w http.ResponseWriter, r *http.Request) {
	reqId, currId, err := selfDeletionDetected(r)
	if err != nil {
		return
	}
	v := r.URL.Query().Get("view")
	if v == "" {
		return
	}
	tmpl, err := view.LoadTemplate(
		"./ui/templates/parts/layout/user-" + v + ".tmpl",
	)
	if err != nil {
		return
	}
	err = handler.UsersService.Ban(reqId)
	if err != nil {
		return
	}
	user, err := handler.UsersService.View(reqId)
	if err != nil {
		return
	}
	currRole, err := handler.UsersService.Role(currId)
	if err != nil {
		return
	}
	data := UserCardData{
		Profiles:      *user,
		CurrentUserId: currId,
		CanDelete: canDelete(currRole, currId, reqId),
		CanBan: canBan(currRole, currId, reqId),
	}
	tmpl.ExecuteTemplate(w, "user-"+v, data)
}

func (handler *UserHandler) Unban(w http.ResponseWriter, r *http.Request) {
	reqId, currId, err := selfDeletionDetected(r)
	if err != nil {
		return
	}
	v := r.URL.Query().Get("view")
	if v == "" {
		return
	}
	tmpl, err := view.LoadTemplate(
		"./ui/templates/parts/layout/user-" + v + ".tmpl",
	)
	if err != nil {
		return
	}
	err = handler.UsersService.Unban(reqId)
	if err != nil {
		return
	}
	user, err := handler.UsersService.View(reqId)
	if err != nil {
		return
	}
	currRole, err := handler.UsersService.Role(currId)
	if err != nil {
		return
	}
	data := UserCardData{
		Profiles:      *user,
		CurrentUserId: currId,
		CanDelete: canDelete(currRole, currId, reqId),
		CanBan: canBan(currRole, currId, reqId),
	}
	tmpl.ExecuteTemplate(w, "user-"+v, data)
}

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

func getPageFromReq(r *http.Request) (int, error) {
	page := strings.TrimSpace(r.PathValue("page"))
	if page == "" {
		return 0, errors.New("Missing param")
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return 0, err
	}
	return p, nil
}

func selfDeletionDetected(r *http.Request) (int, int, error) {
	reqId, err := getIdFromReq(r)
	if err != nil {
		return 0, 0, err
	}
	currId, err := getUserId(r)
	if err != nil {
		return reqId, currId, err
	}
	if reqId == currId {
		return reqId, currId, errors.New("You can't delete yourself.")
	}
	return reqId, currId, nil
}

func canDelete(currentRole string, currentUserId, profileId int) bool {
	if currentRole == "admin" {
		return currentUserId != profileId
	}
	return currentUserId == profileId
}

func canBan(currentRole string, currentUserId, profileId int) bool {
	return currentRole == "admin" && currentUserId != profileId
}