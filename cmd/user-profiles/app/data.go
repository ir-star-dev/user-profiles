package app

import (
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
)

type PageData struct {
	RequestedUserId int
	CurrentUserId   int
	CurrentUserRole string
	UserCards       []UserData
	PostCards       []PostData
	Pagination      Pagination
	Loadmore        Loadmore
	Stats           *Stats
}

type UserData struct {
	Profiles        users.UserWithRole
	CanDelete       bool
	CanBan          bool
	CurrentUserId   int
	CurrentUserRole string
}

type Pagination struct {
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

type PostData struct {
	Posts          posts.PostWithUserName
	CanDeletePost  bool
	CanApprovePost bool
	CanEditPost    bool
}

type Loadmore struct {
	Page    int
	Next    int
	HasMore bool
	Total   int
}
