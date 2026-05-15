package handlers

import (
	"errors"
	"html/template"
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
		handler.TCache.ServerError(w, r, err)
		return
	}

	hasMore := page*limit < total

	postCards := make([]panel.PostData, 0, len(posts))
	for _, post := range posts {
		cutedContent := panel.TruncateContent(string(post.Content), 367)
		post.Content = template.HTML(cutedContent)
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
			handler.TCache.ServerError(w, r, err)
		}
		return
	}
	handler.TCache.Render(w, r, http.StatusOK, "base", "home.tmpl", data)
}

func (handler *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	currUserId, _ := panel.GetUserId(r)
	data := panel.PageData{
		CurrentUserId: currUserId,
	}
	slug := r.PathValue("slug")
	parts := strings.Split(slug, "-")
	idStr := parts[len(parts)-1]
	pId, err := strconv.Atoi(idStr)
	if err != nil {
		handler.TCache.NotFound(w, r, data)
		return
	}
	post, err := handler.PostService.FindById(pId)
	if err != nil {
		handler.TCache.NotFound(w, r, data)
		return
	}
	var postCards []panel.PostData
	postCards = append(postCards, panel.PostData{
		Posts: *post,
	})
	data.PostCards = postCards
	handler.TCache.Render(w, r, http.StatusOK, "base", "post.tmpl", data)
}

func (handler *PostHandler) ViewUserPosts(w http.ResponseWriter, r *http.Request) {
	currUId, _ := panel.GetUserId(r)
	data := panel.PageData{
		CurrentUserId: currUId,
	}
	username := r.PathValue("username")
	if username == "" {
		handler.TCache.NotFound(w, r, data)
		return
	}
	ps, err := handler.PostService.FindByUsername(username)
	if errors.Is(err, posts.PostNotFound) {
		handler.TCache.NotFound(w, r, data)
		return
	}
	if err != nil {
		handler.TCache.NotFound(w, r, data)
		return
	}
	postCards := make([]panel.PostData, 0, len(ps))
	for _, post := range ps {
		cutedContent := panel.TruncateContent(string(post.Content), 367)
		post.Content = template.HTML(cutedContent)
		postCards = append(postCards, panel.PostData{
			Posts: post,
		})
	}
	data.PostCards = postCards
	handler.TCache.Render(w, r, http.StatusOK, "base", "user-posts.tmpl", data)
}
