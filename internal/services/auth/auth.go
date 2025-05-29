package auth

import (
	"errors"
	"github.com/google/uuid"
	"gitlab.com/ictisagora/backend/internal/lib/jwt"
	"gitlab.com/ictisagora/backend/internal/models"
	"gitlab.com/ictisagora/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log       slog.Logger
	repo      repository.Repository
	tokenTTL  time.Duration
	jwtSecret string
}

func New(log slog.Logger, repo repository.Repository, tokenTTL time.Duration, jwtSecret string) *Auth {
	return &Auth{
		log:       log,
		repo:      repo,
		tokenTTL:  tokenTTL,
		jwtSecret: jwtSecret,
	}
}

func (a *Auth) Register(u models.User) error {
	op := "AuthRegister"
	log := a.log.With(
		slog.String("op", op),
		slog.String("uid", u.UID),
	)

	log.Info("attempting to register user")

	hashPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")
		return err
	}

	u.Password = string(hashPass)

	uid := uuid.New().String()

	u.UID = uid

	err = a.repo.InsertUser(u)
	if err != nil {
		log.Error("failed to insert user")
		return err
	}
	log.Info("successfully registered user")

	return nil
}

func (a *Auth) Login(email string, password string) (string, string, error) {
	op := "AuthLogin"
	log := a.log.With(
		slog.String("op", op),
		slog.String("user", email),
	)

	log.Info("attempting to login user")

	user, err := a.repo.SelectUserByEmail(email)
	if err != nil {
		log.Error("error selecting user" + err.Error())
		return "", "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		a.log.Error("incorrect password")
		return "", "", errors.New("incorrect password")
	}
	//TODO: make app provider

	log.Info("user logged in")

	token := jwt.NewToken(user, a.tokenTTL, a.jwtSecret)

	return token, user.UID, nil
}

func (a *Auth) GetJWT() string {
	return a.jwtSecret
}
