package http_server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gitlab.com/ictisagora/backend/internal/models"
	"gitlab.com/ictisagora/backend/internal/services/http-server/mware"
	i "gitlab.com/ictisagora/backend/internal/services/interfaces"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "gitlab.com/ictisagora/backend/docs" // путь, куда будет генерироваться swagger doc
)

// HTTPServer encapsulates the server dependencies and routes.
// for push
type HTTPServer struct {
	ideaService i.IdeaService
	authService i.AuthService
	userService i.UserService
	log         *slog.Logger
}

// NewHTTPServer creates and configures a new HTTPServer instance.
func NewHTTPServer(ideaService i.IdeaService, authService i.AuthService, userService i.UserService, log *slog.Logger) *HTTPServer {
	return &HTTPServer{
		ideaService: ideaService,
		authService: authService,
		userService: userService,
		log:         log,
	}
}

// SetupRoutes builds and returns an http.Handler with all routes and
// middleware.
func (s *HTTPServer) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	s.log.Info("Version 1.3")

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Post("/login", s.handleLogin)
		r.Post("/register", s.handleRegister)
		r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(mware.AuthMiddleware(s.authService.GetJWT(), s.log, s.userService))

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

		r.Get("/users/{uid}", s.handleGetUser)
		r.Post("/users/pfp", s.handleUploadUserPFP)
	})

	return r
}

