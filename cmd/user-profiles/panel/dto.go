package panel

import (
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
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
	Authors           []posts.Authors
	Filters           map[string]string
	HasFilters        bool
	Roles             []users.Roles
	FormValidationErr []FormValidationErr
	UserCredentials	  UserCredentials
}

type UserData struct {
	Profiles        users.UserWithRole
	CurrentUserId   int
	CurrentUserRole string
	Actions         Actions
}

type PostData struct {
	Posts   posts.PostWithUserName
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

type FormValidationErr struct {
	Name    string
	Message string
}

type CreateUserForm struct {
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	Name      string              `json:"name"`
	Role      string              `json:"role"`
	Validator validator.Validator `json:"-"`
}

type UserCredentials struct {
	Email    string
	Password string
}
