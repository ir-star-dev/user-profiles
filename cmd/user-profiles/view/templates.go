package view

import (
	"bytes"
	"html/template"
	"net/http"
)

type Templates struct {
	TemplateCache map[string]*template.Template
}

var commonTemplates = []string{
	"././ui/base.tmpl",
	"././ui/parts/layout/nav.tmpl",
}

func NewTemplates() *Templates {
	return &Templates{
		TemplateCache: make(map[string]*template.Template),
	}
}

func (t *Templates) Render(w http.ResponseWriter, name string, data any, files ...string) error {
	tmpl, ok := t.TemplateCache[name]
	if !ok {
		allFiles := append(commonTemplates, files...)

		parsed, err := template.ParseFiles(allFiles...)
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

func (t *Templates) RenderPartial(w http.ResponseWriter, templateName string, data any,	files ...string) error {
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

func (t *Templates) BuildPagination(currentPage, totalPages int) Pagination {
	p := Pagination{
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

func Truncate(s string, limit int) string {
    r := []rune(s)
    if len(r) > limit {
        return string(r[:limit]) + "..."
    }
    return s
}