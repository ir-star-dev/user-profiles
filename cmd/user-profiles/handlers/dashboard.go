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
	"user-profiles/cmd/user-profiles/app"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
	"user-profiles/configs"

	"github.com/go-chi/chi/v5"
)

type DashboardHandler struct {
	Config       *configs.Config
	AuthService  auth.AuthService
	UsersService users.UsersService
	PostsService posts.PostService
	JWTService   auth.JWTService
	Templates    app.Templates
	Dashboard    app.DS
}

type DashboardHandlerDeps struct {
	Config       *configs.Config
	AuthService  auth.AuthService
	UsersService users.UsersService
	PostsService posts.PostService
	JWTService   auth.JWTService
	Templates    app.Templates
	Dashboard    app.DS
}

func NewDashboardHandler(router chi.Router, deps DashboardHandlerDeps) *DashboardHandler {
	return &DashboardHandler{
		Config:       deps.Config,
		AuthService:  deps.AuthService,
		UsersService: deps.UsersService,
		PostsService: deps.PostsService,
		JWTService:   deps.JWTService,
		Templates:    deps.Templates,
		Dashboard:    deps.Dashboard,
	}
}

func (handler *DashboardHandler) Profile(w http.ResponseWriter, r *http.Request) {
	uIdFromReq, err := app.GetIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	profile, err := handler.UsersService.Get(uIdFromReq)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	currUserId, err := app.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(currUserId)

	var cards []app.UserData
	cards = append(cards, app.UserData{
		Profiles:        *profile,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		CanDelete:       app.CanDelete(currUserRole, currUserId, profile.Id),
		CanBan:          app.CanBan(currUserRole, currUserId, profile.Id),
	})

	stats, err := handler.Dashboard.GetStats()
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
	data := app.PageData{
		RequestedUserId: uIdFromReq,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		UserCards:       cards,
		Stats:           stats,
	}

	err = handler.Templates.Render(w, "profile", data,
		"././ui/pages/panel/profile.tmpl",
		"././ui/parts/layout/panel/user.tmpl",
		"././ui/parts/layout/panel/user-card.tmpl",
		"././ui/parts/layout/panel/stats.tmpl",
		"././ui/parts/layout/panel/delete-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *DashboardHandler) Users(w http.ResponseWriter, r *http.Request) {
	page, err := app.GetPageFromReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if page <= 0 {
		page = 1
	}
	currUserId, err := app.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, err := handler.UsersService.Role(currUserId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	profiles, res, err := handler.UsersService.GetOnPage(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cards []app.UserData
	for _, profile := range profiles {
		cards = append(cards, app.UserData{
			Profiles:  profile,
			CanDelete: app.CanDelete(currRole, currUserId, profile.Id),
			CanBan:    app.CanBan(currRole, currUserId, profile.Id),
		})
	}

	totalPages := (res + 10 - 1) / 10

	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := handler.Templates.BuildPagination(page, totalPages)
	data := app.PageData{
		UserCards:       cards,
		Pagination:      pagination,
		CurrentUserId:   currUserId,
		CurrentUserRole: currRole,
	}
	err = handler.Templates.Render(w, "users", data,
		"././ui/pages/panel/users.tmpl",
		"././ui/parts/layout/panel/user-table.tmpl",
		"././ui/parts/layout/panel/user-pagin.tmpl",
		"././ui/parts/layout/panel/user-row.tmpl",
		"././ui/parts/layout/panel/delete-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *DashboardHandler) Posts(w http.ResponseWriter, r *http.Request) {
	page, err := app.GetPageFromReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if page <= 0 {
		page = 1
	}
	currUserId, err := app.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, err := handler.UsersService.Role(currUserId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	posts, res, err := handler.PostsService.GetOnPage(page, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cards []app.PostData
	for _, post := range posts {
		cards = append(cards, app.PostData{
			Posts:          post,
			CanDeletePost:  app.CanDeletePost(currRole),
			CanApprovePost: app.CanApprovePost(currRole),
			CanEditPost:    app.CanEditPost(currRole, currUserId, post.UserId),
		})
	}

	totalPages := (res + 10 - 1) / 10

	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := handler.Templates.BuildPagination(page, totalPages)
	data := app.PageData{
		PostCards:       cards,
		Pagination:      pagination,
		CurrentUserId:   currUserId,
		CurrentUserRole: currRole,
	}
	err = handler.Templates.Render(w, "posts", data,
		"././ui/pages/panel/posts.tmpl",
		"././ui/parts/layout/panel/post-table.tmpl",
		"././ui/parts/layout/panel/post-pagin.tmpl",
		"././ui/parts/layout/panel/post-row.tmpl",
		"././ui/parts/layout/panel/delete-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
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

// func (handler *UserHandler) DeleteConfirm(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getIdFromReq(r)
// 	if err != nil {
// 		return
// 	}
// 	tmpl, err := app.LoadTemplate(
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
// 	var cards []app.UserCardData
// 	cards = append(cards, app.UserCardData{
// 		Profiles: *profile,
// 	})
// 	data := app.PageData{
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
// 	tmpl, err := app.LoadTemplate(
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
// 	data := app.UserCardData{
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
// 	tmpl, err := app.LoadTemplate(
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
// 	data := app.UserCardData{
// 		Profiles:      *user,
// 		CurrentUserId: currId,
// 		CanDelete:     canDelete(currRole, currId, reqId),
// 		CanBan:        canBan(currRole, currId, reqId),
// 	}
// 	tmpl.ExecuteTemplate(w, "user-"+v, data)
// }
