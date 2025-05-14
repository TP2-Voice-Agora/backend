package ideas

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/ictisagora/backend/internal/models"
	"log/slog"
	"testing"
)

// MockRepository — мок реализации Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ConnectDB(sourceURL string, log slog.Logger) error { return nil }
func (m *MockRepository) CloseConnectDB() error                             { return nil }
func (m *MockRepository) InsertUser(u models.User) error                    { return nil }
func (m *MockRepository) SelectUserByEmail(string) (models.User, error)     { return models.User{}, nil }
func (m *MockRepository) SelectPositions() ([]models.UserPosition, error)   { return nil, nil }
func (m *MockRepository) SelectUserByUID(uid string) (models.User, error)   { return models.User{}, nil }
func (m *MockRepository) InsertIdea(idea models.Idea) error {
	args := m.Called(idea)
	return args.Error(0)
}
func (m *MockRepository) SelectIdeas() ([]models.Idea, error) {
	args := m.Called()
	return args.Get(0).([]models.Idea), args.Error(1)
}
func (m *MockRepository) SelectIdeaByUID(uid string) (models.Idea, error) {
	args := m.Called(uid)
	return args.Get(0).(models.Idea), args.Error(1)
}
func (m *MockRepository) SelectUserIdeas(uid string, limit int) ([]models.Idea, error) {
	args := m.Called(uid, limit)
	return args.Get(0).([]models.Idea), args.Error(1)
}
func (m *MockRepository) InsertIdeaComment(comment models.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}
func (m *MockRepository) InsertCommentReply(reply models.Reply) error {
	args := m.Called(reply)
	return args.Error(0)
}
func (m *MockRepository) SelectIdeaComments(uid string) ([]models.Comment, error) {
	args := m.Called(uid)
	return args.Get(0).([]models.Comment), args.Error(1)
}
func (m *MockRepository) SelectCommentReplies(uid string) ([]models.Reply, error) {
	args := m.Called(uid)
	return args.Get(0).([]models.Reply), args.Error(1)
}
func (m *MockRepository) SelectIdeaCategories() ([]models.IdeaCategory, error) {
	args := m.Called()
	return args.Get(0).([]models.IdeaCategory), args.Error(1)
}
func (m *MockRepository) SelectIdeaStatuses() ([]models.IdeaStatus, error) {
	args := m.Called()
	return args.Get(0).([]models.IdeaStatus), args.Error(1)
}

func setupIdeasWithMocks(t *testing.T) (*Ideas, *MockRepository) {
	repo := new(MockRepository)
	cats := []models.IdeaCategory{{ID: 1, Name: "cat"}}
	stats := []models.IdeaStatus{{ID: 1, Name: "status"}}
	repo.On("SelectIdeaCategories").Return(cats, nil)
	repo.On("SelectIdeaStatuses").Return(stats, nil)
	ideas := New(*slog.Default(), repo)
	assert.NotNil(t, ideas)
	return ideas, repo
}

func TestGetIdeaCategories(t *testing.T) {
	ideas, _ := setupIdeasWithMocks(t)
	cats := ideas.GetIdeaCategories()
	assert.Equal(t, 1, len(cats))
}

func TestGetIdeaStatuses(t *testing.T) {
	ideas, _ := setupIdeasWithMocks(t)
	stats := ideas.GetIdeaStatuses()
	assert.Equal(t, 1, len(stats))
}

