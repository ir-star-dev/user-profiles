package profile

import (
	"user-profiles/internal/domain/user"
)

type Service interface {
	View(uId int) (*user.User, error)
	UpdateName(uId int, name string) (*user.User, error)
	Delete(uId int) error
}

type profileService struct {
	uRepo user.Repository
}

func NewProfileService(uRepo user.Repository) Service {
	return &profileService{
		uRepo: uRepo,
	}
}

func (s *profileService) View(uId int) (*user.User, error) {
	user, err := s.uRepo.FindById(uId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *profileService) UpdateName(uId int, name string) (*user.User, error) {
	user, err := s.uRepo.FindById(uId)
	if user == nil {
		return nil, err
	}
	user.Name = name
	data, err := s.uRepo.UpdateName(user)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *profileService) Delete(uId int) error {
	user, err := s.uRepo.FindById(uId)
	if user == nil {
		return err
	}
	
	err = s.uRepo.Delete(uId)
	if err != nil {
		return err
	}
	return nil
}