package posts

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"user-profiles/internal/http/request"
	"user-profiles/internal/models"
	"user-profiles/internal/sanitizer"
	"user-profiles/internal/templates"
	"user-profiles/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

type PHandler struct {
	PService PostService
	TCache   templates.Templates
}

type PHandlerDeps struct {
	PService PostService
	TCache   templates.Templates
}

func NewPostHandler(router chi.Router, deps PHandlerDeps) *PHandler {
	return &PHandler{
		PService: deps.PService,
		TCache:   deps.TCache,
	}
}

func (h *PHandler) Home(w http.ResponseWriter, r *http.Request) {
	currUserRole, _ := request.GetUserRole(r.Context())
	page := request.GetPageFromReq(r)
	limit := 9
	approved := true
	posts, total, err := h.PService.GetOnPage(page, limit, &approved, "", "")
	if err != nil {
		h.TCache.ServerError(w, r, err)
		return
	}
	hasMore := page*limit < total
	postCards := make([]models.PostData, 0, len(posts))
	for _, post := range posts {
		postCards = append(postCards, models.PostData{
			Posts: models.PostViewTable{
				Post: post,
			},
		})
	}
	currUserId, _ := request.GetUserId(r)
	data := models.PageData{
		CurrentUserId: currUserId,
		CurrentUserRole: currUserRole,
		PostCards:     postCards,
		Loadmore: models.Loadmore{
			Page:    page,
			Next:    page + 1,
			HasMore: hasMore,
			Total:   total,
		},
		CSRFToken: csrf.Token(r),
	}
	isHTMX := r.Header.Get("HX-Request") == "true"
	if page > 1 || isHTMX {
		err = h.TCache.RenderPartial(w, "home.tmpl", "posts-response", data)
		if err != nil {
			h.TCache.ServerError(w, r, err)
		}
		return
	}
	h.TCache.Render(w, r, http.StatusOK, "base", "home.tmpl", data)
}

func (h *PHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	currUserId, _ := request.GetUserId(r)
	data := models.PageData{
		CurrentUserId: currUserId,
		CSRFToken:     csrf.Token(r),
	}
	slug := r.PathValue("slug")
	pId, err := utils.GetPostIdFromSlug(slug)
	if err != nil {
		h.TCache.NotFound(w, r)
		return
	}
	post, err := h.PService.FindById(*pId)
	if err != nil {
		h.TCache.NotFound(w, r)
		return
	}
	var postCards []models.PostData
	postCards = append(postCards, models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
	})
	data.PostCards = postCards
	h.TCache.Render(w, r, http.StatusOK, "base", "post.tmpl", data)
}

func (h *PHandler) ViewUserPosts(w http.ResponseWriter, r *http.Request) {
	currUId, _ := request.GetUserId(r)
	data := models.PageData{
		CurrentUserId: currUId,
		CSRFToken:     csrf.Token(r),
	}
	username := r.PathValue("username")
	if username == "" {
		h.TCache.NotFound(w, r)
		return
	}
	ps, err := h.PService.FindByUsername(username)
	if errors.Is(err, PostNotFound) {
		h.TCache.NotFound(w, r)
		return
	}
	if err != nil {
		h.TCache.NotFound(w, r)
		return
	}
	postCards := make([]models.PostData, 0, len(ps))
	for _, post := range ps {
		postCards = append(postCards, models.PostData{
			Posts: models.PostViewTable{
				Post: post,
			},
		})
	}
	data.PostCards = postCards
	h.TCache.Render(w, r, http.StatusOK, "base", "user-posts.tmpl", data)
}

func (h *PHandler) Posts(w http.ResponseWriter, r *http.Request) {
	page := request.GetPageFromReq(r)
	approved := request.GetFilterValue(r, "approved")
	username := request.GetFilterValue(r, "username")
	search := strings.TrimSpace(request.GetFilterValue(r, "search"))

	currUserId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currRole, _ := request.GetUserRole(r.Context())
	data := models.PageData{
		CurrentUserId:   currUserId,
		CurrentUserRole: currRole,
		CSRFToken:       csrf.Token(r),
	}
	var ap *bool
	if approved != "" {
		v, err := strconv.ParseBool(approved)
		if err != nil {
			ap = nil
		}
		ap = &v
	}
	posts, totalPosts, err := h.PService.GetOnPage(page, 10, ap, username, search)
	if err != nil {
		h.TCache.ServerError(w, r, err)
		return
	}
	authors, _ := h.PService.Authors()

	totalPages := (totalPosts + 10 - 1) / 10
	pages := []int{}
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	var cards []models.PostData
	for i, post := range posts {
		cards = append(cards, models.PostData{
			Posts: models.PostViewTable{
				Post:  post,
				Index: ((page - 1) * 10) + i + 1,
			},
			Actions: models.Actions{
				CanDeletePost:  CanDeletePost(currRole),
				CanApprovePost: CanApprovePost(currRole),
				CanEditPost:    CanEditPost(currRole, currUserId, post.UserId),
			},
		})
	}

	pagination := h.TCache.BuildPagination(page, totalPages, "/panel/posts")
	data.Pagination = pagination
	data.PostCards = cards
	data.Authors = authors
	data.Filters = map[string]string{
		"approved": approved,
		"username": username,
		"search":   search,
	}
	data.HasFilters = page > 1 || approved != "" || username != "" || search != ""

	data.TotalPosts = totalPosts
	data.Page = page
	data.Pages = len(pages)

	if r.Header.Get("HX-Request") == "true" {
		h.TCache.RenderPartial(w, "posts.tmpl", "post-table", data)
		return
	}
	h.TCache.RenderPanel(w, r, http.StatusOK, "posts.tmpl", data)
}

