package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/internal/users"
)

type templateData struct {
	RequestedUserId int
	CurrentUserId   int
	CurrentUserRole string
	User            []userData
	Pagin           paginationData
}

type userData struct {
	Profiles        users.UserWithRole
	CanDelete       bool
	CanBan          bool
	CurrentUserId   int
	CurrentUserRole string
}

type paginationData struct {
	Page       int
	TotalPages int
	Pages      []int
	HasPrev    bool
	HasNext    bool
	PrevPage   int
	NextPage   int
	ShowDots   bool
	LastPage   int
}

func (app *application) getIdFromReq(r *http.Request) (int, error) {
	idStr := strings.TrimSpace(r.PathValue("id"))
	if idStr == "" {
		return 0, errors.New("Missing param")
	}
	uId, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return uId, nil
}

func (app *application) getPageFromReq(r *http.Request) (int, error) {
	page := strings.TrimSpace(r.PathValue("page"))
	if page == "" {
		return 0, errors.New("Missing param")
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return 0, err
	}
	return p, nil
}

func (app *application) selfDeletionDetected(r *http.Request) (int, int, error) {
	reqId, err := app.getIdFromReq(r)
	if err != nil {
		return 0, 0, err
	}
	currId, err := utils.GetUserID(r.Context())
	if err != nil {
		return reqId, currId, err
	}
	if reqId == currId {
		return reqId, currId, errors.New("You can't delete yourself.")
	}
	return reqId, currId, nil
}

func buildPagination(currentPage, totalPages int) paginationData {
	pagin := paginationData{
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
		pagin.Pages = append(pagin.Pages, i)
	}

	if end < totalPages {
		pagin.ShowDots = true
	}

	return pagin
}

func canDelete(currentRole string, currentUserId, profileId int) bool {
	if currentRole == "admin" {
		return currentUserId != profileId
	}
	return currentUserId == profileId
}

func canBan(currentRole string, currentUserId, profileId int) bool {
	return currentRole == "admin" && currentUserId != profileId
}
