package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rajabhishekmaurya/ecom/internal/repo"
	"github.com/rajabhishekmaurya/ecom/libs/config"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repo.UserRepository
	cfg  *config.Config
}
type LoginRequest struct {
	Username string
	Password string
}
type LoginResponse struct {
	Token string
}

func NewAuthService(repo repo.UserRepository, cfg *config.Config) *AuthService {

	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {

	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(s.cfg.JWT.Secret),
	)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
	}, nil
}