func (h *PHandler) Publish(w http.ResponseWriter, r *http.Request) {
	reqUId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	currUId, err := request.GetUserId(r)
	if err != nil {
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	v := request.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	err = h.PService.Publish(pId)
	if err != nil {
		return
	}
	post, err := h.PService.FindById(pId)
	if err != nil {

		return
	}
	data := models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
		Actions: models.Actions{
			CanDeletePost:  CanDeletePost(currUserRole),
			CanEditPost:    CanEditPost(currUserRole, currUId, reqUId),
			CanApprovePost: CanApprovePost(currUserRole),
		},
	}
	page := "post-" + v
	err = h.TCache.RenderPartial(w, "posts.tmpl", page, data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) Review(w http.ResponseWriter, r *http.Request) {
	reqUId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	currUId, err := request.GetUserId(r)
	if err != nil {
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	v := request.GetFilterValue(r, "view")
	if v == "" {
		return
	}
	err = h.PService.Review(pId)
	if err != nil {
		return
	}
	post, err := h.PService.FindById(pId)
	if err != nil {
		return
	}
	data := models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
		Actions: models.Actions{
			CanDeletePost:  CanDeletePost(currUserRole),
			CanEditPost:    CanEditPost(currUserRole, currUId, reqUId),
			CanApprovePost: CanApprovePost(currUserRole),
		},
	}
	page := "post-" + v
	err = h.TCache.RenderPartial(w, "posts.tmpl", page, data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) DeletePostConfirm(w http.ResponseWriter, r *http.Request) {
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	post, err := h.PService.FindById(pId)
	if err != nil {
		return
	}
	var cards []models.PostData
	cards = append(cards, models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
	})
	data := models.PageData{
		PostCards: cards,
		CSRFToken: csrf.Token(r),
	}
	err = h.TCache.RenderPartial(w, "posts.tmpl", "confirm-delete-post-modal", data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	err = h.PService.Delete(pId)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
	w.Header().Set("HX-Trigger", "postDeleted")
}

func (h *PHandler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	data := models.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}
	h.TCache.RenderPanel(w, r, http.StatusOK, "create-post.tmpl", data)
}

func (h *PHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	data := models.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}
	title := strings.TrimSpace(r.FormValue("title"))
	excerpt := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("excerpt")))
	content := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("content")))

	validationErrs, err := h.PService.Create(title, content, excerpt, uId)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := h.TCache.RenderPartial(w, "create-post.tmpl", "form-submit-error", data)
		if err != nil {
			h.TCache.ServerError(w, r, err)
		}
		return
	}
	data.PostCreated = models.PostCreated{
		Message: "Post successfully created and waiting for moderation!",
	}
	err = h.TCache.RenderPartial(w, "create-post.tmpl", "form-submit-success", data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) PreviewPost(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	data := models.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}
	post, err := h.PService.PreviewPost(pId)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}
	var postCards []models.PostData
	postCards = append(postCards, models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
		Actions: models.Actions{
			CanApprovePost: CanApprovePost(currUserRole),
		},
	})
	data.PostCards = postCards
	h.TCache.RenderPanel(w, r, http.StatusOK, "preview-post.tmpl", data)
}

func (h *PHandler) EditPostForm(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	data := models.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		h.TCache.NotFound(w, r)
		return
	}
	post, err := h.PService.PreviewPost(pId)
	if err != nil {
		h.TCache.NotFound(w, r)
		return
	}
	var postCard []models.PostData
	postCard = append(postCard, models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
	})
	data.PostCards = postCard
	h.TCache.RenderPanel(w, r, http.StatusOK, "edit-post.tmpl", data)
}

func (h *PHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	uId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	data := models.PageData{
		CurrentUserId:   uId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("content")))
	excerpt := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("excerpt")))
	createdAt := r.FormValue("created_at")

	validationErrs, err := h.PService.Update(title, content, excerpt, createdAt, pId)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := h.TCache.RenderPartial(w, "edit-post.tmpl", "form-submit-error", data)
		if err != nil {
			h.TCache.ServerError(w, r, err)
		}
		return
	}
	data.PostUpdated = models.PostUpdated{
		Message: "Post updated and waiting for moderation!",
	}
	err = h.TCache.RenderPartial(w, "edit-post.tmpl", "form-submit-success", data)
	if err != nil {
		h.TCache.ServerError(w, r, err)
	}
}
