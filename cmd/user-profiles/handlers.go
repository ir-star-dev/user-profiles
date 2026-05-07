package main

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/cmd/view"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/cookie"
	"user-profiles/internal/http/req"
)

func (app *application) Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	currUserId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := app.UsersService.Role(currUserId)

	data := templateData{
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

func (app *application) LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/login.tmpl",
		"././ui/templates/parts/auth/form-login.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) RegisterForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/register.tmpl",
		"././ui/templates/parts/auth/form-register.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	input := &auth.LoginInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/register.tmpl",
		"././ui/templates/parts/auth/form-register.tmpl",
		"././ui/templates/parts/validation/errors.tmpl",
		"././ui/templates/parts/form-submit/error.tmpl",
	)
	res, err := app.AuthService.Login(input.Email, input.Password)
	if err != nil {
		tmpl.ExecuteTemplate(w, "form-submit-error", err.Error())
		return
	}
	cookie.Set(res.Access, "__up_access_token", 5*time.Minute, w)
	cookie.Set(res.Refresh, "__up_refresh_token", 7*24*time.Hour, w)

	w.Header().Set("HX-Redirect", "/profile/"+strconv.Itoa(res.UserId))
	w.WriteHeader(http.StatusOK)
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	input := &auth.RegisterInput{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Role:     r.FormValue("role"),
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/register.tmpl",
		"././ui/templates/parts/auth/form-register.tmpl",
		"././ui/templates/parts/validation/errors.tmpl",
		"././ui/templates/parts/form-submit/error.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := req.IsValid(input); err != nil {
		var ve req.ValidationErrors
		if errors.As(err, &ve) {
			if err := tmpl.ExecuteTemplate(w, "validation-error", ve); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	if err := app.AuthService.Register(input.Email, input.Password, input.Name, input.Role); err != nil {
		tmpl.ExecuteTemplate(w, "form-submit-error", err.Error())
		return
	}

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	token, err := cookie.Get("__up_refresh_token", r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	err = app.AuthService.Logout(userId, token)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	cookie.Set("", "__up_access_token", -time.Minute, w)
	cookie.Set("", "__up_refresh_token", -time.Hour, w)

	w.Header().Set("HX-Redirect", "/auth/login")
	w.WriteHeader(http.StatusOK)
}

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	uIdFromReq, err := app.getIdFromReq(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	profile, err := app.UsersService.View(uIdFromReq)
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
	currUserId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := app.UsersService.Role(currUserId)

	var userDt []userData
	userDt = append(userDt, userData{
		Profiles:        *profile,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		CanDelete:       canDelete(currUserRole, currUserId, profile.Id),
		CanBan:          canBan(currUserRole, currUserId, profile.Id),
	})
	data := templateData{
		RequestedUserId: uIdFromReq,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		User:            userDt,
	}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (app *application) UserList(w http.ResponseWriter, r *http.Request) {
	page, err := app.getPageFromReq(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if page <= 0 {
		http.Error(w, "Page does not exist", http.StatusInternalServerError)
		return
	}

	uId, err := utils.GetUserID(r.Context())
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, err := app.UsersService.Role(uId)
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
	profiles, res, err := app.UsersService.ViewAll(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var userDt []userData
	for _, profile := range profiles {
		userDt = append(userDt, userData{
			Profiles:  profile,
			CanDelete: canDelete(currRole, uId, profile.Id),
			CanBan:    canBan(currRole, uId, profile.Id),
		})
	}

	totalPages := (res + 10 - 1) / 10

	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	pagination := buildPagination(page, totalPages)
	data := templateData{
		User:          userDt,
		Pagin:         pagination,
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

func (app *application) DeleteConfirm(w http.ResponseWriter, r *http.Request) {
	uId, err := app.getIdFromReq(r)
	if err != nil {
		return
	}
	tmpl, err := view.LoadTemplate(
		"././ui/templates/parts/layout/delete-modal.tmpl",
	)
	if err != nil {
		return
	}
	profile, err := app.UsersService.View(uId)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	var userDt []userData
	userDt = append(userDt, userData{
		Profiles: *profile,
	})
	data := templateData{
		User: userDt,
	}
	tmpl.ExecuteTemplate(w, "delete-modal", data)
}

func (app *application) Delete(w http.ResponseWriter, r *http.Request) {
	currId, err := utils.GetUserID(r.Context())
	if err != nil {
		w.Header().Set("HX-Redirect", "/auth/login")
		w.WriteHeader(http.StatusOK)
		return
	}
	role, err := app.UsersService.Role(currId)
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

func (app *application) Ban(w http.ResponseWriter, r *http.Request) {
	reqId, currId, err := app.selfDeletionDetected(r)
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
	err = app.UsersService.Ban(reqId)
	if err != nil {
		return
	}
	user, err := app.UsersService.View(reqId)
	if err != nil {
		return
	}
	currRole, err := app.UsersService.Role(currId)
	if err != nil {
		return
	}
	data := userData{
		Profiles:      *user,
		CurrentUserId: currId,
		CanDelete:     canDelete(currRole, currId, reqId),
		CanBan:        canBan(currRole, currId, reqId),
	}
	tmpl.ExecuteTemplate(w, "user-"+v, data)
}

func (app *application) Unban(w http.ResponseWriter, r *http.Request) {
	reqId, currId, err := app.selfDeletionDetected(r)
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
	err = app.UsersService.Unban(reqId)
	if err != nil {
		return
	}
	user, err := app.UsersService.View(reqId)
	if err != nil {
		return
	}
	currRole, err := app.UsersService.Role(currId)
	if err != nil {
		return
	}
	data := userData{
		Profiles:      *user,
		CurrentUserId: currId,
		CanDelete:     canDelete(currRole, currId, reqId),
		CanBan:        canBan(currRole, currId, reqId),
	}
	tmpl.ExecuteTemplate(w, "user-"+v, data)
}
