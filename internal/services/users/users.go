package users

import (
	"errors"
	"gitlab.com/ictisagora/backend/internal/models"
	"gitlab.com/ictisagora/backend/internal/repository"
	"log/slog"
)

type Users struct {
	log  slog.Logger
	repo repository.Repository
}

func New(log slog.Logger, repo repository.Repository) *Users {
	return &Users{
		log:  log,
		repo: repo,
	}
}

func (u *Users) GetUserByUID(UID string) (models.User, error) {
	op := "GetUserByUID"
	log := u.log.With(
		slog.String("op", op),
		slog.String("uid", UID),
	)

	if UID == "" {
		log.Error("uid is empty")
		return models.User{}, errors.New("uid is empty")
	}

	log.Debug("fetching user by uid")

	user, err := u.repo.SelectUserByUID(UID)
	if err != nil {
		log.Error("failed to fetch user by uid" + err.Error())
		return models.User{}, err
	}
	log.Info("successfully fetched user")

	return user, nil
}
