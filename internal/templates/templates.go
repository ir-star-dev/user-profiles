package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"user-profiles/internal/models"
	"user-profiles/internal/utils"
	"user-profiles/ui"

	"github.com/lmittmann/tint"
)

type Templates struct {
	templateCache map[string]*template.Template
}

var functions = template.FuncMap{
	"withQuery": utils.WithQuery,
	"add":       utils.Increment,
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

func (t *Templates) RenderSpecialTemplate(w http.ResponseWriter, page, name string) (error, *bytes.Buffer) {
	ts, ok := t.templateCache[page]
	if !ok {
		return fmt.Errorf("Template %s does not exist", page), nil
	}
	var buf bytes.Buffer
	err := ts.ExecuteTemplate(&buf, name, nil)
	if err != nil {
		return err, nil
	}
	return nil, &buf
}

func (t *Templates) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level: slog.LevelDebug,
		}),
	)
	_, file, line, _ := runtime.Caller(0)

	logger.Error("500 Error",
		slog.String("file", file),
		slog.Int("line", line),
		slog.Any("error", err),
	)
	err, buf := t.RenderSpecialTemplate(w, "500.tmpl", "500")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	buf.WriteTo(w)
}

func (t *Templates) NotFound(w http.ResponseWriter, r *http.Request) {
	err, buf := t.RenderSpecialTemplate(w, "404.tmpl", "404")
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	buf.WriteTo(w)
}

func (t *Templates) PanelNotFound(w http.ResponseWriter, r *http.Request) {
	err, buf := t.RenderSpecialTemplate(w, "404.tmpl", "404")
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusForbidden)
	buf.WriteTo(w)
}

func (t *Templates) Forbidden(w http.ResponseWriter, r *http.Request) {
	err, buf := t.RenderSpecialTemplate(w, "403.tmpl", "403")
	if err != nil {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusForbidden)
	buf.WriteTo(w)
}

func (t *Templates) BuildPagination(currentPage, totalPages int, path string) models.Pagination {
	p := models.Pagination{
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