func TestGetAllIdeas_Success(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	expected := []models.Idea{{IdeaUID: uuid.New().String(), Name: "name"}}
	repo.On("SelectIdeas").Return(expected, nil)
	got, err := ideas.GetAllIdeas()
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetAllIdeas_RepoError(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	repo.On("SelectIdeas").Return([]models.Idea{}, errors.New("err"))
	got, err := ideas.GetAllIdeas()
	assert.Error(t, err)
	assert.Equal(t, []models.Idea{}, got) // исправлено: ожидаем пустой срез, а не nil
}

func TestGetIdeaByUID_EmptyUID(t *testing.T) {
	ideas, _ := setupIdeasWithMocks(t)
	result, err := ideas.GetIdeaByUID("")
	assert.NoError(t, err)
	assert.Equal(t, models.IdeaComment{}, result)
}

func TestGetIdeaByUID_Success(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	idea := models.Idea{IdeaUID: "id1"}
	comments := []models.Comment{{CommentUID: "c1", IdeaUID: "id1"}}
	replies := []models.Reply{{ReplyUID: "r1", CommentUID: "c1"}}

	repo.On("SelectIdeaByUID", "id1").Return(idea, nil)
	repo.On("SelectIdeaComments", "id1").Return(comments, nil)
	repo.On("SelectCommentReplies", "c1").Return(replies, nil)

	ic, err := ideas.GetIdeaByUID("id1")
	assert.NoError(t, err)
	assert.Equal(t, idea, ic.Idea)
	assert.Len(t, ic.CommentReplies, 1)
	assert.Equal(t, comments[0], ic.CommentReplies[0].Comment)
	assert.Equal(t, replies, ic.CommentReplies[0].Replies)
}

func TestGetAuthorIdeas(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	authors := []models.Idea{{IdeaUID: "i1"}, {IdeaUID: "i2"}}
	repo.On("SelectUserIdeas", "u1", 2).Return(authors, nil)
	result, err := ideas.GetAuthorIdeas("u1", 2)
	assert.NoError(t, err)
	assert.Equal(t, authors, result)
}

func TestGetIdeaComments(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	comms := []models.Comment{{CommentUID: "c"}}
	repo.On("SelectIdeaComments", "i1").Return(comms, nil)
	result, err := ideas.GetIdeaComments("i1")
	assert.NoError(t, err)
	assert.Equal(t, comms, result)
}

func TestGetCommentReplies(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	repls := []models.Reply{{ReplyUID: "r"}}
	repo.On("SelectCommentReplies", "c1").Return(repls, nil)
	result, err := ideas.GetCommentReplies("c1")
	assert.NoError(t, err)
	assert.Equal(t, repls, result)
}

func TestInsertIdea_Validation(t *testing.T) {
	ideas, _ := setupIdeasWithMocks(t)
	id, err := ideas.InsertIdea("", "body", "author", 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, models.Idea{}, id)
}

func TestInsertIdea_Success(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	repo.On("InsertIdea", mock.AnythingOfType("models.Idea")).Return(nil)
	idea, err := ideas.InsertIdea("n", "t", "a", 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "n", idea.Name)
	assert.Equal(t, "a", idea.Author)
}

func TestInsertComment_Validate(t *testing.T) {
	ideas, _ := setupIdeasWithMocks(t)
	c, err := ideas.InsertComment("", "uid", "txt")
	assert.NoError(t, err)
	assert.Equal(t, models.Comment{}, c)
}

func TestInsertComment_Success(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	repo.On("InsertIdeaComment", mock.AnythingOfType("models.Comment")).Return(nil)
	c, err := ideas.InsertComment("i", "a", "t")
	assert.NoError(t, err)
	assert.Equal(t, "i", c.IdeaUID)
}

func TestInsertReply_Validate(t *testing.T) {
	ideas, _ := setupIdeasWithMocks(t)
	r, err := ideas.InsertReply("", "author", "txt")
	assert.NoError(t, err)
	assert.Equal(t, models.Reply{}, r)
}

func TestInsertReply_Success(t *testing.T) {
	ideas, repo := setupIdeasWithMocks(t)
	repo.On("InsertCommentReply", mock.AnythingOfType("models.Reply")).Return(nil)
	r, err := ideas.InsertReply("comm", "auth", "txt")
	assert.NoError(t, err)
	assert.Equal(t, "auth", r.AuthorID)
}
