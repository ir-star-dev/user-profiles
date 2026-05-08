package users

import (
	"errors"
)

type usersService struct {
	uRepo Repository
}

func NewUsersService(uRepo Repository) UsersService {
	return &usersService{uRepo: uRepo}
}

func (s *usersService) Get(uId int) (*UserWithRole, error) {
	user, err := s.uRepo.FindById(uId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *usersService) ChangeName(uId int, name *string) (*UserWithRole, error) {
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

func (s *usersService) Delete(uId int) error {
	err := s.uRepo.Delete(uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *usersService) GetOnPage(page int) ([]UserWithRole, int, error) {
	users, res, err := s.uRepo.GetOnPage(page)
	if err != nil {
		return nil, 0, err
	}
	return users, res, nil
}

func (s *usersService) GetAll() ([]UserWithRole, error) {
	users, err := s.uRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *usersService) Ban(uId int) error {
	err := s.uRepo.Ban(uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *usersService) Unban(uId int) error {
	err := s.uRepo.Unban(uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *usersService) Role(uId int) (string, error) {
	role, err := s.uRepo.FindRoleByUserId(uId)
	if err != nil {
		return "", err
	}
	return role, nil
}