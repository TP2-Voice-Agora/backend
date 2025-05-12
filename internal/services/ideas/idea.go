package ideas

import (
	"github.com/google/uuid"
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

func (i *Ideas) GetIdeaByUID(uid string) (models.IdeaComment, error) {
	op := "IdeaGetByUID"
	log := i.log.With(
		slog.String("op", op),
		slog.String("uid", uid),
	)

	if uid == "" {
		log.Error("uid is empty")
		return models.IdeaComment{}, nil
	}

	log.Debug("fetching idea by uid")

	idea, err := i.repo.SelectIdeaByUID(uid)
	if err != nil {
		log.Error("failed to fetch idea by uid" + err.Error())
		return models.IdeaComment{}, err
	}
	log.Info("successfully fetched idea")
	log.Info("fetching comments for idea")

	comments, err := i.repo.SelectIdeaComments(idea.IdeaUID)
	if err != nil {
		log.Error("failed to fetch comments for idea" + err.Error())
		return models.IdeaComment{}, err
	}

	commentsReplies := make([]models.CommentReply, len(comments))

	// putting together comments and replies
	for j, comment := range comments {
		log.Info("fetching replies for comment")

		replies, err := i.repo.SelectCommentReplies(comment.CommentUID)
		if err != nil {
			log.Error("failed to fetch replies for comment" + comment.CommentUID + err.Error())
		}
		commentsReplies[j] = models.CommentReply{
			Comment: comment,
			Replies: replies,
		}
	}

	ideaComment := models.IdeaComment{
		Idea:           idea,
		CommentReplies: commentsReplies,
	}

	return ideaComment, nil
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

func (i *Ideas) GetCommentReplies(uid string) ([]models.Reply, error) {
	op := "IdeasGetCommentReplies"
	log := i.log.With(
		slog.String("op", op),
		slog.String("uid", uid),
	)
	log.Debug("fetching all idea replies")

	replies, err := i.repo.SelectCommentReplies(uid)
	if err != nil {
		log.Error("failed to fetch idea replies" + err.Error())
		return nil, err
	}

	log.Info("successfully fetched comment replies")

	return replies, nil
}

func (i *Ideas) InsertIdea(name, text, author string, status, category int) (models.Idea, error) {
	op := "IdeasInsertIdea"
	log := i.log.With(slog.String("op", op),
		slog.String("name", name),
	)
	log.Debug("inserting idea")

	if name == "" {
		log.Error("idea name is null")
		return models.Idea{}, nil
	}
	if text == "" {
		log.Error("idea text is null")
		return models.Idea{}, nil
	}
	if author == "" {
		log.Error("idea author is null")
		return models.Idea{}, nil
	}

	ideaUID := uuid.New().String()

	//timestamp - in PSQL
	idea := models.Idea{
		IdeaUID:    ideaUID,
		Name:       name,
		Text:       text,
		Author:     author,
		StatusID:   status,
		CategoryID: category,
	}

	err := i.repo.InsertIdea(idea)
	if err != nil {
		log.Error("failed to insert idea" + err.Error())
		return models.Idea{}, err
	}

	log.Info("successfully inserted idea, UUID:" + idea.IdeaUID)
	return idea, nil
}

func (i *Ideas) InsertComment(ideaUID, authorUID, commentText string) (models.Comment, error) {
	op := "IdeasInsertComment"
	log := i.log.With(slog.String("op", op),
		slog.String("ideaUID", ideaUID),
	)
	log.Debug("inserting comment")

	if ideaUID == "" {
		log.Error("ideaUID is null")
		return models.Comment{}, nil
	}
	if authorUID == "" {
		log.Error("authorUID is null")
		return models.Comment{}, nil
	}
	if commentText == "" {
		log.Error("commentText is null")
		return models.Comment{}, nil
	}

	commentUID := uuid.New().String()

	//timestamp - in PSQL
	comment := models.Comment{
		CommentUID:  commentUID,
		IdeaUID:     ideaUID,
		CommentText: commentText,
	}

	err := i.repo.InsertIdeaComment(comment)
	if err != nil {
		log.Error("failed to insert comment" + err.Error())
		return models.Comment{}, err
	}

	log.Info("successfully inserted comment, UUID:" + comment.CommentUID)

	return comment, nil
}

func (i *Ideas) InsertReply(commentUID, authorID, replyText string) (models.Reply, error) {
	op := "IdeasInsertReply"
	log := i.log.With(slog.String("op", op),
		slog.String("commentUID", commentUID),
	)
	log.Debug("inserting reply")

	if commentUID == "" {
		log.Error("commentUID is null")
		return models.Reply{}, nil
	}
	if authorID == "" {
		log.Error("AuthorID is null")
		return models.Reply{}, nil
	}
	if replyText == "" {
		log.Error("ReplyText is null")
		return models.Reply{}, nil
	}

	replyUID := uuid.New().String()

	reply := models.Reply{
		ReplyUID:   replyUID,
		CommentUID: commentUID,
		AuthorID:   authorID,
		ReplyText:  replyText,
	}

	err := i.repo.InsertCommentReply(reply)
	if err != nil {
		log.Error("failed to insert reply" + err.Error())
		return models.Reply{}, err
	}

	log.Info("successfully inserted reply, UUID:" + reply.ReplyUID)

	return reply, nil
}
