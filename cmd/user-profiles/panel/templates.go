package panel

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/ui"

	"github.com/microcosm-cc/bluemonday"
)

type Templates struct {
	templateCache  map[string]*template.Template
}

var functions = template.FuncMap{
	"withQuery": WithQuery,
}

func NewTemplateCache() (*Templates, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	panelPages, err := fs.Glob(ui.Files, "pages/panel/*.tmpl")
	if err != nil {
		return nil, err
	}
	pages = append(pages, panelPages...)

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"base.tmpl",
			"panel-base.tmpl",
			"pages/*.tmpl",
			"pages/panel/*.tmpl",
			"parts/*/*.tmpl",
			"parts/*/*/*.tmpl",
			"parts/*/*/*/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return &Templates{
		templateCache: cache,
	}, nil
}

func (t *Templates) Render(w http.ResponseWriter, r *http.Request, status int, layout string, page string, data any) {
	ts, ok := t.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", page)
		t.ServerError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, layout, data)
	if err != nil {
		t.ServerError(w, r, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (t *Templates) RenderPanel(w http.ResponseWriter, r *http.Request, status int, page string, data any) {
    t.Render(w, r, status, "panel-base", page, data)
}

func (t *Templates) RenderPartial(w http.ResponseWriter, page string, templateName string, data any) error {
	ts, ok := t.templateCache[page]
	if !ok {
		return fmt.Errorf("Template %s does not exist", page)
	}
	var buf bytes.Buffer
	err := ts.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func (t *Templates) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	data := PageData{}
	uId, _ := GetUserId(r)
	data.CurrentUserId = uId
	t.Render(w,	r, http.StatusNotFound,	"base", "500.tmpl", data)
}

func (t *Templates) NotFound(w http.ResponseWriter, r *http.Request, data any) {
	t.Render(w,	r, http.StatusNotFound,	"base", "404.tmpl",	data)
}

func (t *Templates) PanelNotFound(w http.ResponseWriter, r *http.Request) {
	data := PageData{}
	t.Render(w,	r, http.StatusNotFound,	"panel-base", "404.tmpl",	data)
}

func (t *Templates) BuildPagination(currentPage, totalPages int, path string) Pagination {
	p := Pagination{
		Path:       path,
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

func (t *Templates) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {

				w.Header().Set("Connection", "close")

				t.ServerError(w, r, fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func WithQuery(base string, params map[string]string, key string, value any) string {
	q := url.Values{}

	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}

	if value == nil {
		return base + "?" + q.Encode()
	}

	q.Set(key, fmt.Sprint(value))

	return base + "?" + q.Encode()
}

func TruncateContent(s string, limit int) string {
	r := []rune(s)
	if len(r) > limit {
		return string(r[:limit]) + "..."
	}
	return s
}

func SanitizeContent(html string) string {
	var policy = bluemonday.UGCPolicy()
	policy.AllowAttrs("href").OnElements("a")
	policy.RequireNoFollowOnLinks(true)
	
	return policy.Sanitize(html)
}

func GetIdFromReq(r *http.Request) (int, error) {
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

func GetUserId(r *http.Request) (int, error) {
	uId, err := utils.GetUserID(r.Context())
	if err != nil {
		return 0, err
	}
	return uId, nil
}

func GetPageFromReq(r *http.Request) int {
	page := r.URL.Query().Get("page")
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		return 1
	}
	return p
}

func GetFilterValue(r *http.Request, key string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return ""
	}
	return value
}

func SelfDeletionDetected(r *http.Request) (int, int, error) {
	reqId, err := GetIdFromReq(r)
	if err != nil {
		return 0, 0, err
	}
	currId, err := GetUserId(r)
	if err != nil {
		return reqId, currId, err
	}
	if reqId == currId {
		return reqId, currId, errors.New("You can't delete yourself.")
	}
	return reqId, currId, nil
}

func CanDeleteUser(currentRole string, currentUserId, profileId int) bool {
	if currentRole == "admin" {
		return currentUserId != profileId
	}
	return currentUserId == profileId
}

func CanBanUser(currentRole string, currentUserId, profileId int) bool {
	return currentRole == "admin" && currentUserId != profileId
}

func CanDeletePost(currentRole string) bool {
	if currentRole == "admin" || currentRole == "moderator" {
		return true
	}
	return false
}

func CanApprovePost(currentRole string) bool {
	if currentRole == "admin" || currentRole == "moderator" {
		return true
	}
	return false
}

func CanEditPost(currentRole string, currentUserId, authorId int) bool {
	if currentRole == "admin" || currentRole == "moderator" {
		return true
	}
	return currentUserId == authorId
}
