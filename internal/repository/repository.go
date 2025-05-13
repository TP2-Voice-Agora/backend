package repository

import (
	"gitlab.com/ictisagora/backend/internal/models"
	"log/slog"
)

// Repository - interface to work with DB
type Repository interface {
	ConnectDB(sourceURL string, log slog.Logger) error
	CloseConnectDB() error

	InsertUser(models.User) error
	SelectUserByEmail(string) (models.User, error)

	SelectPositions() ([]models.UserPosition, error)

	InsertIdea(models.Idea) error
	SelectIdeas() ([]models.Idea, error)
	SelectIdeaByUID(uid string) (models.Idea, error)
	SelectUserIdeas(string, int) ([]models.Idea, error)
	InsertIdeaComment(models.Comment) error
	InsertCommentReply(models.Reply) error
	SelectIdeaComments(string) ([]models.Comment, error)
	SelectCommentReplies(string) ([]models.Reply, error)

	SelectIdeaCategories() ([]models.IdeaCategory, error)
	SelectIdeaStatuses() ([]models.IdeaStatus, error)
}
