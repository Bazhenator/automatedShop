package services

import (
	customErr "automatedShop/internal/errors"
	"automatedShop/internal/repository"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"unicode/utf8"
)

type AuthService struct {
	l        *slog.Logger
	AuthRepo repository.IAuthRepository
}

func NewAuthService(repo repository.IAuthRepository) *AuthService {
	var l *slog.Logger

	return &AuthService{
		l:        l,
		AuthRepo: repo,
	}
}

// AuthoriseUser checks if user with given credentials exists in the system and returns bool result.
//
// If user exists in system and password is correct returns true.
// If user exists, but password is incorrect, returns false.
// If user doesn't exist, returns false.
func (s *AuthService) AuthoriseUser(ctx context.Context, login string, pwd string) bool {
	const op = "AuthService.AuthoriseUser"

	user, err := s.AuthRepo.FindUser(ctx, login)
	if err != nil {
		fmt.Printf("error occurred in: %v: %v\n", op, err)
		return false
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(pwd)); err != nil {
		fmt.Printf("error occurred in: %v: %v\n", op, err)

		return false
	}

	fmt.Printf("%v: user logged in successfully!\n", op)
	return true
}

// RegisterUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (s *AuthService) RegisterUser(ctx context.Context, login string, pass string) error {
	const op = "Auth.RegisterNewUser"

	if utf8.RuneCountInString(pass) < 5 {
		return fmt.Errorf("%s: %w", op, customErr.ErrTooSmallPwdLen)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.AuthRepo.FindUser(ctx, login)
	if err == nil {
		return errors.New("user already exists")
	}

	err = s.AuthRepo.SaveUser(ctx, login, passHash)
	if err != nil {
		//log.Error("failed to save user", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// IsRootUser checks if user is admin.
func (s *AuthService) IsRootUser(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsRoot"

	log := s.l.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is root")

	isRoot, err := s.AuthRepo.IsRoot(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is root", slog.Bool("is_root", isRoot))

	return isRoot, nil
}
