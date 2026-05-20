package users

import "user-profiles/internal/models"

type AccessData struct {
	Index  int
	ReqUID int
	Data   models.UserData
	View   string
}
