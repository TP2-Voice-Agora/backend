package interfaces

import (
	"github.com/TP2-Voice-Agora/backend/internal/models"
	"mime/multipart"
)

type IdeaService interface {
	GetIdeaCategories() []models.IdeaCategory
	GetIdeaStatuses() []models.IdeaStatus
	GetAllIdeas() ([]models.Idea, error)
	GetIdeaByUID(uid string) (models.IdeaComment, error)
	GetAuthorIdeas(uid string, limit int) ([]models.Idea, error)
	InsertIdea(name string, text string, author string, status int, category int) (models.Idea, error)
	InsertComment(ideaUID, authorUID, commentText string) (models.Comment, error)
	InsertReply(commentUID, authorID, replyText string) (models.Reply, error)
	IncrementDislikes(ideaUID string) error
	IncrementLikes(ideaUID string) error
}

type AuthService interface {
	Register(u models.User) error
	Login(email string, password string) (string, string, error)
	GetJWT() string
}

type UserService interface {
	GetUserByUID(uid string) (models.User, error)
	UploadPFP(file multipart.File, header *multipart.FileHeader, UID string) (string, error)
	GetPositions() ([]models.UserPosition, error)
}
