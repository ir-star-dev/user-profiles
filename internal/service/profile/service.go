package profile_service

import (
	"user-profiles/internal/domain/user"
	"errors"
)

type profileService struct {
	uRepo user.Repository
}

func New(uRepo user.Repository) Service {
	return &profileService{
		uRepo: uRepo,
	}
}

func (s *profileService) View(uId int) (*user.UserWithRole, error) {
	user, err := s.uRepo.FindById(uId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *profileService) ViewAll() ([]*user.UserWithRole, error) {
	users, err := s.uRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *profileService) ChangeName(uId int, name *string) (*user.UserWithRole, error) {
	user, err := s.uRepo.FindById(uId)
	if user == nil {
		return nil, err
	}
	if name == nil {
		return nil, errors.New(MissingName)
	}
	if len(*name) < 2 {
		return nil, errors.New(ShortName)
	}
	data, err := s.uRepo.UpdateName(*name, uId)
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

func (s *profileService) Ban(uId int) error {
	user, err := s.uRepo.FindById(uId)
	if user == nil {
		return err
	}
	err = s.uRepo.Ban(uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *profileService) Unban(uId int) error {
	user, err := s.uRepo.FindById(uId)
	if user == nil {
		return err
	}
	err = s.uRepo.Unban(uId)
	if err != nil {
		return err
	}
	return nil
}