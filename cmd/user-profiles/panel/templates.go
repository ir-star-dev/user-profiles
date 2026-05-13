package panel

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"user-profiles/cmd/user-profiles/utils"
)

type Templates struct {
	TemplateCache map[string]*template.Template
	funcMap       template.FuncMap
}

var commonTemplates = []string{
	"././ui/base.tmpl",
	"././ui/parts/layout/nav.tmpl",
}

func NewTemplates() *Templates {
	return &Templates{
		TemplateCache: make(map[string]*template.Template),
		funcMap: template.FuncMap{
			"withQuery": WithQuery,
		},
	}
}

func (t *Templates) Render(w http.ResponseWriter, name string, data any, files ...string) error {
	tmpl, ok := t.TemplateCache[name]
	if !ok {
		allFiles := append(commonTemplates, files...)

		parsed, err := template.New("base").Funcs(t.funcMap).ParseFiles(allFiles...)
		if err != nil {
			return err
		}

		tmpl = parsed
		t.TemplateCache[name] = tmpl
	}
	var buf bytes.Buffer

	err := tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func (t *Templates) RenderPartial(w http.ResponseWriter, templateName string, data any, files ...string) error {
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func (t *Templates) ServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
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
