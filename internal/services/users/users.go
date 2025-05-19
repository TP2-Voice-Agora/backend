package users

import (
	"errors"
	"fmt"
	"github.com/TP2-Voice-Agora/backend/internal/models"
	"github.com/TP2-Voice-Agora/backend/internal/repository"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
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

func (u *Users) UploadPFP(
	file multipart.File,
	header *multipart.FileHeader,
	UID string) (string, error) {

	isValid := func(h *multipart.FileHeader) bool {
		ext := strings.ToLower(filepath.Ext(h.Filename))
		return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
	}

	if !isValid(header) {
		return "", errors.New("invalid file type, only jpeg and png allowed")
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("user_%s%s", UID, ext)
	filePath := filepath.Join("uploads", filename)

	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", errors.New("could not create or overwrite file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", errors.New("failed to write file")
	}

	err = u.repo.UpdateUserPfpURL(UID, filePath)
	if err != nil {
		u.log.Error("failed to update user pfp url", slog.String("error", err.Error()))
		return "", err
	}

	return filePath, nil
}
