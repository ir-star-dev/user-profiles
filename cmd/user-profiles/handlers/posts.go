package handlers

import (
	//"bytes"
	"net/http"
	"strconv"
	"strings"

	//"html/template"
	//"time"
	"user-profiles/cmd/user-profiles/view"
	"user-profiles/configs"

	//"user-profiles/internal/http/cookie"

	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/posts"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	Config      *configs.Config
	PostService posts.PostService
	AuthService auth.AuthService
	JWTService  auth.JWTService
	Templates   view.Templates
}

type PostHandlerDeps struct {
	Config      *configs.Config
	PostService posts.PostService
	AuthService auth.AuthService
	JWTService  auth.JWTService
	Templates   view.Templates
}

func NewPostHandler(router chi.Router, deps PostHandlerDeps) *PostHandler {
	return &PostHandler{
		Config:      deps.Config,
		PostService: deps.PostService,
		AuthService: deps.AuthService,
		JWTService:  deps.JWTService,
		Templates:   deps.Templates,
	}
}

func (handler *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
	pageFromReq := r.URL.Query().Get("page")
	page := 1
	if pageFromReq != "" {
		if p, err := strconv.Atoi(pageFromReq); err == nil {
			page = p
		}
	}
	if page <= 0 {
		page = 1
	}
	limit := 9
	posts, total, err := handler.PostService.GetAll(page, limit)
	if err != nil {
		handler.Templates.ServerError(w, err)
		return
	}

	hasMore := page*limit < total

	postCards := make([]view.PostData, 0, len(posts))
	for _, post := range posts {
		cutedContent := view.Truncate(post.Content, 367)
		post.Content = cutedContent
		postCards = append(postCards, view.PostData{
			Posts: post,
		})
	}

	data := view.PageData{
		PostCards: postCards,
		Loadmore: view.Loadmore{
			Page:    page,
			Next:    page + 1,
			HasMore: hasMore,
			Total:   total,
		},
	}
	isHTMX := r.Header.Get("HX-Request") == "true"

	if page > 1 || isHTMX {
		err = handler.Templates.RenderPartial(
			w,
			"posts-response",
			data,
			"././ui/parts/layout/posts-response.tmpl",
			"././ui/parts/layout/posts-list.tmpl",
			"././ui/parts/layout/post-card.tmpl",
			"././ui/parts/layout/loadmore.tmpl",
		)
		if err != nil {
			handler.Templates.ServerError(w, err)
		}
		return
	}
	err = handler.Templates.Render(w, "home", data,
		"././ui/pages/home.tmpl",
		"././ui/parts/layout/posts-list.tmpl",
		"././ui/parts/layout/post-card.tmpl",
		"././ui/parts/layout/loadmore.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	parts := strings.Split(slug, "-")
	idStr := parts[len(parts)-1]
	pId, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	post, err := handler.PostService.FindById(pId)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var postCards []view.PostData
	postCards = append(postCards, view.PostData{
		Posts: *post,
	})
	data := view.PageData{
		PostCards: postCards,
	}
	err = handler.Templates.Render(w, "post", data, "././ui/pages/post.tmpl")
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
}

func (handler *PostHandler) ViewUserPosts(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		http.NotFound(w, r)
		return
	}

	posts, err := handler.PostService.FindByUsername(username)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}
	postCards := make([]view.PostData, 0, len(posts))
	for _, post := range posts {
		cutedContent := view.Truncate(post.Content, 367)
		post.Content = cutedContent
		postCards = append(postCards, view.PostData{
			Posts: post,
		})
	}
	data := view.PageData{
		PostCards: postCards,
	}
	err = handler.Templates.Render(w, "user-posts", data,
		"././ui/pages/user-posts.tmpl",
		"././ui/parts/layout/posts-list.tmpl",
		"././ui/parts/layout/post-card.tmpl",
	)
	if err != nil {
		handler.Templates.ServerError(w, err)
	}

}

// func (handler *PostHandler) PostList(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }

// func (handler *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }

// func (handler *PostHandler) Publish(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }

// func (handler *PostHandler) Review(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }

// func (handler *PostHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }

// func (handler *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }

// func (handler *PostHandler) EditForm(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }

// func (handler *PostHandler) Edit(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := view.LoadTemplate(
// 		"././ui/templates/base.tmpl",
// 		"././ui/templates/parts/layout/nav.tmpl",
// 		"././ui/templates/pages/index.tmpl",
// 		"././ui/templates/parts/layout/posts.tmpl",
// 	)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	data := view.PageData{}
// 	var buf bytes.Buffer

// 	err = tmpl.ExecuteTemplate(&buf, "base", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(buf.Bytes())
// }
