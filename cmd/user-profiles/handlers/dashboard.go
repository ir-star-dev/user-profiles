package handlers

import (
	//"errors"
	//"html/template"
	"net/http"
	"strconv"

	// "strconv"
	// "strings"
	// "time"

	// "html/template"
	// "time"

	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"

	"github.com/go-chi/chi/v5"
)

type DashboardHandler struct {
	UsersService users.UsersService
	PostsService posts.PostService
	Templates    panel.Templates
	Dashboard    panel.DS
}

type DashboardHandlerDeps struct {
	UsersService users.UsersService
	PostsService posts.PostService
	Templates    panel.Templates
	Dashboard    panel.DS
}

func NewDashboardHandler(router chi.Router, deps DashboardHandlerDeps) *DashboardHandler {
	return &DashboardHandler{
		UsersService: deps.UsersService,
		PostsService: deps.PostsService,
		Templates:    deps.Templates,
		Dashboard:    deps.Dashboard,
	}
}

func (handler *DashboardHandler) Profile(w http.ResponseWriter, r *http.Request) {
	uIdFromReq, err := panel.GetIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	profile, err := handler.UsersService.Get(uIdFromReq)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	currUserId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(currUserId)

	var cards []panel.UserData
	cards = append(cards, panel.UserData{
		Profiles:        *profile,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		Actions: panel.Actions{
			CanDeleteUser: panel.CanDeleteUser(currUserRole, currUserId, profile.Id),
			CanBanUser:    panel.CanBanUser(currUserRole, currUserId, profile.Id),
		},
	})

	stats, err := handler.Dashboard.GetStats()
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
	data := panel.PageData{
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
		"././ui/parts/layout/panel/modals/delete-modal.tmpl",
		"././ui/parts/layout/panel/modals/update-name-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *DashboardHandler) Users(w http.ResponseWriter, r *http.Request) {
	page := panel.GetPageFromReq(r)
	role := panel.GetFilterValue(r, "role")
	banned := panel.GetFilterValue(r, "banned")

	currUserId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, err := handler.UsersService.Role(currUserId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	var bn *bool
	if banned != "" {
		v, err := strconv.ParseBool(banned)
		if err != nil {
			bn = nil
		}
		bn = &v
	}

	profiles, res, err := handler.UsersService.GetOnPage(page, 10, bn, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cards []panel.UserData
	for _, profile := range profiles {
		cards = append(cards, panel.UserData{
			Profiles: profile,
			Actions: panel.Actions{
				CanDeleteUser: panel.CanDeleteUser(currRole, currUserId, profile.Id),
				CanBanUser:    panel.CanBanUser(currRole, currUserId, profile.Id),
			},
		})
	}

	totalPages := (res + 10 - 1) / 10

	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := handler.Templates.BuildPagination(page, totalPages, "/panel/users")
	roles, _ := handler.Dashboard.GetRoles()
	data := panel.PageData{
		UserCards:       cards,
		Pagination:      pagination,
		CurrentUserId:   currUserId,
		CurrentUserRole: currRole,
		Roles:           roles,
		Filters: map[string]string{
			"banned": banned,
			"role": role,
		},
		HasFilters: page > 1 || banned != "" || role != "",
	}
	err = handler.Templates.Render(w, "users", data,
		"././ui/pages/panel/users.tmpl",
		"././ui/parts/layout/panel/user-table.tmpl",
		"././ui/parts/layout/panel/filters/user-filter.tmpl",
		"././ui/parts/layout/panel/post-pagin.tmpl",
		"././ui/parts/layout/panel/user-row.tmpl",
		"././ui/parts/layout/panel/modals/delete-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *DashboardHandler) Posts(w http.ResponseWriter, r *http.Request) {
	page := panel.GetPageFromReq(r)
	approved := panel.GetFilterValue(r, "approved")
	username := panel.GetFilterValue(r, "username")

	currUserId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, err := handler.UsersService.Role(currUserId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	var ap *bool
	if approved != "" {
		v, err := strconv.ParseBool(approved)
		if err != nil {
			ap = nil
		}
		ap = &v
	}
	posts, res, err := handler.PostsService.GetOnPage(page, 10, ap, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authors, _ := handler.Dashboard.GetAuthors()

	var cards []panel.PostData
	for _, post := range posts {
		cards = append(cards, panel.PostData{
			Posts: post,
			Actions: panel.Actions{
				CanDeletePost:  panel.CanDeletePost(currRole),
				CanApprovePost: panel.CanApprovePost(currRole),
				CanEditPost:    panel.CanEditPost(currRole, currUserId, post.UserId),
			},
		})
	}

	totalPages := (res + 10 - 1) / 10

	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := handler.Templates.BuildPagination(page, totalPages, "/panel/posts")
	data := panel.PageData{
		PostCards:       cards,
		Pagination:      pagination,
		CurrentUserId:   currUserId,
		CurrentUserRole: currRole,
		Authors:         authors,
		Filters: map[string]string{
			"approved": approved,
			"username": username,
		},
		HasFilters: page > 1 || approved != "" || username != "",
	}
	err = handler.Templates.Render(w, "posts", data,
		"././ui/pages/panel/posts.tmpl",
		"././ui/parts/layout/panel/post-table.tmpl",
		"././ui/parts/layout/panel/filters/post-filter.tmpl",
		"././ui/parts/layout/panel/post-pagin.tmpl",
		"././ui/parts/layout/panel/post-row.tmpl",
		"././ui/parts/layout/panel/modals/delete-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *DashboardHandler) UpdateNameModal(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	profile, err := handler.UsersService.Get(uId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	var cards []panel.UserData
	cards = append(cards, panel.UserData{
		Profiles: *profile,
	})
	data := panel.PageData{
		UserCards: cards,
	}
	err = handler.Templates.RenderPartial(
		w,
		"update-name-modal",
		data,
		"././ui/parts/layout/panel/modals/update-name-modal.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *DashboardHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	name := r.FormValue("name")
	res, err := handler.UsersService.ChangeName(uId, name)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
	w.Header().Set("HX-Redirect", "/panel/profile/"+*res)
}

// func (handler *UserHandler) DeleteConfirm(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getIdFromReq(r)
// 	if err != nil {
// 		return
// 	}
// 	tmpl, err := panel.LoadTemplate(
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
// 	var cards []panel.UserCardData
// 	cards = append(cards, panel.UserCardData{
// 		Profiles: *profile,
// 	})
// 	data := panel.PageData{
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
// 	tmpl, err := panel.LoadTemplate(
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
// 	data := panel.UserCardData{
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
// 	tmpl, err := panel.LoadTemplate(
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
// 	data := panel.UserCardData{
// 		Profiles:      *user,
// 		CurrentUserId: currId,
// 		CanDelete:     canDelete(currRole, currId, reqId),
// 		CanBan:        canBan(currRole, currId, reqId),
// 	}
// 	tmpl.ExecuteTemplate(w, "user-"+v, data)
// }
