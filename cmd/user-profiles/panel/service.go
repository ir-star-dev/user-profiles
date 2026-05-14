package panel

import (
	"errors"
	"strings"
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
	"user-profiles/internal/validator"

	"golang.org/x/crypto/bcrypt"
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

func (s *DS) CreateUser(email, name, password, role string) ([]FormValidationErr, error) {
	var formValiErr []FormValidationErr
	form := CreateUserForm{
		Email:     email,
		Name:      name,
		Role:      role,
		Password:  password,
		Validator: validator.Validator{},
	}
	form.Validator.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email is not valid")
	form.Validator.CheckField(validator.MaxChars(form.Password, 8), "password", "Password length more than 8 characters")
	form.Validator.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.Validator.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.Validator.CheckField(validator.PermittedValue(form.Role, "user", "admin", "moderator"), "role", "Role can be user, admin or moderator")

	if !form.Validator.Valid() {
		for key, value := range form.Validator.FieldErrors {
			formValiErr = append(formValiErr, FormValidationErr{
				Name:    key,
				Message: value,
			})
		}
		return formValiErr, errors.New("Invalid form")
	}

	existedUser, err := s.uRepo.FindByEmail(email)
	if existedUser != nil {
		formValiErr = append(formValiErr, FormValidationErr{
			Name:    "exist",
			Message: "User already exists",
		})
		return formValiErr, errors.New("User already exists")
	}

	username, _, _ := strings.Cut(email, "@")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		formValiErr = append(formValiErr, FormValidationErr{
			Name:    "password",
			Message: "Password error",
		})
		return formValiErr, errors.New("Password error")
	}
	user := &users.UserWithRole{
		Name:     name,
		Username: username,
		Email:    email,
		Role:     role,
		Password: string(hash),
	}
	_, err = s.uRepo.Create(user)
	if err != nil {
		formValiErr = append(formValiErr, FormValidationErr{
			Name:    "form",
			Message: "Something went wrong. Try later, please!",
		})
		return formValiErr, err
	}
	return nil, nil
}
