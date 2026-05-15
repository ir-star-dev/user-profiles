package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/internal/http/cookie"

	"github.com/go-chi/chi/v5"
)

type DashboardHandler struct {
	UsersService users.UsersService
	PostsService posts.PostService
	TCache       panel.Templates
	Dashboard    panel.DS
}

type DashboardHandlerDeps struct {
	UsersService users.UsersService
	PostsService posts.PostService
	TCache       panel.Templates
	Dashboard    panel.DS
}

func NewDashboardHandler(router chi.Router, deps DashboardHandlerDeps) *DashboardHandler {
	return &DashboardHandler{
		UsersService: deps.UsersService,
		PostsService: deps.PostsService,
		TCache:       deps.TCache,
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
		handler.TCache.ServerError(w, r, err)
	}
	data := panel.PageData{
		RequestedUserId: uIdFromReq,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		UserCards:       cards,
		Stats:           stats,
	}

	handler.TCache.RenderPanel(w, r, http.StatusOK, "profile.tmpl", data)
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
	pagination := handler.TCache.BuildPagination(page, totalPages, "/panel/users")
	roles, _ := handler.Dashboard.GetRoles()
	data := panel.PageData{
		UserCards:       cards,
		Pagination:      pagination,
		CurrentUserId:   currUserId,
		CurrentUserRole: currRole,
		Roles:           roles,
		Filters: map[string]string{
			"banned": banned,
			"role":   role,
		},
		HasFilters: page > 1 || banned != "" || role != "",
	}
	handler.TCache.RenderPanel(w, r, http.StatusOK, "users.tmpl", data)
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
	pagination := handler.TCache.BuildPagination(page, totalPages, "/panel/posts")
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
	handler.TCache.RenderPanel(w, r, http.StatusOK, "posts.tmpl", data)
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
	err = handler.TCache.RenderPartial(w, "profile.tmpl", "update-name-modal", data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	res, err := handler.UsersService.ChangeName(uId, name)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
	w.Header().Set("HX-Redirect", "/panel/profile/"+*res)
}

func (handler *DashboardHandler) Ban(w http.ResponseWriter, r *http.Request) {
	reqId, currId, err := panel.SelfDeletionDetected(r)
	if err != nil {
		return
	}
	v := panel.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	err = handler.UsersService.Ban(reqId)
	if err != nil {
		return
	}
	profile, err := handler.UsersService.Get(reqId)
	if err != nil {
		return
	}
	currRole, err := handler.UsersService.Role(currId)
	if err != nil {
		return
	}
	data := panel.UserData{
		Profiles:        *profile,
		CurrentUserId:   currId,
		CurrentUserRole: currRole,
		Actions: panel.Actions{
			CanDeleteUser: panel.CanDeleteUser(currRole, currId, reqId),
			CanBanUser:    panel.CanBanUser(currRole, currId, reqId),
		},
	}
	page := "user-" + v
	err = handler.TCache.RenderPartial(w, "users.tmpl", page, data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) Unban(w http.ResponseWriter, r *http.Request) {
	reqId, currId, err := panel.SelfDeletionDetected(r)
	if err != nil {
		return
	}
	v := panel.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	err = handler.UsersService.Unban(reqId)
	if err != nil {
		return
	}
	profile, err := handler.UsersService.Get(reqId)
	if err != nil {
		return
	}
	currRole, err := handler.UsersService.Role(currId)
	if err != nil {
		return
	}
	data := panel.UserData{
		Profiles:        *profile,
		CurrentUserId:   currId,
		CurrentUserRole: currRole,
		Actions: panel.Actions{
			CanDeleteUser: panel.CanDeleteUser(currRole, currId, reqId),
			CanBanUser:    panel.CanBanUser(currRole, currId, reqId),
		},
	}
	page := "user-" + v
	err = handler.TCache.RenderPartial(w, "users.tmpl", page, data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) DeleteUserConfirm(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetIdFromReq(r)
	if err != nil {
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
	err = handler.TCache.RenderPartial(w, "users.tmpl", "confirm-delete-user-modal", data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	currId, err := panel.GetUserId(r)
	if err != nil {
		w.Header().Set("HX-Redirect", "/auth/login")
		return
	}
	role, _ := handler.UsersService.Role(currId)
	if role == "admin" {
		uId, err := panel.GetIdFromReq(r)
		if err != nil {
			return
		}
		err = handler.UsersService.Delete(uId)
		if err != nil {
			return
		}
		w.Header().Set("HX-Trigger", "userDeleted")
		return
	} else {
		err = handler.UsersService.Delete(currId)
		if err != nil {
			return
		}
		cookie.Set("", "__up_access_token", -time.Minute, w)
		cookie.Set("", "__up_refresh_token", -time.Hour, w)

		w.Header().Set("HX-Redirect", "/auth/login")
		return
	}
}

func (handler *DashboardHandler) Publish(w http.ResponseWriter, r *http.Request) {
	reqUId, err := panel.GetIdFromReq(r)
	if err != nil {
		return
	}
	currUId, err := panel.GetUserId(r)
	if err != nil {
		return
	}
	currUserRole, _ := handler.UsersService.Role(currUId)
	pId, err := panel.GetIdFromReq(r)
	if err != nil {
		return
	}
	v := panel.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	err = handler.PostsService.Publish(pId)
	if err != nil {
		return
	}
	post, err := handler.PostsService.FindById(pId)
	if err != nil {
		return
	}
	data := panel.PostData{
		Posts: *post,
		Actions: panel.Actions{
			CanDeletePost:  panel.CanDeletePost(currUserRole),
			CanEditPost:    panel.CanEditPost(currUserRole, currUId, reqUId),
			CanApprovePost: panel.CanApprovePost(currUserRole),
		},
	}
	page := "post-" + v
	err = handler.TCache.RenderPartial(w, "posts.tmpl", page, data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) Review(w http.ResponseWriter, r *http.Request) {
	reqUId, err := panel.GetIdFromReq(r)
	if err != nil {
		return
	}
	currUId, err := panel.GetUserId(r)
	if err != nil {
		return
	}
	currUserRole, _ := handler.UsersService.Role(currUId)
	pId, err := panel.GetIdFromReq(r)
	if err != nil {
		return
	}
	v := panel.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	err = handler.PostsService.Review(pId)
	if err != nil {
		return
	}
	post, err := handler.PostsService.FindById(pId)
	if err != nil {
		return
	}
	data := panel.PostData{
		Posts: *post,
		Actions: panel.Actions{
			CanDeletePost:  panel.CanDeletePost(currUserRole),
			CanEditPost:    panel.CanEditPost(currUserRole, currUId, reqUId),
			CanApprovePost: panel.CanApprovePost(currUserRole),
		},
	}
	page := "post-" + v
	err = handler.TCache.RenderPartial(w, "posts.tmpl", page, data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) DeletePostConfirm(w http.ResponseWriter, r *http.Request) {
	pId, err := panel.GetIdFromReq(r)
	if err != nil {
		return
	}
	post, err := handler.PostsService.FindById(pId)
	if err != nil {
		return
	}
	var cards []panel.PostData
	cards = append(cards, panel.PostData{
		Posts: *post,
	})
	data := panel.PageData{
		PostCards: cards,
	}
	err = handler.TCache.RenderPartial(w, "posts.tmpl", "confirm-delete-post-modal", data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	pId, err := panel.GetIdFromReq(r)
	if err != nil {
		return
	}
	err = handler.PostsService.Delete(pId)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
	w.Header().Set("HX-Trigger", "postDeleted")
}

func (handler *DashboardHandler) CreateUserForm(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(uId)
	data := panel.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
	}
	handler.TCache.RenderPanel(w, r, http.StatusOK, "create-user.tmpl", data)
}

func (handler *DashboardHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(uId)
	data := panel.PageData{
		CurrentUserId:     uId,
		CurrentUserRole:   currUserRole,
		FormValidationErr: []panel.FormValidationErr{},
	}
	email := strings.TrimSpace(r.FormValue("email"))
	name := strings.TrimSpace(r.FormValue("name"))
	role := strings.TrimSpace(r.FormValue("role"))
	password := strings.TrimSpace(r.FormValue("password"))

	validationErrs, err := handler.Dashboard.CreateUser(email, name, password, role)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := handler.TCache.RenderPartial(w, "create-user.tmpl", "form-submit-error", data)
		if err != nil {
			handler.TCache.ServerError(w, r, err)
		}
		return
	}
	data.UserCredentials = panel.UserCredentials{
		Email:    email,
		Password: password,
	}
	err = handler.TCache.RenderPartial(w, "create-user.tmpl", "form-submit-success", data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	password, err := utils.GeneratePassword(8)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
		return
	}
	w.Write([]byte(password))
}

func (handler *DashboardHandler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(uId)
	data := panel.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
	}
	handler.TCache.RenderPanel(w, r, http.StatusOK, "create-post.tmpl", data)
}

func (handler *DashboardHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(uId)
	data := panel.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
	}
	title := strings.TrimSpace(r.FormValue("title"))
	excerpt := panel.SanitizeContent(strings.TrimSpace(r.FormValue("excerpt")))
	content := panel.SanitizeContent(strings.TrimSpace(r.FormValue("content")))

	_, validationErrs, err := handler.Dashboard.CreatePost(title, content, excerpt, uId)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := handler.TCache.RenderPartial(w, "create-post.tmpl", "form-submit-error", data)
		if err != nil {
			handler.TCache.ServerError(w, r, err)
		}
		return
	}
	data.PostCreated = panel.PostCreated{
		Message: template.HTML("Post successfully created and waiting for moderation!"),
	}
	err = handler.TCache.RenderPartial(w, "create-post.tmpl", "form-submit-success", data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}

func (handler *DashboardHandler) PreviewPost(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(uId)
	pId, err := panel.GetIdFromReq(r)
	if err != nil {
		handler.TCache.PanelNotFound(w, r)
		return
	}
	post, err := handler.Dashboard.PreviewPost(pId)
	if err != nil {
		handler.TCache.PanelNotFound(w, r)
		return
	}
	var postCards []panel.PostData
	postCards = append(postCards, panel.PostData{
		Posts: *post,
	})
	data := panel.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
		PostCards:       postCards,
	}
	handler.TCache.RenderPanel(w, r, http.StatusOK, "preview-post.tmpl", data)
}

func (handler *DashboardHandler) EditPostForm(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(uId)
	data := panel.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
	}
	pId, err := panel.GetIdFromReq(r)
	if err != nil {
		handler.TCache.NotFound(w, r, data)
		return
	}
	post, err := handler.Dashboard.PreviewPost(pId)
	if err != nil {
		handler.TCache.NotFound(w, r, data)
		return
	}
	var postCard []panel.PostData
	postCard = append(postCard, panel.PostData{
		Posts: *post,
	})
	data.PostCards = postCard
	handler.TCache.RenderPanel(w, r, http.StatusOK, "edit-post.tmpl", data)
}

func (handler *DashboardHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	uId, err := panel.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := handler.UsersService.Role(uId)
	pId, err := panel.GetIdFromReq(r)
	if err != nil {
		handler.TCache.PanelNotFound(w, r)
		return
	}
	data := panel.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
	}
	title := strings.TrimSpace(r.FormValue("title"))
	content := panel.SanitizeContent(strings.TrimSpace(r.FormValue("content")))
	excerpt := panel.SanitizeContent(strings.TrimSpace(r.FormValue("excerpt")))
	createdAt := r.FormValue("created_at")

	_, validationErrs, err := handler.Dashboard.EditPost(title, content, excerpt, createdAt, pId)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := handler.TCache.RenderPartial(w, "edit-post.tmpl", "form-submit-error", data)
		if err != nil {
			handler.TCache.ServerError(w, r, err)
		}
		return
	}
	data.PostUpdated = panel.PostUpdated{
		Message: template.HTML("<p>Post updated and waiting for moderation!</p>"),
	}
	err = handler.TCache.RenderPartial(w, "edit-post.tmpl", "form-submit-success", data)
	if err != nil {
		handler.TCache.ServerError(w, r, err)
	}
}
