package users

import (
	"errors"
	"strconv"
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

func (s *usersService) ChangeName(uId int, name string) (*string, error) {
	if len(name) < 2 {
		return nil, errors.New(ShortName)
	}
	data, err := s.uRepo.UpdateName(name, uId)
	if err != nil {
		return nil, err
	}
	id := strconv.Itoa(*data)
	return &id, nil
}

func (s *usersService) Delete(uId int) error {
	err := s.uRepo.Delete(uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *usersService) GetOnPage(page int, limit int, banned *bool, role string) ([]UserWithRole, int, error) {
	users, res, err := s.uRepo.GetOnPage(page, limit, banned, role)
	if err != nil {
		return nil, 0, err
	}
	return users, res, nil
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
