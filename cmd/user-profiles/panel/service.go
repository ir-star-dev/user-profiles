package panel

import (
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
)

type DS struct {
	uRepo users.Repository
	pRepo posts.Repository
}

type Stats struct {
	UserStats []userStats
	Banned    int
	PostStats []postStats
}

type userStats struct {
	Role  string
	Count int
}

type postStats struct {
	Published int
	Pending   int
	Authors   []string
}

type statUsers struct {
	Count  int    `db:"count"`
	Role   string `db:"role"`
	Banned int    `db:"banned"`
}

func NewDashboardService(uRepo users.Repository, pRepo posts.Repository) *DS {
	return &DS{uRepo: uRepo, pRepo: pRepo}
}

func (s *DS) GetStats() (*Stats, error) {
	var uS []userStats
	uStats, err := s.uRepo.CountRoles()
	if err != nil {
		return nil, err
	}
	banned := 0
	for _, stat := range uStats {
		banned += stat.Banned
		uS = append(uS, userStats{
			Role:  stat.Role,
			Count: stat.Count,
		},
		)
	}
	pStats, err := s.pRepo.PostsStatus()
	if err != nil {
		return nil, err
	}
	var pS []postStats
	for _, stat := range pStats {
		pS = append(pS, postStats{
			Published: stat.Published,
			Pending:   stat.Pending,
		},
		)
	}

	stats := &Stats{
		UserStats: uS,
		Banned:    banned,
		PostStats: pS,
	}
	return stats, nil
}

func (s *DS) GetAuthors() ([]posts.Authors, error) {
	authors, err := s.pRepo.Authors()
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (s *DS) GetRoles() ([]users.Roles, error) {
	roles, err := s.uRepo.Roles()
	if err != nil {
		return nil, err
	}
	return roles, nil
}