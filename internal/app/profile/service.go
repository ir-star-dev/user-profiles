package profile

import (
	"user-profiles/internal/domain/user"
	"errors"
)

type Service interface {
	View(uId int) (*user.User, error)
	ViewAll() ([]*user.User, error)
	ChangeName(uId int, name *string) (*user.User, error)
	Delete(uId int) error
	Ban(uId int) error
	Unban(uId int) error
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

func (s *profileService) ViewAll() ([]*user.User, error) {
	users, err := s.uRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *profileService) ChangeName(uId int, name *string) (*user.User, error) {
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