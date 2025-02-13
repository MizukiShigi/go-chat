package usecase

import (
	"context"
	"mychat/internal/domain"
	"mychat/internal/infrastructure/redis"
	"mychat/internal/repository"
)

type AuthUsecase struct {
	userRepository repository.UserRepository
	redis *redis.Session
}

func NewAuthUsecase(userRepository repository.UserRepository, redis *redis.Session) *AuthUsecase {
	return &AuthUsecase{
		userRepository: userRepository,
		redis: redis,
	}
}

func (au *AuthUsecase) Register(user *domain.User) error {
	err := user.HashPassword()
	if err != nil {
		return err
	}

	return au.userRepository.Create(user)
}

func (au *AuthUsecase) Login(username, password string) (string, error) {
	user, err := au.userRepository.FindByUserName(username)
	if err != nil {
		return "", err
	}

	err = user.CheckPassword(password)
	if err != nil {
		return "", err
	}

	return au.redis.CreateSession(context.Background(), user.ID)
}