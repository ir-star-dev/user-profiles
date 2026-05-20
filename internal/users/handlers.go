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
)

type UHandler struct {
	UService UsersService
	*templates.BaseHandler
}

type UHandlerDeps struct {
	UService UsersService
	*templates.BaseHandler
}

func NewUserHandler(router chi.Router, deps UHandlerDeps) *UHandler {
	return &UHandler{
		UService:    deps.UService,
		BaseHandler: deps.BaseHandler,
	}
}

func (h *UHandler) Users(w http.ResponseWriter, r *http.Request) {
	page := request.GetPageFromReq(r)
	role := request.GetFilterValue(r, "role")
	banned := request.GetFilterValue(r, "banned")
	search := strings.TrimSpace(request.GetFilterValue(r, "search"))

	base := h.NewBasePageData(r)

	if base.CurrentUserId == 0 {
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
	roles, _ := h.UService.GetRoles()
	data := models.UsersDashboardData{
		BasePageData: base,
		Filters: map[string]string{
			"banned": banned,
			"role":   role,
			"search": search,
		},
		HasFilters: page > 1 || banned != "" || role != "" || search != "",
		Roles: roles,
	}
	filters := models.UserFilters{
		Role:   role,
		Search: search,
		Banned: bn,
	}
	profiles, totalUsers, err := h.UService.GetOnPage(page, 10, filters)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
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
				CanDeleteUser: CanDeleteUser(base.CurrentUserRole, base.CurrentUserId, profile.Id),
				CanBanUser:    CanBanUser(base.CurrentUserRole, base.CurrentUserId, profile.Id),
			},
		})
	}

	totalPages := (totalUsers + 10 - 1) / 10
	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := h.BaseHandler.TCache.BuildPagination(page, totalPages, "/panel/users")

	data.UserCards = cards
	data.Pagination = pagination
	data.TotalUsers = totalUsers
	data.Page = page
	data.Pages = len(pages)

	if r.Header.Get("HX-Request") == "true" {
		h.BaseHandler.TCache.RenderPartial(w, "users.tmpl", "user-table", data)
		return
	}
	h.BaseHandler.TCache.RenderPanel(w, r, http.StatusOK, "users.tmpl", data)
}

