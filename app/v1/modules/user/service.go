package user

type Service interface {
	RegisterUser(input *AddUserForm) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input *AddUserForm) (User, error) {
	user := User{}
	user.Name = input.Name

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}
