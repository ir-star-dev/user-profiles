package handlers

import (

	//"errors"
	//"html/template"
	"net/http"
	// "strconv"
	// "strings"
	// "time"

	// "html/template"
	// "time"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/users"
	"user-profiles/cmd/user-profiles/view"
	"user-profiles/configs"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	Config       *configs.Config
	AuthService  auth.AuthService
	UsersService users.UsersService
	JWTService   auth.JWTService
	Templates    view.Templates
}

type UserHandlerDeps struct {
	Config       *configs.Config
	AuthService  auth.AuthService
	UsersService users.UsersService
	JWTService   auth.JWTService
	Templates    view.Templates
}

func NewUserHandler(router chi.Router, deps UserHandlerDeps) *UserHandler {
	return &UserHandler{
		Config:       deps.Config,
		AuthService:  deps.AuthService,
		UsersService: deps.UsersService,
		JWTService:   deps.JWTService,
		Templates:    deps.Templates,
	}
}

func (handler *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	uIdFromReq, err := view.GetIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	profile, err := handler.UsersService.Get(uIdFromReq)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	currUserId, err := view.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(currUserId)

	var cards []view.UserData
	cards = append(cards, view.UserData{
		Profiles:        *profile,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		CanDelete:       view.CanDelete(currUserRole, currUserId, profile.Id),
		CanBan:          view.CanBan(currUserRole, currUserId, profile.Id),
	})
	data := view.PageData{
		RequestedUserId: uIdFromReq,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		UserCards:       cards,
	}

	err = handler.Templates.Render(w, "profile", data,
		"././ui/pages/panel/profile.tmpl",
		"././ui/parts/layout/panel/user.tmpl",
		"././ui/parts/layout/panel/user-card.tmpl",
		"././ui/parts/layout/panel/delete-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *UserHandler) statis() (*view.Stats, error) {
	users, err := handler.UsersService.GetAll()
	if err != nil {
		return nil, err
	}
	posts, err := handler.
	if err != nil {
		return nil, err
	}
}

// func (handler *UserHandler) UserListPage(w http.ResponseWriter, r *http.Request) {
// 	page, err := getPageFromReq(r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if page <= 0 {
// 		http.Error(w, "Page does not exist", http.StatusInternalServerError)
// 		return
// 	}

// 	uId, err := getUserId(r)
// 	if err != nil {
// 		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
// 		return
// 	}
// 	currRole, err := handler.UsersService.Role(uId)
// 	if err != nil {
// 		return
// 	}
// 	var tmpl *template.Template
// 	tmpl, err = view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/users-profile.tmpl",
// 		"././ui/templates/parts/layout/users.tmpl",
// 		"././ui/templates/parts/layout/pagination.tmpl",
// 		"././ui/templates/parts/layout/user-row.tmpl",
// 		"././ui/templates/parts/layout/delete-modal.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	profiles, res, err := handler.UsersService.ViewAll(page)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	var cards []view.UserCardData
// 	for _, profile := range profiles {
// 		cards = append(cards, view.UserCardData{
// 			Profiles:  profile,
// 			CanDelete: canDelete(currRole, uId, profile.Id),
// 			CanBan:    canBan(currRole, uId, profile.Id),
// 		})
// 	}

// 	totalPages := (res + 10 - 1) / 10

// 	pages := []int{}
// 	for i := 1; i <= totalPages; i++ {
// 		pages = append(pages, i)
// 	}
// 	pagination := BuildPagination(page, totalPages)
// 	data := view.PageData{
// 		Cards:           cards,
// 		Pagination:      pagination,
// 		CurrentUserId:   uId,
// 		CurrentUserRole: currRole,
// 	}
// 	var buf bytes.Buffer
// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	w.Write(buf.Bytes())
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

// func (handler *UserHandler) DeleteConfirm(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getIdFromReq(r)
// 	if err != nil {
// 		return
// 	}
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/parts/layout/delete-modal.tmpl",
// 	)
// 	if err != nil {
// 		return
// 	}
// 	profile, err := handler.UsersService.View(uId)
// 	if err != nil {
// 		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
// 		return
// 	}
// 	var cards []view.UserCardData
// 	cards = append(cards, view.UserCardData{
// 		Profiles: *profile,
// 	})
// 	data := view.PageData{
// 		Cards: cards,
// 	}
// 	tmpl.ExecuteTemplate(w, "delete-modal", data)
// }

// func (handler *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	currId, err := getUserId(r)
// 	if err != nil {
// 		w.Header().Set("HX-Redirect", "/auth/login")
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
// 	role, err := handler.UsersService.Role(currId)
// 	if err != nil {
// 		return
// 	}
// 	if role == "admin" {
// 		w.Header().Set("HX-Trigger", "userDeleted")
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	} else {
// 		cookie.Set("", "__up_access_token", -time.Minute, w)
// 		cookie.Set("", "__up_refresh_token", -time.Hour, w)

// 		w.Header().Set("HX-Redirect", "/auth/login")
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
// }

// func (handler *UserHandler) Ban(w http.ResponseWriter, r *http.Request) {
// 	reqId, currId, err := selfDeletionDetected(r)
// 	if err != nil {
// 		return
// 	}
// 	v := r.URL.Query().Get("view")
// 	if v == "" {
// 		return
// 	}
// 	tmpl, err := view.LoadTemplate(
// 		"./ui/templates/parts/layout/user-" + v + ".tmpl",
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
// 	currRole, err := handler.UsersService.Role(currId)
// 	if err != nil {
// 		return
// 	}
// 	data := view.UserCardData{
// 		Profiles:      *user,
// 		CurrentUserId: currId,
// 		CanDelete:     canDelete(currRole, currId, reqId),
// 		CanBan:        canBan(currRole, currId, reqId),
// 	}
// 	tmpl.ExecuteTemplate(w, "user-"+v, data)
// }

// func (handler *UserHandler) Unban(w http.ResponseWriter, r *http.Request) {
// 	reqId, currId, err := selfDeletionDetected(r)
// 	if err != nil {
// 		return
// 	}
// 	v := r.URL.Query().Get("view")
// 	if v == "" {
// 		return
// 	}
// 	tmpl, err := view.LoadTemplate(
// 		"./ui/templates/parts/layout/user-" + v + ".tmpl",
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
// 	currRole, err := handler.UsersService.Role(currId)
// 	if err != nil {
// 		return
// 	}
// 	data := view.UserCardData{
// 		Profiles:      *user,
// 		CurrentUserId: currId,
// 		CanDelete:     canDelete(currRole, currId, reqId),
// 		CanBan:        canBan(currRole, currId, reqId),
// 	}
// 	tmpl.ExecuteTemplate(w, "user-"+v, data)
// }
