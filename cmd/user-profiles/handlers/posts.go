package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/posts"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	PostService posts.PostService
	TCache      panel.Templates
}

type PostHandlerDeps struct {
	PostService posts.PostService
	TCache      panel.Templates
}

func NewPostHandler(router chi.Router, deps PostHandlerDeps) *PostHandler {
	return &PostHandler{
		PostService: deps.PostService,
		TCache:      deps.TCache,
	}
}

func (handler *PostHandler) Home(w http.ResponseWriter, r *http.Request) {
	page := panel.GetPageFromReq(r)
	limit := 9
	approved := true
	posts, total, err := handler.PostService.GetOnPage(page, limit, &approved, "")
	if err != nil {
		handler.TCache.ServerError(w, err)
		return
	}

	hasMore := page*limit < total

	postCards := make([]panel.PostData, 0, len(posts))
	for _, post := range posts {
		cutedContent := panel.TruncateContent(post.Content, 367)
		post.Content = cutedContent
		postCards = append(postCards, panel.PostData{
			Posts: post,
		})
	}
	currUserId, _ := panel.GetUserId(r)

	data := panel.PageData{
		CurrentUserId: currUserId,
		PostCards:     postCards,
		Loadmore: panel.Loadmore{
			Page:    page,
			Next:    page + 1,
			HasMore: hasMore,
			Total:   total,
		},
	}
	isHTMX := r.Header.Get("HX-Request") == "true"

	if page > 1 || isHTMX {
		err = handler.TCache.RenderPartial(w, "home.tmpl", "posts-response", data)
		if err != nil {
			handler.TCache.ServerError(w, err)
		}
		return
	}
	handler.TCache.Render(w, r, http.StatusOK, "base", "home.tmpl", data)
}

func (handler *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	parts := strings.Split(slug, "-")
	idStr := parts[len(parts)-1]
	pId, err := strconv.Atoi(idStr)
	if err != nil {
		handler.TCache.NotFound(w, r)
		return
	}
	post, err := handler.PostService.FindById(pId)
	if err != nil {
		handler.TCache.NotFound(w, r)
		return
	}
	var postCards []panel.PostData
	postCards = append(postCards, panel.PostData{
		Posts: *post,
	})
	currUserId, _ := panel.GetUserId(r)
	data := panel.PageData{
		CurrentUserId: currUserId,
		PostCards:     postCards,
	}
	handler.TCache.Render(w, r, http.StatusOK, "base", "post.tmpl", data)
}

func (handler *PostHandler) ViewUserPosts(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		handler.TCache.NotFound(w, r)
		return
	}
	currUId, _ := panel.GetUserId(r)
	ps, err := handler.PostService.FindByUsername(username)
	if errors.Is(err, posts.PostNotFound) {
		handler.TCache.NotFound(w, r)
		return
	}
	if err != nil {
		handler.TCache.NotFound(w, r)
		return
	}
	postCards := make([]panel.PostData, 0, len(ps))
	for _, post := range ps {
		cutedContent := panel.TruncateContent(post.Content, 367)
		post.Content = cutedContent
		postCards = append(postCards, panel.PostData{
			Posts: post,
		})
	}
	data := panel.PageData{
		PostCards: postCards,
		CurrentUserId: currUId,
	}
	handler.TCache.Render(w, r, http.StatusOK, "base", "user-posts.tmpl", data)
}

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
