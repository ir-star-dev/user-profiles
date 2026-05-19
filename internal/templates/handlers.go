package templates

import (
	"net/http"
	"user-profiles/internal/http/request"
	"user-profiles/internal/models"

	"github.com/gorilla/csrf"
)

type BaseHandler struct {
	TCache *Templates
}

func (h *BaseHandler) NewBasePageData(r *http.Request) models.BasePageData {
	currentUserId, _ := request.GetUserId(r)
	currentUserRole, _ := request.GetUserRole(r.Context())

	return models.BasePageData{
		CurrentUserId:   currentUserId,
		CurrentUserRole: currentUserRole,
		CSRFToken:       csrf.Token(r),
	}
}
