package handlers

import (
	"bytes"
	"net/http"

	//"html/template"
	//"time"
	"user-profiles/cmd/view"
	"user-profiles/configs"

	//"user-profiles/internal/http/cookie"

	"user-profiles/internal/auth"
	"user-profiles/internal/posts"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	Config      *configs.Config
	PostService posts.PostService
	AuthService auth.AuthService
	JWTService  auth.JWTService
}

type PostHandlerDeps struct {
	Config      *configs.Config
	PostService posts.PostService
	AuthService auth.AuthService
	JWTService  auth.JWTService
}

func NewPostHandler(router chi.Router, deps PostHandlerDeps) *PostHandler {
	return &PostHandler{
		Config:      deps.Config,
		PostService: deps.PostService,
		AuthService: deps.AuthService,
		JWTService:  deps.JWTService,
	}
}

func (handler *PostHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/post.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) ViewUserPosts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/post.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) PostList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) Publish(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) Review(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}

func (handler *PostHandler) Edit(w http.ResponseWriter, r *http.Request) {
	tmpl, err := view.LoadTemplate(
		"././ui/templates/base.tmpl",
		"././ui/templates/pages/index.tmpl",
		"././ui/templates/parts/layout/posts.tmpl",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := PageData{}
	var buf bytes.Buffer

	err = tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buf.Bytes())
}