// handleLogin
// @Summary      Аутентификация
// @Description  Аутентификация, возвращает jwt токен, который прикладывается ко всем (secure) рутам.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body  models.LoginRequest true  "Login data"
// @Success      200  {string}  string  "JWT token"
// @Failure      400  {string}  string  "Bad request"
// @Failure      405  {string}  string  "Invalid method"
// @Failure      500  {string}  string  "Failed to login"
// @Router       /login [post]
func (s *HTTPServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var body models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.log.Error("failed to decode request body", slog.String("error", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	jwtToken, userUID, err := s.authService.Login(body.Email, body.Password)
	if err != nil {
		s.log.Error("failed to log in user", slog.String("email", body.Email), slog.String("error", err.Error()))
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}
	var respBody struct {
		Token string `json:"Token"`
		Uid   string `json:"Uid"`
	}
	respBody.Token = jwtToken
	respBody.Uid = userUID
	resp, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)

}

// handleRegister
// @Summary      Регистрация
// @Description  Регистрация - будет только в админке
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        registerRequest  body  models.RegisterRequest true "Register data"
// @Success      200  {object}  map[string]string  "message: ok"
// @Failure      400  {string}  string  "Bad request"
// @Failure      405  {string}  string  "Invalid method"
// @Failure      500  {string}  string  "Failed to register"
// @Router       /register [post]
func (s *HTTPServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var body models.RegisterRequest
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
// @Summary      Категории идей(secure)
// @Description  Ручка категорий идей, в теории дергается один раз при первой загрузке страницы,
// так как никогда не обновляется
// @Tags         ideas
// @Produce      json
// @Success      200  {array}   models.IdeaCategory
// @Failure      405  {string}  string  "Invalid method"
// @Router       /ideas/categories [get]
func (s *HTTPServer) handleGetIdeaCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	categories := s.ideaService.GetIdeaCategories()
	resp, _ := json.Marshal(categories)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// handleGetIdeaStatuses
// @Summary      Статусы идей(secure)
// @Description  Ручка статусов идей, в теории дергается один раз при первой загрузке страницы,
// // так как никогда не обновляется
// @Tags         ideas
// @Produce      json
// @Success      200  {array}   models.IdeaStatus
// @Failure      405  {string}  string  "Invalid method"
// @Router       /ideas/statuses [get]
func (s *HTTPServer) handleGetIdeaStatuses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	statuses := s.ideaService.GetIdeaStatuses()
	resp, _ := json.Marshal(statuses)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// handleGetAllIdeas
// @Summary      Все идеи(secure)
// @Description  Возвращает все идеи списков без комментариев\ответов.
// @Tags         ideas
// @Produce      json
// @Success      200  {array}   models.Idea
// @Failure      500  {string}  string  "Failed to get ideas"
// @Router       /ideas [get]
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
// @Summary      Конкретная идея(secure)
// @Description  Возвращает идею по UID, уже с комментариями\ответами
// @Tags         ideas
// @Produce      json
// @Param        uid   path      string  true  "Idea UID"
// @Success      200   {object}  models.IdeaComment
// @Failure      405   {string}  string  "Invalid method"
// @Failure      500   {string}  string  "Failed to get idea by UID"
// @Router       /ideas/{uid} [get]
func (s *HTTPServer) handleGetIdeaByUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

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

// handleInsertIdea
// @Summary      Вставка новой идеи(secure)
// @Description  Вставляет идею, и возвращает ее со всеми заполненными полями
// @Tags         ideas
// @Accept       json
// @Produce      json
// @Param        idea  body  models.InsertIdeaRequest true  "Idea data"
// @Success      201  {object}  models.Idea
// @Failure      400  {string}  string  "Bad request"
// @Failure      500  {string}  string  "Failed to create idea"
// @Router       /ideas [post]
func (s *HTTPServer) handleInsertIdea(w http.ResponseWriter, r *http.Request) {
	var body models.InsertIdeaRequest
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

// handleInsertComment
// @Summary      Вставка комментария(secure)
// @Description  Вставляет коммент и возвращает его.
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        comment body models.InsertCommentRequest true "Comment data"
// @Success      201  {object}  models.Comment
// @Failure      400  {string}  string  "Bad request"
// @Failure      405  {string}  string  "Invalid method"
// @Failure      500  {string}  string  "Failed to create comment"
// @Router       /comments [post]
func (s *HTTPServer) handleInsertComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var body models.InsertCommentRequest
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

// handleInsertReply
// @Summary      Вставка ответа
// @Description  Вставляет новый ответ, и возвращает его
// @Tags         replies
// @Accept       json
// @Produce      json
// @Param        reply  body  models.InsertReplyRequest true "Reply data"
// @Security     JWTAuth
// @Success      201  {object}  models.Reply
// @Failure      400  {string}  string  "Bad request"
// @Failure      405  {string}  string  "Invalid method"
// @Failure      500  {string}  string  "Failed to create reply"
// @Router       /replies [post]
func (s *HTTPServer) handleInsertReply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var body models.InsertReplyRequest
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

// handleGetUser
// @Summary      Получение юзера по UID
// @Description  Возвращает данные пользователя по UID.
// @Tags         users
// @Produce      json
// @Param        uid   path      string  true  "User UID"
// @Success      200   {object}  models.User
// @Failure      405   {string}  string  "Invalid method"
// @Failure      500   {string}  string  "Failed to get user"
// @Router       /users/{uid} [get]
func (s *HTTPServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	uid := chi.URLParam(r, "uid")

	user, err := s.userService.GetUserByUID(uid)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
	}

	resp, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

// handleUploadUserPFP
// @Summary      Загрузка PFP
// @Description  Загрузка новой аватарки для юзера.
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        profile_picture  formData  file  true  "Profile picture file"
// @Success      200  {object}  map[string]string  "url to uploaded picture"
// @Failure      400  {string}  string  "No file uploaded or bad request"
// @Failure      405  {string}  string  "Invalid method"
// @Failure      500  {string}  string  "Failed to upload file"
// @Router       /users/pfp [post]
func (s *HTTPServer) handleUploadUserPFP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
	}
	defer file.Close()

	userUID := r.Context().Value(mware.ContextUserUID).(string)

	url, err := s.userService.UploadPFP(file, header, userUID)
	if err != nil {
		s.log.Error("failed to upload user profile", slog.String("error", err.Error()))
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
	}

	resp, _ := json.Marshal(map[string]string{"url": url})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}
