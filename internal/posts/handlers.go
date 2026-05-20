package posts

import (
	"net/http"
	"strconv"
	"strings"

	"user-profiles/internal/http/request"
	"user-profiles/internal/models"
	"user-profiles/internal/sanitizer"
	"user-profiles/internal/templates"
	"user-profiles/internal/utils"

	"github.com/go-chi/chi/v5"
)

type PHandler struct {
	PService PostService
	TCache   templates.Templates
	*templates.BaseHandler
}

type PHandlerDeps struct {
	PService PostService
	TCache   templates.Templates
	*templates.BaseHandler
}

func NewPostHandler(router chi.Router, deps PHandlerDeps) *PHandler {
	return &PHandler{
		PService:    deps.PService,
		BaseHandler: deps.BaseHandler,
	}
}

func (h *PHandler) Home(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(request.GetFilterValue(r, "search"))
	sort := strings.TrimSpace(request.GetFilterValue(r, "sort"))

	base := h.NewBasePageData(r)
	page := request.GetPageFromReq(r)
	limit := 9
	approved := true
	filters := models.PostFilters{
		Approved: &approved,
		Search:   search,
		Sort:     sort,
	}
	posts, total, err := h.PService.GetOnPage(page, limit, filters)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
		return
	}
	hasMore := page*limit < total
	postCards := make([]models.PostData, 0, total)
	for _, post := range posts {
		postCards = append(postCards, models.PostData{
			Posts: models.PostViewTable{
				Post: post,
			},
		})
	}

	data := models.HomeData{
		BasePageData: base,
		PostCards:    postCards,
		Loadmore: models.Loadmore{
			Page:    page,
			Next:    page + 1,
			HasMore: hasMore,
			Total:   total,
		},
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	isLoadMore := isHTMX && page > 1
	isFilters := isHTMX && (search != "" || sort != "")

	if isLoadMore {
		err = h.BaseHandler.TCache.RenderPartial(w, "home.tmpl", "posts-loadmore", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	}
	if isFilters {
		data.Filters = map[string]string{
			"search": search,
			"sort":   sort,
		}
		data.HasFilters = search != "" || sort != ""
		err = h.BaseHandler.TCache.RenderPartial(w, "home.tmpl", "posts-response", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	}
	h.BaseHandler.TCache.Render(w, r, http.StatusOK, "base", "home.tmpl", data)
}

func (h *PHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)

	slug := r.PathValue("slug")
	pId, err := utils.GetPostIdFromSlug(slug)
	if err != nil {
		h.BaseHandler.TCache.NotFound(w, r)
		return
	}
	post, err := h.PService.FindById(*pId)
	if err != nil {
		h.BaseHandler.TCache.NotFound(w, r)
		return
	}
	var postCards []models.PostData
	postCards = append(postCards, models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
	})
	data := models.PostSingleData{
		PostCards:    postCards,
		BasePageData: base,
	}
	h.BaseHandler.TCache.Render(w, r, http.StatusOK, "base", "post.tmpl", data)
}

func (h *PHandler) ViewUserPosts(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(request.GetFilterValue(r, "search"))
	sort := strings.TrimSpace(request.GetFilterValue(r, "sort"))
	base := h.NewBasePageData(r)

	username := r.PathValue("username")
	if username == "" {
		h.BaseHandler.TCache.NotFound(w, r)
		return
	}
	data := models.UserPostsData{
		BasePageData: base,
		Filters: map[string]string{
			"search":   search,
			"sort":     sort,
			"username": username,
		},
		HasFilters: search != "" || sort != "" || username != "",
	}
	approved := true
	filters := models.PostFilters{
		Approved: &approved,
		Username: username,
		Search:   search,
		Sort:     sort,
	}
	ps, total, err := h.PService.GetOnPage(-1, -1, filters)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
		return
	}
	postCards := make([]models.PostData, 0, total)
	for _, post := range ps {
		postCards = append(postCards, models.PostData{
			Posts: models.PostViewTable{
				Post: post,
			},
		})
	}
	data.PostCards = postCards

	isHTMX := r.Header.Get("HX-Request") == "true"
	isFilters := isHTMX && (search != "" || sort != "")

	if isFilters {
		err = h.BaseHandler.TCache.RenderPartial(w, "user-posts.tmpl", "user-posts-response", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	}
	h.BaseHandler.TCache.Render(w, r, http.StatusOK, "base", "user-posts.tmpl", data)
}

