package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/rajabhishekmaurya/ecom/internal/model"
	repo "github.com/rajabhishekmaurya/ecom/internal/repo"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, user *model.User) error {
	// Check if username already exists
	existingUser, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Save user
	return s.repo.Create(ctx, user)
}
