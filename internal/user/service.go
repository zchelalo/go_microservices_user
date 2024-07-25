package user

import (
	"log"

	"github.com/zchelalo/go_microservices_domain/domain"
)

type (
	Filters struct {
		FirstName string
		LastName  string
	}

	Service interface {
		Create(firstName, lastName, email, phone string) (*domain.User, error)
		GetAll(filters Filters, offset, limit int) ([]domain.User, error)
		Get(id string) (*domain.User, error)
		Update(id string, firstName, lastName, email, phone *string) error
		Delete(id string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log        *log.Logger
		repository Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:        log,
		repository: repo,
	}
}

func (srv *service) Create(firstName, lastName, email, phone string) (*domain.User, error) {
	srv.log.Println("create user service")
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
	if err := srv.repository.Create(&user); err != nil {
		// srv.log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (srv *service) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	srv.log.Println("get all users service")
	users, err := srv.repository.GetAll(filters, offset, limit)
	if err != nil {
		// srv.log.Println(err)
		return nil, err
	}
	return users, nil
}

func (srv *service) Get(id string) (*domain.User, error) {
	srv.log.Println("get user service")
	user, err := srv.repository.Get(id)
	if err != nil {
		// srv.log.Println(err)
		return nil, err
	}
	return user, nil
}

func (srv *service) Update(id string, firstName, lastName, email, phone *string) error {
	srv.log.Println("update user service")
	return srv.repository.Update(id, firstName, lastName, email, phone)
}

func (srv *service) Delete(id string) error {
	srv.log.Println("delete user service")
	return srv.repository.Delete(id)
}

func (srv *service) Count(filters Filters) (int, error) {
	srv.log.Println("count user service")
	return srv.repository.Count(filters)
}
