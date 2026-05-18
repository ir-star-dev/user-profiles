package users

import (
	"errors"
	"strconv"
	"strings"
	"user-profiles/internal/models"
	"user-profiles/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	uRepo Repository
}

func NewUsersService(uRepo Repository) UsersService {
	return &usersService{uRepo: uRepo}
}

func (s *usersService) CreateUser(email, name, password, role string) ([]validator.FormValidationErr, error) {
	var formValiErr []validator.FormValidationErr
	var v validator.Validator
	v.CheckField(validator.Matches(email, validator.EmailRX), "email", EmailNotValid)
	v.CheckField(validator.MaxChars(password, 8), "password", LongPass)
	v.CheckField(validator.NotBlank(password), "password", PassEmpty)
	v.CheckField(validator.NotBlank(name), "name", NameEmpty)
	v.CheckField(validator.PermittedValue(role, "user", "admin", "moderator"), "role", WrongRole)

	if !v.Valid() {
		for key, value := range v.FieldErrors {
			formValiErr = append(formValiErr, validator.FormValidationErr{
				Name:    key,
				Message: value,
			})
		}
		return formValiErr, errors.New(InvalidForm)
	}

	existedUser, err := s.uRepo.FindByEmail(email)
	if existedUser != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "exist",
			Message: UserExists,
		})
		return formValiErr, errors.New(UserExists)
	}
	username, _, _ := strings.Cut(email, "@")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "password",
			Message: PasswordError,
		})
		return formValiErr, errors.New(PasswordError)
	}
	user := &models.UserWithRole{
		Name:     name,
		Username: username,
		Email:    email,
		Role:     role,
		Password: string(hash),
	}
	_, err = s.uRepo.Create(user)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return formValiErr, err
	}
	return nil, nil
}

func (s *usersService) Get(uId int) (*models.UserWithRole, error) {
	user, err := s.uRepo.FindById(uId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *usersService) ChangeName(uId int, name string) ([]validator.FormValidationErr, *string, error) {
	var formValiErr []validator.FormValidationErr
	var v validator.Validator
	v.CheckField(validator.MinChars(name, 3), "name", ShortName)
	v.CheckField(validator.NotBlank(name), "name", NameEmpty)
	if !v.Valid() {
		for key, value := range v.FieldErrors {
			formValiErr = append(formValiErr, validator.FormValidationErr{
				Name:    key,
				Message: value,
			})
		}
		return formValiErr, nil, errors.New(InvalidForm)
	}

	data, err := s.uRepo.UpdateName(name, uId)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name: "update",
			Message: SomethingWrong,
		})
		return formValiErr, nil, errors.New(SomethingWrong)
	}
	id := strconv.Itoa(*data)
	return nil, &id, nil
}

func (s *usersService) Delete(uId int) error {
	err := s.uRepo.Delete(uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *usersService) GetOnPage(page, limit int, banned *bool, role, search string) ([]models.UserWithRole, int, error) {
	if len(search) < 2 {
		search = ""
	}
	if len(search) > 100 {
		search = search[:100]
	}
	users, res, err := s.uRepo.GetOnPage(page, limit, banned, role, search)
	if err != nil {
		return nil, 0, err
	}
	return users, res, nil
}

func (s *usersService) Role(uId int) (string, error) {
	role, err := s.uRepo.FindRoleByUserId(uId)
	if err != nil {
		return "", err
	}
	return role, nil
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

func (s *usersService) CountRoles() ([]models.StatUsers, error) {
	stat, err := s.uRepo.CountRoles()
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (s *usersService) GetRoles() ([]models.Roles, error) {
	role, err := s.uRepo.Roles()
	if err != nil {
		return nil, err
	}
	return role, nil
}
