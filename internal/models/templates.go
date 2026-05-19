package models

import (
	"user-profiles/internal/validator"
)

// General for all pages
type BasePageData struct {
	CurrentUserId   int
	CurrentUserRole string
	CSRFToken       string
}

// Auth
type AuthData struct {
	BasePageData
	FormValidationErr []validator.FormValidationErr
}

// ===================== Front =======================
// Home
type HomeData struct {
	BasePageData
	PostCards  []PostData
	Loadmore   Loadmore
	Filters    map[string]string
	HasFilters bool
}

// Post
type PostSimpleData struct {
	BasePageData
	PostCards []PostData
}

// User Posts
type UserPostsData struct {
	BasePageData
	PostCards  []PostData
	Filters    map[string]string
	HasFilters bool
}

// ===================== Front =======================

// Dashboard
type ProfileData struct {
	BasePageData
	UserCards       []UserData
	RequestedUserId int
}

type StaticticsData struct {
	BasePageData
	Stats  *Stats
	Logins []LoginsResponse
}

type PostsDashboardData struct {
	BasePageData
	PostCards         []PostData
	Pagination        Pagination
	Authors           []Authors
	Filters           map[string]string
	HasFilters        bool
	FormValidationErr []validator.FormValidationErr
	TotalPosts        int
	Page              int
	Pages             int
	Stats             *Stats
}

type CreatePostFormData struct {
	BasePageData
	FormValidationErr []validator.FormValidationErr
	PostCreated       PostCreated
}

type PreviewPostData struct {
	BasePageData
	PostCards []PostData
}

type PostDeleteData struct {
	BasePageData
	PostCards []PostData
}

type EditPostFormData struct {
	BasePageData
	FormValidationErr []validator.FormValidationErr
	PostUpdated       PostUpdated
	PostCards         []PostData
}

type UsersDashboardData struct {
	BasePageData
	UserCards  []UserData
	Pagination Pagination
	Stats      *Stats
	Filters    map[string]string
	HasFilters bool
	Roles      []Roles
	TotalUsers int
	Page       int
	Pages      int
}

type UpdateNameModalData struct {
	BasePageData
	UserCards []UserData
}

type UserDeleteData struct {
	BasePageData
	UserCards []UserData
}

type CreateUserFormData struct {
	BasePageData
	FormValidationErr []validator.FormValidationErr
	UserCredentials   UserCredentials
}

// =====================================

type UserData struct {
	BasePageData
	Profiles        UserViewTable
	Actions         Actions
	CreatedPosts    GroupedPosts
	CurrentUserId   int
	CurrentUserRole string
}

type PostData struct {
	Posts   PostViewTable
	Actions Actions
}

type PostViewTable struct {
	Post  PostWithUserName
	Index int
}

type UserViewTable struct {
	User  UserWithRole
	Index int
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

type GroupedPosts struct {
	Published []PostRows
	Pending   []PostRows
}
