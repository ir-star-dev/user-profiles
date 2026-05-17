package models

import (
	"user-profiles/internal/validator"
)

type PageData struct {
	RequestedUserId   int
	CurrentUserId     int
	CurrentUserRole   string
	UserCards         []UserData
	PostCards         []PostData
	Pagination        Pagination
	Loadmore          Loadmore
	Stats             *Stats
	Authors           []Authors
	Filters           map[string]string
	HasFilters        bool
	Roles             []Roles
	FormValidationErr []validator.FormValidationErr
	UserCredentials   UserCredentials
	PostCreated       PostCreated
	PostUpdated       PostUpdated
	Logins            []LoginsResponse
}

type UserData struct {
	Profiles        UserWithRole
	CurrentUserId   int
	CurrentUserRole string
	Actions         Actions
	CreatedPosts    []PostRows
}

type PostData struct {
	Posts   PostWithUserName
	Actions Actions
}

type Actions struct {
	CanDeletePost  bool
	CanApprovePost bool
	CanEditPost    bool
	CanDeleteUser  bool
	CanBanUser     bool
}

type Pagination struct {
	Path       string
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

type Loadmore struct {
	Page    int
	Next    int
	HasMore bool
	Total   int
}

type UserCredentials struct {
	Email    string
	Password string
}

type PostCreated struct {
	Message string
}

type PostUpdated struct {
	Message string
}

type Stats struct {
	UserStats []StatUsers
	Banned    int
	PostStats []PostsStatus
}
