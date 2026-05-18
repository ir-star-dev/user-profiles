package dashboard

import (
	"net/http"

	"user-profiles/internal/http/request"
	"user-profiles/internal/models"
	"user-profiles/internal/posts"
	"user-profiles/internal/templates"
	"user-profiles/internal/users"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

type DHandler struct {
	UService users.UsersService
	PService posts.PostService
	TCache   templates.Templates
	DService DS
}

type DHandlerDeps struct {
	UService users.UsersService
	PService posts.PostService
	TCache   templates.Templates
	DService DS
}

func NewDashboardHandler(router chi.Router, deps DHandlerDeps) *DHandler {
	return &DHandler{
		UService: deps.UService,
		PService: deps.PService,
		TCache:   deps.TCache,
		DService: deps.DService,
	}
}

func (h *DHandler) Profile(w http.ResponseWriter, r *http.Request) {
	currUserId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := h.UService.Role(currUserId)

	uIdFromReq, err := request.GetIdFromReq(r)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}
	data := models.PageData{
		RequestedUserId: uIdFromReq,
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}

	profile, err := h.UService.Get(uIdFromReq)
	if err != nil {
		h.TCache.PanelNotFound(w, r)
		return
	}

	uPosts, err := h.PService.CreatedPosts(uIdFromReq)
	if err != nil {
		uPosts = make([]models.PostRows, 0)
	}
	groupedP := GroupPosts(uPosts)

	var cards []models.UserData
	cards = append(cards, models.UserData{
		Profiles: models.UserViewTable{
			User:  *profile,
		},
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		Actions: models.Actions{
			CanDeleteUser: users.CanDeleteUser(currUserRole, currUserId, profile.Id),
			CanBanUser:    users.CanBanUser(currUserRole, currUserId, profile.Id),
		},
		CreatedPosts: groupedP,
	})

	data.UserCards = cards

	h.TCache.RenderPanel(w, r, http.StatusOK, "profile.tmpl", data)
}

func (h *DHandler) Statistics(w http.ResponseWriter, r *http.Request) {
	currUserId, err := request.GetUserId(r)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}
	currUserRole, _ := h.UService.Role(currUserId)

	data := models.PageData{
		CurrentUserId:   currUserId,
		CurrentUserRole: currUserRole,
		CSRFToken:       csrf.Token(r),
	}

	stats, err := h.DService.GetStats()
	if err != nil {
		stats = &models.Stats{}
	}
	data.Stats = stats

	logins, err := h.DService.GetLoginsLog()
	if err != nil {
		logins = []models.LoginsResponse{}
	}
	data.Logins = logins
	h.TCache.RenderPanel(w, r, http.StatusOK, "dashboard.tmpl", data)
}

func GroupPosts(posts []models.PostRows) models.GroupedPosts {
	var grouped models.GroupedPosts

	for _, post := range posts {
		if post.Approved {
			grouped.Published = append(grouped.Published, post)
		} else {
			grouped.Pending = append(grouped.Pending, post)
		}
	}
	return grouped
}