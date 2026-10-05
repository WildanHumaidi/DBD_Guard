package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (string, *domain.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	// 1. Find user by email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("email atau kata sandi salah")
	}

	// 2. Check if user is active (LOGIC-06)
	if !user.ActiveStatus {
		return "", nil, errors.New("akun anda tidak aktif, silakan hubungi administrator")
	}

	// 3. Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("email atau kata sandi salah")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role.RoleName,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	if user.LocationID != nil {
		claims["location_id"] = user.LocationID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Env.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
