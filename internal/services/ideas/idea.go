package ideas

import (
	"gitlab.com/ictisagora/backend/internal/models"
	"gitlab.com/ictisagora/backend/internal/repository"
	"log/slog"
)

type Ideas struct {
	log             slog.Logger
	repo            repository.Repository
	ideasCategories []models.IdeaCategory
	ideasStatuses   []models.IdeaStatus
}

func New(log slog.Logger, repo repository.Repository) *Ideas {
	i := &Ideas{
		log:  log,
		repo: repo,
	}

	var err error

	i.ideasCategories, err = i.repo.SelectIdeaCategories()
	if err != nil {
		log.Error("failed to fetch ideas categories" + err.Error())
		return nil
	}

	i.ideasStatuses, err = i.repo.SelectIdeaStatuses()
	if err != nil {
		log.Error("failed to fetch ideas statuses" + err.Error())
	}

	return i
}

func (i *Ideas) GetIdeaCategories() []models.IdeaCategory {
	return i.ideasCategories
}

func (i *Ideas) GetIdeaStatuses() []models.IdeaStatus {
	return i.ideasStatuses
}

func (i *Ideas) GetAllIdeas() ([]models.Idea, error) {
	op := "IdeasGetAll"
	log := i.log.With(slog.String("op", op))
	log.Debug("fetching all ideas")

	ideas, err := i.repo.SelectIdeas()
	if err != nil {
		log.Error("failed to fetch ideas" + err.Error())
		return nil, err
	}

	log.Info("successfully fetched ideas")

	return ideas, nil
}

func (i *Ideas) GetAuthorIdeas(uid string, limit int) ([]models.Idea, error) {
	op := "IdeasGetAuthorIdeas"
	log := i.log.With(
		slog.String("op", op),
		slog.String("uid", uid),
		slog.Int("limit", limit),
	)
	log.Debug("fetching all author ideas")

	ideas, err := i.repo.SelectUserIdeas(uid, limit)

	if err != nil {
		log.Error("failed to fetch author ideas" + err.Error())
		return nil, err
	}

	log.Info("successfully fetched author ideas")

	return ideas, nil
}

func (i *Ideas) GetIdeaComments(uid string) ([]models.Comment, error) {
	op := "IdeasGetIdeaComments"
	log := i.log.With(
		slog.String("op", op),
		slog.String("uid", uid),
	)
	log.Debug("fetching all idea comments")

	ideas, err := i.repo.SelectIdeaComments(uid)

	if err != nil {
		log.Error("failed to fetch idea comments" + err.Error())
		return nil, err
	}

	log.Info("successfully fetched idea comments")

	return ideas, nil
}
