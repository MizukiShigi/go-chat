package postgresql

import (
	"mychat/internal/domain"
	"mychat/internal/repository"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) Create(user *domain.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepositoryImpl) FindByUserName(username string) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
