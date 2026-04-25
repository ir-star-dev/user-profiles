package users

// import (
// 	"errors"
// )

// type adminService struct {
// 	uRepo Repository
// }

// func NewAdminService(uRepo Repository) AdminService {
// 	return &adminService{uRepo: uRepo}
// }

// func (s *adminService) View(uId int) (*UserWithRole, error) {
// 	user, err := s.uRepo.FindById(uId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// // func (s *adminService) ViewAll() ([]*UserWithRole, error) {
// // 	users, err := s.uRepo.GetAll()
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return users, nil
// // }

// func (s *adminService) ChangeNameById(uId int, name *string) (*UserWithRole, error) {
// 	if name == nil {
// 		return nil, errors.New(MissingName)
// 	}
// 	if len(*name) < 2 {
// 		return nil, errors.New(ShortName)
// 	}
// 	data, err := s.uRepo.UpdateName(*name, uId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

// func (s *adminService) DeleteById(uId int) error {
// 	err := s.uRepo.Delete(uId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

