package posts

import "user-profiles/internal/models"

type ModerationData struct {
	Index   int
	ReqUI   int
	PID     int
	CurrUID int
	View    string
	Data    models.PostData
}
