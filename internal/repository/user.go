package repository

import (
	"mychat/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByUserName(username string) (*domain.User, error)
}