// Dashboard
func (h *PHandler) Posts(w http.ResponseWriter, r *http.Request) {
	page := request.GetPageFromReq(r)
	approved := request.GetFilterValue(r, "approved")
	username := request.GetFilterValue(r, "username")
	search := strings.TrimSpace(request.GetFilterValue(r, "search"))

	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	authors, _ := h.PService.Authors()
	data := models.PostsDashboardData{
		BasePageData: base,
		Filters: map[string]string{
			"approved": approved,
			"username": username,
			"search":   search,
		},
		HasFilters: page > 1 || approved != "" || username != "" || search != "",
		Authors:    authors,
	}
	var ap *bool
	if approved != "" {
		v, err := strconv.ParseBool(approved)
		if err != nil {
			ap = nil
		}
		ap = &v
	}
	filters := models.PostFilters{
		Approved: ap,
		Username: username,
		Search:   search,
	}
	posts, totalPosts, err := h.PService.GetOnPage(page, 10, filters)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
		return
	}

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
				CanDeletePost:  CanDeletePost(base.CurrentUserRole),
				CanApprovePost: CanApprovePost(base.CurrentUserRole),
				CanEditPost:    CanEditPost(base.CurrentUserRole, base.CurrentUserId, post.UserId),
			},
		})
	}

	pagination := h.BaseHandler.TCache.BuildPagination(page, totalPages, "/panel/posts")
	data.Pagination = pagination
	data.PostCards = cards
	data.TotalPosts = totalPosts
	data.Page = page
	data.Pages = len(pages)

	if r.Header.Get("HX-Request") == "true" {
		h.BaseHandler.TCache.RenderPartial(w, "posts.tmpl", "post-table", data)
		return
	}
	h.BaseHandler.TCache.RenderPanel(w, r, http.StatusOK, "posts.tmpl", data)
}

func (h *PHandler) Publish(w http.ResponseWriter, r *http.Request) {
	modData := h.PrepareToModeration(w, r)
	if modData == nil {
		return
	}
	err := h.PService.Publish(modData.PID)
	if err != nil {
		return
	}
	post, err := h.PService.FindById(modData.PID)
	if err != nil {
		return
	}
	data := models.PostData{
		Posts: models.PostViewTable{
			Post:  *post,
			Index: modData.Index,
		},
		Actions: models.Actions{
			CanDeletePost:  CanDeletePost(modData.Data.CurrentUserRole),
			CanEditPost:    CanEditPost(modData.Data.CurrentUserRole, modData.CurrUID, modData.ReqUI),
			CanApprovePost: CanApprovePost(modData.Data.CurrentUserRole),
		},
	}
	page := "post-" + modData.View
	err = h.BaseHandler.TCache.RenderPartial(w, "posts.tmpl", page, data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) Review(w http.ResponseWriter, r *http.Request) {
	modData := h.PrepareToModeration(w, r)
	if modData == nil {
		return
	}
	err := h.PService.Review(modData.PID)
	if err != nil {
		return
	}
	post, err := h.PService.FindById(modData.PID)
	if err != nil {
		return
	}
	data := models.PostData{
		Posts: models.PostViewTable{
			Post:  *post,
			Index: modData.Index,
		},
		Actions: models.Actions{
			CanDeletePost:  CanDeletePost(modData.Data.CurrentUserRole),
			CanEditPost:    CanEditPost(modData.Data.CurrentUserRole, modData.CurrUID, modData.ReqUI),
			CanApprovePost: CanApprovePost(modData.Data.CurrentUserRole),
		},
	}
	page := "post-" + modData.View
	err = h.BaseHandler.TCache.RenderPartial(w, "posts.tmpl", page, data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) DeletePostConfirm(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
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
	data := models.PostDeleteData{
		PostCards:    cards,
		BasePageData: base,
	}
	err = h.BaseHandler.TCache.RenderPartial(w, "posts.tmpl", "confirm-delete-post-modal", data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		return
	}
	err = h.PService.Delete(pId)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
	w.Header().Set("HX-Trigger", "postDeleted")
}

func (h *PHandler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	data := models.CreatePostFormData{
		BasePageData: base,
	}
	h.BaseHandler.TCache.RenderPanel(w, r, http.StatusOK, "create-post.tmpl", data)
}

func (h *PHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	data := models.CreatePostFormData{
		BasePageData: base,
	}
	title := strings.TrimSpace(r.FormValue("title"))
	excerpt := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("excerpt")))
	content := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("content")))

	validationErrs, err := h.PService.Create(title, content, excerpt, base.CurrentUserId)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := h.BaseHandler.TCache.RenderPartial(w, "create-post.tmpl", "form-submit-error", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	}
	data.Success.PostCreated = &models.PostCreated{
		Message: "Post successfully created and waiting for moderation!",
	}

	err = h.BaseHandler.TCache.RenderPartial(w, "create-post.tmpl", "form-submit-success", data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) PreviewPost(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := request.GetUserRole(r.Context())
	data := models.PreviewPostData{
		BasePageData: base,
	}
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		h.BaseHandler.TCache.PanelNotFound(w, r)
		return
	}
	post, err := h.PService.PreviewPost(pId)
	if err != nil {
		h.BaseHandler.TCache.PanelNotFound(w, r)
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
	h.BaseHandler.TCache.RenderPanel(w, r, http.StatusOK, "preview-post.tmpl", data)
}

