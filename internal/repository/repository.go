package repository

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ictisagora/backend/internal/models"
)

// Repository - interface to work with DB
type Repository interface {
	ConnectDB(string) (*sqlx.DB, error)
	CloseConnectDB() error

	InsertUser(models.User) error
	SelectUserByEmail(string) (models.User, error)

	SelectPositions() ([]models.UserPosition, error)

	InsertIdea(models.Idea) error
	SelectIdeas() ([]models.Idea, error)
	SelectUserIdeas(string, int) ([]models.Idea, error)
	InsertIdeaComment(models.Comment) error
	InsertCommentReply(models.Reply) error
	SelectIdeaComments(string) ([]models.Comment, error)
	SelectCommentReplies(int) ([]models.Reply, error)

	SelectIdeaCategories() ([]models.IdeaCategory, error)
	SelectIdeaStatuses() ([]models.IdeaStatus, error)
}
