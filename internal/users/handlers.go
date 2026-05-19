package users

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-profiles/internal/http/request"
	"user-profiles/internal/models"
	"user-profiles/internal/templates"
	"user-profiles/internal/utils"
	"user-profiles/internal/validator"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

type UHandler struct {
	UService UsersService
	TCache   templates.Templates
}

type UHandlerDeps struct {
	UService UsersService
	TCache   templates.Templates
}

func NewUserHandler(router chi.Router, deps UHandlerDeps) *UHandler {
	return &UHandler{
		UService: deps.UService,
		TCache:   deps.TCache,
	}
}

func (h *UHandler) Users(w http.ResponseWriter, r *http.Request) {
	page := request.GetPageFromReq(r)
	role := request.GetFilterValue(r, "role")
	banned := request.GetFilterValue(r, "banned")
	search := strings.TrimSpace(request.GetFilterValue(r, "search"))

	currUserId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, _ := h.UService.Role(currUserId)
	var bn *bool
	if banned != "" {
		v, err := strconv.ParseBool(banned)
		if err != nil {
			bn = nil
		}
		bn = &v
	}
	data := models.PageData{
		CurrentUserId:   currUserId,
		CurrentUserRole: currRole,
		CSRFToken:       csrf.Token(r),
	}
	profiles, totalUsers, err := h.UService.GetOnPage(page, 10, bn, role, search)
	if err != nil {
		h.TCache.ServerError(w, r, err)
		return
	}
	var cards []models.UserData
	for i, profile := range profiles {
		cards = append(cards, models.UserData{
			Profiles: models.UserViewTable{
				User:  profile,
				Index: ((page - 1) * 10) + i + 1,
			},
			Actions: models.Actions{
				CanDeleteUser: CanDeleteUser(currRole, currUserId, profile.Id),
				CanBanUser:    CanBanUser(currRole, currUserId, profile.Id),
			},
		})
	}

	totalPages := (totalUsers + 10 - 1) / 10

	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := h.TCache.BuildPagination(page, totalPages, "/panel/users")
	roles, _ := h.UService.GetRoles()

	data.UserCards = cards
	data.Pagination = pagination
	data.Roles = roles
	data.Filters = map[string]string{
		"banned": banned,
		"role":   role,
		"search": search,
	}
	data.TotalUsers = totalUsers
	data.Page = page
	data.Pages = len(pages)

	data.HasFilters = page > 1 || banned != "" || role != "" || search != ""

	if r.Header.Get("HX-Request") == "true" {
		h.TCache.RenderPartial(w, "users.tmpl", "user-table", data)
		return
	}

	h.TCache.RenderPanel(w, r, http.StatusOK, "users.tmpl", data)
}

func (h *UHandler) UpdateNameModal(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetIdFromReq(r)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}
	profile, err := h.UService.Get(uId)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}

	var cards []models.UserData
	cards = append(cards, models.UserData{
		Profiles: models.UserViewTable{
			User:  *profile,
		},
	})
	data := models.PageData{
		UserCards: cards,
		CSRFToken: csrf.Token(r),
	}
	err = h.TCache.RenderPartial(w, "profile.tmpl", "update-name-modal", data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	_, res, err := h.UService.ChangeName(uId, name)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
	w.Header().Set("HX-Redirect", "/panel/profile/"+*res)
}

