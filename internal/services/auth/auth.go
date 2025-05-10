package auth

import (
	"context"
	"errors"
	"gitlab.com/ictisagora/backend/internal/lib/jwt"
	"gitlab.com/ictisagora/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log      *slog.Logger
	repo     repository.Repository
	tokenTTL time.Duration
}

func New(log *slog.Logger, repo repository.Repository, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		repo:     repo,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	op := "AuthLogin"
	log := a.log.With(
		slog.String("op", op),
		slog.String("user", email),
	)

	log.Info("attempting to login user")

	user, err := a.repo.SelectUserByEmail(email)
	if err != nil {
		log.Error("user not found")
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		a.log.Error("incorrect password")
		return "", errors.New("incorrect password")
	}
	//TODO: make app provider

	log.Info("user logged in")

	token := jwt.NewToken(user, a.tokenTTL, "dasdbasdhas")

	return token, nil
}
