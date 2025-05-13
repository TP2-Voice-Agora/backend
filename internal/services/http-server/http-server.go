package http_server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/ictisagora/backend/internal/models"
	"gitlab.com/ictisagora/backend/internal/services/http-server/mware"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// IdeaService is a placeholder interface representing your business logic layer.
type IdeaService interface {
	GetIdeaCategories() []models.IdeaCategory
	GetIdeaStatuses() []models.IdeaStatus
	GetAllIdeas() ([]models.Idea, error)
	GetIdeaByUID(uid string) (models.IdeaComment, error)
	GetAuthorIdeas(uid string, limit int) ([]models.Idea, error)
	InsertIdea(name string, text string, author string, status int, category int) (models.Idea, error)
	InsertComment(ideaUID, authorUID, commentText string) (models.Comment, error)
	InsertReply(commentUID, authorID, replyText string) (models.Reply, error)
}

type AuthService interface {
	Register(u models.User) error
	Login(email string, password string) (string, error)
	GetJWT() string
}

// HTTPServer encapsulates the server dependencies and routes.
type HTTPServer struct {
	ideaService IdeaService
	authService AuthService
	log         *slog.Logger
}

// NewHTTPServer creates and configures a new HTTPServer instance.
func NewHTTPServer(ideaService IdeaService, authService AuthService, log *slog.Logger) *HTTPServer {
	return &HTTPServer{
		ideaService: ideaService,
		authService: authService,
		log:         log,
	}
}

// SetupRoutes builds and returns an http.Handler with all routes and
// middleware.
func (s *HTTPServer) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Post("/login", s.handleLogin)
		r.Post("/register", s.handleRegister)
	})

	r.Group(func(r chi.Router) {
		r.Use(mware.AuthMiddleware(s.authService.GetJWT(), s.log))

		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)

		r.Get("/ideas/categories", s.handleGetIdeaCategories)
		r.Get("/ideas/statuses", s.handleGetIdeaStatuses)
		r.Get("/ideas", s.handleGetAllIdeas)
		r.Get("/ideas/{uid}", s.handleGetIdeaByUID)
		r.Post("/ideas", s.handleInsertIdea)
		r.Post("/comments", s.handleInsertComment)
		r.Post("/replies", s.handleInsertReply)
	})

	return r
}

// handleLogin
// returns jwt token
func (s *HTTPServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.log.Error("failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	jwtToken, err := s.authService.Login(body.Email, body.Password)
	if err != nil {
		s.log.Error("failed to log in user", slog.String("email", body.Email), slog.String("error", err.Error()))
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(jwtToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)

}

// handeRegister
func (s *HTTPServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		PositionID int    `json:"positionID"`
		Name       string `json:"name"`
		Surname    string `json:"surname"`
	}
	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err := s.authService.Register(models.User{
		Email:      body.Email,
		Password:   body.Password,
		PositionID: body.PositionID,
		Name:       body.Name,
		Surname:    body.Surname,
	})
	if err != nil {
		s.log.Error("failed to register user", slog.String("email", body.Email), slog.String("error", err.Error()))
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(map[string]string{"message": "ok"})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

// handleGetIdeaCategories
// returns []models.IdeaCategory
func (s *HTTPServer) handleGetIdeaCategories(w http.ResponseWriter, r *http.Request) {
	categories := s.ideaService.GetIdeaCategories()
	resp, _ := json.Marshal(categories)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// handleGetIdeaStatuses
// returns []models.IdeaStatus
func (s *HTTPServer) handleGetIdeaStatuses(w http.ResponseWriter, r *http.Request) {
	statuses := s.ideaService.GetIdeaStatuses()
	resp, _ := json.Marshal(statuses)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// handleGetAllIdeas
// returns models.Ideas
func (s *HTTPServer) handleGetAllIdeas(w http.ResponseWriter, r *http.Request) {
	ideas, err := s.ideaService.GetAllIdeas()
	if err != nil {
		http.Error(w, "Failed to get ideas", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(ideas)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// handleGetIdeaByUID
// return models.IdeaComment
func (s *HTTPServer) handleGetIdeaByUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")

	ideaComment, err := s.ideaService.GetIdeaByUID(uid)
	if err != nil {
		http.Error(w, "Failed to get idea by UID", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(ideaComment)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// handleInsertIdea handles creation of a new idea.
// gets requestBody
// returns models.Idea
func (s *HTTPServer) handleInsertIdea(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Name     string `json:"name"`
		Text     string `json:"text"`
		Author   string `json:"author"`
		Status   int    `json:"status"`
		Category int    `json:"category"`
	}

	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.log.Error("failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	body.Author = r.Context().Value(mware.ContextUserUID).(string)

	newIdea, err := s.ideaService.InsertIdea(
		body.Name, body.Text, body.Author, body.Status, body.Category,
	)
	if err != nil {
		http.Error(w, "Failed to create idea", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(newIdea)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}

func (s *HTTPServer) handleInsertComment(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		IdeaUID     string `json:"ideaUID"`
		CommentText string `json:"commentText"`
	}

	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	authorUID := r.Context().Value(mware.ContextUserUID).(string)

	newComment, err := s.ideaService.InsertComment(body.IdeaUID, authorUID, body.CommentText)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	resp, _ := json.Marshal(newComment)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)

}

func (s *HTTPServer) handleInsertReply(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		CommentUID string `json:"commentUID"`
		ReplyText  string `json:"replyText"`
	}

	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	authorID := r.Context().Value(mware.ContextUserUID).(string)
	newReply, err := s.ideaService.InsertReply(body.CommentUID, authorID, body.ReplyText)
	if err != nil {
		http.Error(w, "Failed to create reply", http.StatusInternalServerError)
	}

	resp, _ := json.Marshal(newReply)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(resp)
}