func (h *UHandler) UpdateNameModal(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	uId, err := request.GetIdFromReq(r)
	if err != nil {
		h.BaseHandler.TCache.PanelNotFound(w, r)
		return
	}
	profile, err := h.UService.Get(uId)
	if err != nil {
		h.BaseHandler.TCache.PanelNotFound(w, r)
		return
	}

	var cards []models.UserData
	cards = append(cards, models.UserData{
		Profiles: models.UserViewTable{
			User: *profile,
		},
	})
	data := models.UpdateNameModalData{
		UserCards:    cards,
		BasePageData: base,
	}
	err = h.BaseHandler.TCache.RenderPartial(w, "profile.tmpl", "update-name-modal", data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
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
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
	w.Header().Set("HX-Redirect", "/panel/profile/"+*res)
}

func (h *UHandler) Ban(w http.ResponseWriter, r *http.Request) {
	aData := h.PrepareToRestrictAccess(w, r)
	if aData == nil {
		return
	}
	err := h.UService.Ban(aData.ReqUID)
	if err != nil {
		return
	}
	profile, err := h.UService.Get(aData.ReqUID)
	if err != nil {
		http.Redirect(w, r, "/panel/users", http.StatusSeeOther)
		return
	}
	data := models.UserData{
		Profiles: models.UserViewTable{
			User: *profile,
			Index: aData.Index,
		},
		Actions: models.Actions{
			CanDeleteUser: CanDeleteUser(aData.Data.CurrentUserRole, aData.Data.CurrentUserId, aData.ReqUID),
			CanBanUser:    CanBanUser(aData.Data.CurrentUserRole, aData.Data.CurrentUserId, aData.ReqUID),
		},
	}
	page := "user-" + aData.View
	err = h.BaseHandler.TCache.RenderPartial(w, "users.tmpl", page, data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) Unban(w http.ResponseWriter, r *http.Request) {
	aData := h.PrepareToRestrictAccess(w, r)
	if aData == nil {
		return
	}
	err := h.UService.Unban(aData.ReqUID)
	if err != nil {
		return
	}
	profile, err := h.UService.Get(aData.ReqUID)
	if err != nil {
		http.Redirect(w, r, "/panel/users", http.StatusSeeOther)
		return
	}
	data := models.UserData{
		Profiles: models.UserViewTable{
			User: *profile,
			Index: aData.Index,
		},
		Actions: models.Actions{
			CanDeleteUser: CanDeleteUser(aData.Data.CurrentUserRole, aData.Data.CurrentUserId, aData.ReqUID),
			CanBanUser:    CanBanUser(aData.Data.CurrentUserRole, aData.Data.CurrentUserId, aData.ReqUID),
		},
	}
	page := "user-" + aData.View
	err = h.BaseHandler.TCache.RenderPartial(w, "users.tmpl", page, data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) DeleteUserConfirm(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	uId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	profile, err := h.UService.Get(uId)
	if err != nil {
		h.BaseHandler.TCache.PanelNotFound(w, r)
		return
	}
	var cards []models.UserData
	cards = append(cards, models.UserData{
		Profiles: models.UserViewTable{
			User: *profile,
		},
	})
	data := models.UserDeleteData{
		UserCards:    cards,
		BasePageData: base,
	}
	err = h.BaseHandler.TCache.RenderPartial(w, "users.tmpl", "confirm-delete-user-modal", data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		w.Header().Set("HX-Redirect", "/auth/login")
		return
	}
	reqId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	_, _, err = SelfDeletionDetected(reqId, base.CurrentUserId)
	if err != nil {
		return
	}
	if base.CurrentUserRole == "admin" {
		err = h.UService.Delete(reqId)
		if err != nil {
			return
		}
		w.Header().Set("HX-Trigger", "userDeleted")
		return
	} else {
		err := h.UService.Delete(base.CurrentUserId)
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
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	data := models.CreateUserFormData{
		BasePageData: base,
	}
	h.BaseHandler.TCache.RenderPanel(w, r, http.StatusOK, "create-user.tmpl", data)
}

func (h *UHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	data := models.CreateUserFormData{
		BasePageData:      base,
		FormValidationErr: []validator.FormValidationErr{},
	}
	email := strings.TrimSpace(r.FormValue("email"))
	name := strings.TrimSpace(r.FormValue("name"))
	role := strings.TrimSpace(r.FormValue("role"))
	password := strings.TrimSpace(r.FormValue("password"))

	validationErrs, err := h.UService.CreateUser(email, name, password, role)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := h.BaseHandler.TCache.RenderPartial(w, "create-user.tmpl", "form-submit-error", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	}
	data.Success.UserCredentials = &models.UserCredentials{
		Email:    email,
		Password: password,
	}

	err = h.BaseHandler.TCache.RenderPartial(w, "create-user.tmpl", "form-submit-success", data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *UHandler) GeneratePassword(w http.ResponseWriter, r *http.Request) {
	password, err := utils.GeneratePassword(8)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
		return
	}
	w.Write([]byte(password))
}

func (h *UHandler) PrepareToRestrictAccess(w http.ResponseWriter, r *http.Request) *AccessData {
	index, _ := strconv.Atoi(request.GetFilterValue(r, "index"))
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return nil
	}
	reqId, err := request.GetIdFromReq(r)
	if err != nil {
		return nil
	}
	view := request.GetFilterValue(r, "view")
	if view == "" {
		return nil
	}
	data := models.UserData{
		BasePageData:    base,
		CurrentUserId:   base.CurrentUserId,
		CurrentUserRole: base.CurrentUserRole,
	}
	access := AccessData{
		Index:  index,
		ReqUID: reqId,
		Data:   data,
		View:   view,
	}
	return &access
}
