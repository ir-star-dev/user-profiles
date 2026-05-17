package dashboard

import (
	"user-profiles/internal/auth"
	"user-profiles/internal/models"
	"user-profiles/internal/posts"
	"user-profiles/internal/users"
)

type DS struct {
	uService users.UsersService
	pService posts.PostService
	aService auth.AuthService
}

func NewDashboardService(uS users.UsersService, pS posts.PostService, aS auth.AuthService) *DS {
	return &DS{uService: uS, pService: pS, aService: aS}
}

func (s *DS) GetStats() (*models.Stats, error) {
	var uS []models.StatUsers
	uStats, err := s.uService.CountRoles()
	if err != nil {
		return nil, err
	}
	banned := 0
	for _, stat := range uStats {
		banned += stat.Banned
		uS = append(uS, models.StatUsers{
			Role:  stat.Role,
			Count: stat.Count,
			},
		)
	}
	pStats, err := s.pService.PostsStatus()
	if err != nil {
		return nil, err
	}
	stats := &models.Stats{
		UserStats: uS,
		Banned:    banned,
		PostStats: pStats,
	}
	return stats, nil
}

func (s *DS) GetAuthors() ([]models.Authors, error) {
	authors, err := s.pService.Authors()
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (s *DS) GetLoginsLog() ([]models.LoginsResponse, error) {
	res, err := s.aService.GetLoginLog()
	if err != nil {
		return nil, err
	}
	var logs []models.LoginsResponse
	for _, r := range res {
		logs = append(logs, models.LoginsResponse{
			Username:  r.Username,
			Device:    r.Device,
			Ip:        r.Ip,
			UserAgent: r.UserAgent,
			LoginAt:   r.LoginAt,
		})
	}
	return logs, nil
}
