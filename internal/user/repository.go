package user

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/zchelalo/go_microservices_domain/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repository struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{
		log: log,
		db:  db,
	}
}

func (repo *repository) Create(ctx context.Context, user *domain.User) error {
	if err := repo.db.WithContext(ctx).Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("user created with id: ", user.Id)
	return nil
}

func (repo *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error) {
	var users []domain.User

	tx := repo.db.WithContext(ctx).Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	// if err := repo.db.Model(&users).Select("id, first_name, email, created_at").Order("created_at desc").Find(&users).Error; err != nil {
	if err := tx.Order("created_at desc").Find(&users).Error; err != nil {
		repo.log.Println(err)
		return nil, err
	}

	return users, nil
}

func (repo *repository) Get(ctx context.Context, id string) (*domain.User, error) {
	user := domain.User{
		Id: id,
	}

	if err := repo.db.WithContext(ctx).Model(&user).First(&user).Error; err != nil {
		repo.log.Println(err)
		return nil, err
	}

	return &user, nil
}

func (repo *repository) Update(ctx context.Context, id string, firstName, lastName, email, phone *string) error {
	values := make(map[string]interface{})

	if firstName != nil {
		values["first_name"] = *firstName
	}

	if lastName != nil {
		values["last_name"] = *lastName
	}

	if email != nil {
		values["email"] = *email
	}

	if phone != nil {
		values["phone"] = *phone
	}

	if err := repo.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	return nil
}

func (repo *repository) Delete(ctx context.Context, id string) error {
	user := domain.User{
		Id: id,
	}

	if err := repo.db.WithContext(ctx).Model(&user).Delete(&user).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	return nil
}

func (repo *repository) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(&domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}

	return tx
}
