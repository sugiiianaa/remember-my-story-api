package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	repositories "github.com/sugiiianaa/remember-my-story/internal/Repositories"
	"github.com/sugiiianaa/remember-my-story/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(email, fullName, password string) (uint, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return 0, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		Email:    email,
		FullName: fullName,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.generateJWTToken(user)
}

func (s *AuthService) generateJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", user.ID),
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}
