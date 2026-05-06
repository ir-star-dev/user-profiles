package view

import "user-profiles/cmd/user-profiles/users"

type PageData struct {
	RequestedUserId int
	CurrentUserId   int
	CurrentUserRole string
	Cards           []UserCardData
	Pagination      PaginationData
}

type UserCardData struct {
	Profiles        users.UserWithRole
	CanDelete       bool
	CanBan          bool
	CurrentUserId   int
	CurrentUserRole string
}

type PaginationData struct {
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