func (h *UHandler) Ban(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(request.GetFilterValue(r, "index"))
	currId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	reqId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	_, _, err = SelfDeletionDetected(reqId, currId)
	if err != nil {
		return
	}
	currRole, _ := h.UService.Role(currId)
	v := request.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	data := models.UserData{
		CurrentUserId:   currId,
		CurrentUserRole: currRole,
	}
	err = h.UService.Ban(reqId)
	if err != nil {
		return
	}
	profile, err := h.UService.Get(reqId)
	if err != nil {
		http.Redirect(w, r, "/panel/users", http.StatusSeeOther)
		return
	}
	data.Profiles = models.UserViewTable{
		User: *profile,
		Index: index,
	}
	data.Actions = models.Actions{
		CanDeleteUser: CanDeleteUser(currRole, currId, reqId),
		CanBanUser:    CanBanUser(currRole, currId, reqId),
	}
	page := "user-" + v
	err = h.TCache.RenderPartial(w, "users.tmpl", page, data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) Unban(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(request.GetFilterValue(r, "index"))
	currId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	reqId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	_, _, err = SelfDeletionDetected(reqId, currId)
	if err != nil {
		return
	}
	currRole, _ := h.UService.Role(currId)
	v := request.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	data := models.UserData{
		CurrentUserId:   currId,
		CurrentUserRole: currRole,
	}
	err = h.UService.Unban(reqId)
	if err != nil {
		return
	}
	profile, err := h.UService.Get(reqId)
	if err != nil {
		http.Redirect(w, r, "/panel/users", http.StatusSeeOther)
		return
	}
	data.Profiles = models.UserViewTable{
		User: *profile,
		Index: index,
	}
	data.Actions = models.Actions{
		CanDeleteUser: CanDeleteUser(currRole, currId, reqId),
		CanBanUser:    CanBanUser(currRole, currId, reqId),
	}

	page := "user-" + v
	err = h.TCache.RenderPartial(w, "users.tmpl", page, data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) DeleteUserConfirm(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	profile, err := h.UService.Get(uId)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}
	var cards []models.UserData
	cards = append(cards, models.UserData{
		Profiles: models.UserViewTable{
			User: *profile,
		},
	})
	data := models.PageData{
		UserCards: cards,
		CSRFToken: csrf.Token(r),
	}
	err = h.TCache.RenderPartial(w, "users.tmpl", "confirm-delete-user-modal", data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	currId, err := request.GetUserId(r)
	if err != nil {
		w.Header().Set("HX-Redirect", "/auth/login")
		return
	}
	role, _ := h.UService.Role(currId)
	if role == "admin" {
		uId, err := request.GetIdFromReq(r)
		if err != nil {
			return
		}
		err = h.UService.Delete(uId)
		if err != nil {
			return
		}
		w.Header().Set("HX-Trigger", "userDeleted")
		return
	} else {
		err = h.UService.Delete(currId)
		if err != nil {
			return
		}
		request.SetCookie("", "__up_access_token", -time.Minute, w)
		request.SetCookie("", "__up_refresh_token", -time.Hour, w)

		w.Header().Set("HX-Redirect", "/auth/login")
		return
	}
}

func (h *UHandler) CreateUserForm(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := h.UService.Role(uId)
	data := models.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}
	h.TCache.RenderPanel(w, r, http.StatusOK, "create-user.tmpl", data)
}

func (h *UHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := h.UService.Role(uId)
	data := models.PageData{
		CurrentUserId:     uId,
		CurrentUserRole:   currUserRole,
		FormValidationErr: []validator.FormValidationErr{},
		CSRFToken:         csrf.Token(r),
	}
	email := strings.TrimSpace(r.FormValue("email"))
	name := strings.TrimSpace(r.FormValue("name"))
	role := strings.TrimSpace(r.FormValue("role"))
	password := strings.TrimSpace(r.FormValue("password"))

	validationErrs, err := h.UService.CreateUser(email, name, password, role)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := h.TCache.RenderPartial(w, "create-user.tmpl", "form-submit-error", data)
		if err != nil {
			h.TCache.ServerError(w, r, err)
		}
		return
	}
	data.UserCredentials = models.UserCredentials{
		Email:    email,
		Password: password,
	}
	err = h.TCache.RenderPartial(w, "create-user.tmpl", "form-submit-success", data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	password, err := utils.GeneratePassword(8)
	if err != nil {
		h.TCache.ServerError(w, r, err)
		return
	}
	w.Write([]byte(password))
}
