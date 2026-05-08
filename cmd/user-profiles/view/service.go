package view

func (s *DashboardService) GetAdminStats() (*Stats, error) {
    users, err := s.uRepo.GetAll()
	if err != nil {
		return nil, err
	}
	posts, err := s.pRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var Stats []Stats
	count := 0
	for _, user := range users {
		
	}
}