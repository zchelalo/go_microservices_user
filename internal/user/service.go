package user

import (
	"context"
	"log"

	"github.com/zchelalo/go_microservices_domain/domain"
)

type (
	Filters struct {
		FirstName string
		LastName  string
	}

	Service interface {
		Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
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

func (srv *service) Create(ctx context.Context, firstName, lastName, email, phone string) (*domain.User, error) {
	srv.log.Println("create user service")
	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
	if err := srv.repository.Create(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (srv *service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	srv.log.Println("get all users service")
	users, err := srv.repository.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (srv *service) Get(ctx context.Context, id string) (*domain.User, error) {
	srv.log.Println("get user service")
	user, err := srv.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (srv *service) Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error {
	srv.log.Println("update user service")
	return srv.repository.Update(ctx, id, firstName, lastName, email, phone)
}

func (srv *service) Delete(ctx context.Context, id string) error {
	srv.log.Println("delete user service")
	return srv.repository.Delete(ctx, id)
}

func (srv *service) Count(ctx context.Context, filters Filters) (int, error) {
	srv.log.Println("count user service")
	return srv.repository.Count(ctx, filters)
}