func (h *PHandler) EditPostForm(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	data := models.EditPostFormData{
		BasePageData: base,
	}
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		h.BaseHandler.TCache.NotFound(w, r)
		return
	}
	post, err := h.PService.PreviewPost(pId)
	if err != nil {
		h.BaseHandler.TCache.NotFound(w, r)
		return
	}
	var postCard []models.PostData
	postCard = append(postCard, models.PostData{
		Posts: models.PostViewTable{
			Post: *post,
		},
	})
	data.PostCards = postCard
	h.BaseHandler.TCache.RenderPanel(w, r, http.StatusOK, "edit-post.tmpl", data)
}

func (h *PHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	base := h.NewBasePageData(r)
	if base.CurrentUserId == 0 {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	data := models.EditPostFormData{
		BasePageData: base,
	}
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		h.BaseHandler.TCache.PanelNotFound(w, r)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("content")))
	excerpt := sanitizer.SanitizeContent(strings.TrimSpace(r.FormValue("excerpt")))
	createdAt := r.FormValue("created_at")

	validationErrs, err := h.PService.Update(title, content, excerpt, createdAt, pId)
	if err != nil {
		data.FormValidationErr = validationErrs
		err := h.BaseHandler.TCache.RenderPartial(w, "edit-post.tmpl", "form-submit-error", data)
		if err != nil {
			h.BaseHandler.TCache.ServerError(w, r, err)
		}
		return
	}
	data.Success.PostUpdated = &models.PostUpdated{
		Message: "Post updated and waiting for moderation!",
	}

	err = h.BaseHandler.TCache.RenderPartial(w, "edit-post.tmpl", "form-submit-success", data)
	if err != nil {
		h.BaseHandler.TCache.ServerError(w, r, err)
	}
}

func (h *PHandler) PrepareToModeration(w http.ResponseWriter, r *http.Request) *ModerationData {
	index, _ := strconv.Atoi(request.GetFilterValue(r, "index"))
	reqUId, err := request.GetIdFromReq(r)
	if err != nil {
		return nil
	}
	base := h.NewBasePageData(r)
	currUId, err := request.GetUserId(r)
	if base.CurrentUserId == 0 {
		return nil
	}
	pId, err := request.GetIdFromReq(r)
	if err != nil {
		return nil
	}
	view := request.GetFilterValue(r, "view")
	if view == "" {
		return nil
	}
	data := models.PostData{
		BasePageData: base,
	}
	modData := ModerationData{
		Index:   index,
		ReqUI:   reqUId,
		CurrUID: currUId,
		PID:     pId,
		View:    view,
		Data:    data,
	}

	return &modData
}
